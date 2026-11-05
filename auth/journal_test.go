package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A batch not taken is sent again as is; a P-256 signed batch opens.
func TestJournal(t *testing.T) {
	var sent []JournalBatch
	fail := true
	j := NewJournal(func(b JournalBatch) error {
		sent = append(sent, b)
		if fail {
			return errors.New("auth away")
		}
		return nil
	})
	for range journalBatchLines + 1 {
		j.Add(Line{Type: "relay.connect"})
	}
	assert.False(t, j.flush())
	j.Add(Line{Type: "relay.close"})
	fail = false
	assert.True(t, j.flush())
	assert.True(t, j.flush())
	assert.False(t, j.flush())
	require.Len(t, sent, 3)
	assert.Equal(t, sent[0], sent[1], "sent again as is")
	assert.Len(t, sent[2].Lines, 2)
	assert.Greater(t, sent[2].Seq, sent[1].Seq)

	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	raw, _ := priv.PublicKey.Bytes()
	pub, err := ParsePublicKey(P256, raw)
	require.NoError(t, err)
	doc, err := SignDocumentP256(priv, JournalPrefix, sent[2])
	require.NoError(t, err)
	var back JournalBatch
	require.NoError(t, doc.Open(JournalPrefix, []PublicKey{pub}, &back))
	assert.Equal(t, sent[2], back)
}

// Refusals repeated within a minute: one line, then "repeated N"; the queue
// kept in a file across a restart.
func TestJournalRepeatsAndFile(t *testing.T) {
	var sent []JournalBatch
	j := NewJournal(func(b JournalBatch) error { sent = append(sent, b); return nil })
	file := t.TempDir() + "/journal.json"
	j.Persist(file)
	refused := Line{Type: "guard.denied", Result: "denied", Reason: "step 2", Detail: "action open"}
	for range 5 {
		j.Add(refused)
	}
	j.Add(Line{Type: "guard.denied", Result: "denied", Reason: "step 7"})
	assert.Len(t, j.lines, 2, "the repeats counted, not queued")
	for _, r := range j.repeats {
		r.since = r.since.Add(-repeatWindow)
	}
	j.endRepeats(true)
	require.Len(t, j.lines, 3)
	assert.Equal(t, "action open repeated 4", j.lines[2].Detail)

	j.save()
	again := NewJournal(nil)
	again.Persist(file)
	assert.Equal(t, j.lines, again.lines, "kept across a restart")
	assert.True(t, j.flush())
	j.save()
	again = NewJournal(nil)
	again.Persist(file)
	assert.Empty(t, again.lines, "sent: not kept")
	assert.Equal(t, sent[0].Seq, again.seq)
}
