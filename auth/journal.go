package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Journal lines of the other sources (AUTH-DESIGN §4.5, step 5e-3): the
// master and eria-ui-server hand theirs to auth in signed batches. auth sets
// the source from the key that signed, and the retention from the type.
// Never a PIN, a code, a whole token, a key, a challenge nor a client's IP.

// JournalPrefix is the domain of a signed JournalBatch.
const JournalPrefix = "eria-journal-batch-v1"

// Line is a journal line as a source writes it.
type Line struct {
	At      int64  `json:"at"`                // unix seconds, the source's clock
	Type    string `json:"type"`              // order.validated, relay.denied…
	Subject string `json:"subject,omitempty"` // profil:<id>, term:<id>…
	Object  string `json:"object,omitempty"`
	Result  string `json:"result,omitempty"` // ok (default), denied
	Reason  string `json:"reason,omitempty"`
	Jti     string `json:"jti,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

// JournalBatch is what a source signs. Seq only goes up: a batch sent
// again (its answer lost) is taken once.
type JournalBatch struct {
	Seq   int64  `json:"seq"`
	Lines []Line `json:"lines"`
}

const (
	journalBatchLines = 32   // with the fields cut, under auth's 64 KiB per batch
	journalQueued     = 1024 // the oldest lines go past it
	journalField      = 128  // bytes kept of each field, journalDetail for Detail
	journalDetail     = 256
)

// cut keeps the first n bytes of s, on a character boundary.
func cut(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return strings.ToValidUTF8(s[:n], "")
}

// Journal queues lines, never blocking the caller, and hands them in order
// to send, a batch at a time, again every minute until it's taken. A refusal
// repeated within a minute (same type, subject, object, reason) is counted,
// not queued: one line "repeated N" follows once the minute is over.
// ponytail: in memory unless Persist names a file; and a clock stepped back
// makes the next batches look sent again.
type Journal struct {
	send    func(JournalBatch) error
	mu      sync.Mutex
	lines   []Line
	pending *JournalBatch // sent, not yet taken: sent again as is
	seq     int64
	wake    chan struct{}
	repeats map[string]*repeat // refusals of the last minute
	file    string             // "": nothing kept
	dirty   bool               // changed since written
}

type repeat struct {
	line  Line
	since time.Time
	n     int // times again since
}

// repeatWindow is how long a refusal coalesces its repeats.
const repeatWindow = time.Minute

func NewJournal(send func(JournalBatch) error) *Journal {
	return &Journal{send: send, wake: make(chan struct{}, 1), repeats: map[string]*repeat{}}
}

// Add queues a line, stamped now when At is zero.
func (j *Journal) Add(l Line) {
	if l.At == 0 {
		l.At = time.Now().Unix()
	}
	l.Type, l.Subject, l.Object = cut(l.Type, journalField), cut(l.Subject, journalField), cut(l.Object, journalField)
	l.Result, l.Reason, l.Jti = cut(l.Result, journalField), cut(l.Reason, journalField), cut(l.Jti, journalField)
	l.Detail = cut(l.Detail, journalDetail)
	j.mu.Lock()
	if l.Result == "denied" {
		key := l.Type + "\x00" + l.Subject + "\x00" + l.Object + "\x00" + l.Reason
		if r := j.repeats[key]; r != nil && time.Since(r.since) < repeatWindow {
			r.n++
			j.mu.Unlock()
			return
		}
		j.endRepeats(true)
		j.repeats[key] = &repeat{line: l, since: time.Now()}
	}
	j.queue(l)
	j.mu.Unlock()
	select {
	case j.wake <- struct{}{}:
	default:
	}
}

// queue appends l, the oldest lines going past journalQueued; j.mu held.
func (j *Journal) queue(l Line) {
	j.lines = append(j.lines, l)
	if len(j.lines) > journalQueued {
		j.lines = j.lines[len(j.lines)-journalQueued:]
	}
	j.dirty = true
}

// endRepeats queues "repeated N" for the refusals whose minute is over (all
// of them unless onlyOver); j.mu held.
func (j *Journal) endRepeats(onlyOver bool) {
	for key, r := range j.repeats {
		if onlyOver && time.Since(r.since) < repeatWindow {
			continue
		}
		if r.n > 0 {
			l := r.line
			l.At, l.Jti, l.Detail = time.Now().Unix(), "", cut(fmt.Sprintf("%s repeated %d", l.Detail, r.n), journalDetail)
			j.queue(l)
		}
		delete(j.repeats, key)
	}
}

// Persist keeps the queue in file: read now, written after each round of
// sends. Call before Run.
func (j *Journal) Persist(file string) {
	var saved struct {
		Seq     int64
		Pending *JournalBatch
		Lines   []Line
	}
	if data, err := os.ReadFile(file); err == nil && json.Unmarshal(data, &saved) == nil {
		j.seq, j.pending, j.lines = saved.Seq, saved.Pending, append(saved.Lines, j.lines...)
	}
	j.file = file
}

// save writes the queue when it changed.
func (j *Journal) save() {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.file == "" || !j.dirty {
		return
	}
	data, _ := json.Marshal(map[string]any{"Seq": j.seq, "Pending": j.pending, "Lines": j.lines})
	tmp := j.file + ".tmp"
	if os.WriteFile(tmp, data, 0o600) == nil && os.Rename(tmp, j.file) == nil {
		j.dirty = false
	}
}

// Run sends the queued lines until ctx ends.
func (j *Journal) Run(ctx context.Context) {
	tick := time.NewTicker(time.Minute)
	defer tick.Stop()
	for {
		j.mu.Lock()
		j.endRepeats(true)
		j.mu.Unlock()
		for j.flush() {
		}
		j.save()
		select {
		case <-ctx.Done():
			return
		case <-j.wake:
		case <-tick.C:
		}
	}
}

// flush sends one batch, the pending one or the next lines; false when
// there was none or it wasn't taken.
func (j *Journal) flush() bool {
	j.mu.Lock()
	if j.pending == nil && len(j.lines) > 0 {
		n := min(len(j.lines), journalBatchLines)
		j.seq = max(j.seq+1, time.Now().UnixMicro())
		j.pending = &JournalBatch{Seq: j.seq, Lines: append([]Line(nil), j.lines[:n]...)}
		j.lines = j.lines[n:]
	}
	b := j.pending
	j.mu.Unlock()
	if b == nil || j.send(*b) != nil {
		return false
	}
	j.mu.Lock()
	j.pending, j.dirty = nil, true
	j.mu.Unlock()
	return true
}

// SignDocumentP256 signs v like SignDocument, with a P-256 key (a relay's
// TLS key): r ‖ s over SHA-256(prefix ‖ payload).
func SignDocumentP256(priv *ecdsa.PrivateKey, prefix string, v any) (Signed, error) {
	payload, err := json.Marshal(v)
	if err != nil {
		return Signed{}, err
	}
	h := sha256.Sum256(append([]byte(prefix), payload...))
	r, s, err := ecdsa.Sign(rand.Reader, priv, h[:])
	if err != nil {
		return Signed{}, err
	}
	sig := make([]byte, 64)
	r.FillBytes(sig[:32])
	s.FillBytes(sig[32:])
	return Signed{Payload: b64.EncodeToString(payload), Sig: b64.EncodeToString(sig)}, nil
}

// JournalErasePrefix is the domain of a JournalErase, signed by a deployment
// key (the master), once the admin validated the order (§4.5 « effacement »).
const JournalErasePrefix = "eria-journal-erase-v1"

// JournalErase erases whole UTC days up to Until (YYYY-MM-DD, included):
// the comfort lines (30 days), or all of them. At (unix seconds) only goes
// up: a message replayed erases nothing more.
type JournalErase struct {
	Until string `json:"until"`
	All   bool   `json:"all"`
	At    int64  `json:"at"`
}
