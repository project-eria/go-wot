package protocolWebSocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/interaction"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeSendsToken(t *testing.T) {
	got := make(chan string, 1)
	upgrader := websocket.Upgrader{Subprotocols: []string{"webthing"}}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got <- r.Header.Get("Authorization")
		if c, err := upgrader.Upgrade(w, r, nil); err == nil {
			defer c.Close()
			c.ReadMessage() // until the client stops
		}
	}))
	defer srv.Close()

	c := NewClient()
	c.Token = func(string) string { return "tok" }
	form := &interaction.Form{Href: "ws" + strings.TrimPrefix(srv.URL, "http")}
	c.SubscribeResource(form, nil, &consumer.Subscription{}, func(any, error) {})
	select {
	case auth := <-got:
		assert.Equal(t, "Bearer tok", auth, "the service token goes in the handshake")
	case <-time.After(5 * time.Second):
		t.Fatal("no handshake")
	}
	c.Stop()
}
