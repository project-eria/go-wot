package producer

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    []string{"webthing"},
}

type wsConnection struct {
	mu sync.RWMutex
	*websocket.Conn
}

func (c *wsConnection) jsonWSRenderer(content interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.WriteJSON(content)
}

func (c *wsConnection) errorWSRenderer(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.WriteJSON(map[string]string{"error": message})
}

func (c *wsConnection) Close() error {
	closeNormalClosure := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	if err := c.WriteControl(websocket.CloseMessage, closeNormalClosure, time.Now().Add(time.Second)); err != nil {
		return err
	}
	c.Close()
	return nil
}

func (t *ExposedThing) gracefullWSShutdown() {
	for _, p := range t.exposedProperties {
		p.mu.RLock()
		conns := p.observersProperties
		p.mu.RUnlock()
		for key, wsConn := range conns {
			log.Trace().Str("key", key).Msg("[ExposedProperty:gracefullWSShutdown] Send Close message")
			err := wsConn.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				time.Time{})
			if err != nil {
				log.Error().Str("key", key).Err(err).Msg("[ExposedProperty:gracefullWSShutdown] Sending error")
			}
			delete(p.observersProperties, key)
			t._wait.Done()
		}
	}
}
