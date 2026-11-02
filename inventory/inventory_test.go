package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrigin(t *testing.T) {
	assert.Equal(t, "http://etage:80", Origin("http://etage/bureau-volet"))
	assert.Equal(t, "http://etage:80", Origin("ws://ETAGE:80/bureau-volet/properties/x"), "a Thing's WebSocket is its HTTP server")
	assert.Equal(t, "http://10.0.0.5:8080", Origin("http://10.0.0.5:8080/"))
	assert.NotEqual(t, Origin("http://etage/"), Origin("http://etage:8080/"))
	assert.Empty(t, Origin("/relative"))
}
