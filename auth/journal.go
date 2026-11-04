package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
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
// to send, a batch at a time, again every minute until it's taken.
// ponytail: in memory, the lines of a restart before auth took them are
// lost; and a clock stepped back makes the next batches look sent again.
type Journal struct {
	send    func(JournalBatch) error
	mu      sync.Mutex
	lines   []Line
	pending *JournalBatch // sent, not yet taken: sent again as is
	seq     int64
	wake    chan struct{}
}

func NewJournal(send func(JournalBatch) error) *Journal {
	return &Journal{send: send, wake: make(chan struct{}, 1)}
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
	j.lines = append(j.lines, l)
	if len(j.lines) > journalQueued {
		j.lines = j.lines[len(j.lines)-journalQueued:]
	}
	j.mu.Unlock()
	select {
	case j.wake <- struct{}{}:
	default:
	}
}

// Run sends the queued lines until ctx ends.
func (j *Journal) Run(ctx context.Context) {
	tick := time.NewTicker(time.Minute)
	defer tick.Stop()
	for {
		for j.flush() {
		}
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
	j.pending = nil
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
