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
