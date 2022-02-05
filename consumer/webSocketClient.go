package consumer

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type wsConn struct {
	mu          sync.RWMutex
	wsURL       string
	isConnected bool
	dialErr     error
	connWait    connWait
	sub         *subscription
	listener    Listener
	*websocket.Conn
}

// connectWebSocket Connect the thing using the WebSocket access
func connectWebSocket(wsURL string, sub *subscription, listener Listener, ctx context.Context) {
	wsc := &wsConn{
		wsURL:    wsURL,
		connWait: newConnWait(),
		sub:      sub,
		listener: listener,
	}
	for {
		select {
		case <-ctx.Done():
			log.Warn().Str("url", wsURL).Msg("[consumer:ConnectWebSocket] Connecting interrupted by user")
			return
		case <-wsc.connect():
			if wsc.IsConnected() { // Should come here connected
				select {
				case <-ctx.Done():
					wsc.gracefullyShutdown()
					return
				case err := <-wsc.read():
					if err != nil {
						if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
							log.Info().Str("url", wsURL).Msg("[consumer:ConnectWebSocket] Normal Closure")
						} else {
							log.Error().Err(err).Str("url", wsURL).Msg("[consumer:ConnectWebSocket]")
						}
						wsc.close()
					}
					break // Break the read message loop to go back on connect() to try to reconnect
				}
			}
			// If we go out from the listen() loop we try to reconnect
			log.Info().Str("url", wsURL).Msg("[consumer:ConnectWebSocket] Trying to reconnect...")
		}
	}
}

func (c *wsConn) connect() <-chan bool {
	if c == nil {
		log.Error().Msg("[consumer:connect] nil connection")
		return nil
	}
	success := make(chan bool, 1)
	go func() {
		for {
			ws := websocket.DefaultDialer
			wsConn, _, err := ws.Dial(c.wsURL, http.Header{"Sec-WebSocket-Protocol": []string{"webthing"}})

			c.mu.Lock()
			c.Conn = wsConn
			c.dialErr = err
			c.isConnected = err == nil
			c.mu.Unlock()

			if err == nil {
				log.Info().Str("url", c.wsURL).Msg("[consumer:connect] WebSocket connection successfully established")
				c.connWait.reset()
				success <- true
				return
			}
			// If err != nil
			nextDuration := c.connWait.nextDuration()
			log.Error().Err(err).Str("url", c.wsURL).Msgf("[consumer:connect] WebSocket connect will try again in %s.", nextDuration.String())

			time.Sleep(nextDuration)
		}
	}()
	return success
}

// read monitor the WebSocket Messages
func (c *wsConn) read() <-chan error {
	if c == nil {
		log.Error().Msg("[consumer:read] nil connection")
		return nil
	}
	result := make(chan error, 1)
	go func() {
		for {
			var data interface{}
			err := c.Conn.ReadJSON(&data)
			if err != nil {
				go c.listener(nil, err)
				result <- err
				return
			}

			log.Trace().Interface("message", data).Msgf("[consumer:read] Received from WebSocket")
			if c.listener != nil {
				go c.listener(data, nil)
			}
		}
	}()
	return result
}

func (c *wsConn) get() *websocket.Conn {
	if c == nil {
		log.Error().Msg("[consumer:get] nil connection")
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Conn
}

// setIsConnected sets state for isConnected
func (c *wsConn) setIsConnected(state bool) {
	if c == nil {
		log.Error().Msg("[consumer:setIsConnected] nil connection")
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isConnected = state
}

// IsConnected returns the WebSocket connection state
func (c *wsConn) IsConnected() bool {
	if c == nil {
		log.Error().Msg("[consumer:IsConnected] nil connection")
		return false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.isConnected
}

// gracefullyShutdown cleanly close the connection by sending a close message
// [TODO] and then waiting (with timeout) for the server to close the connection.
func (c *wsConn) gracefullyShutdown() {
	if c == nil {
		log.Error().Msg("[consumer:gracefullyShutdown] nil connection")
		return
	}
	log.Info().Msg("[consumer:gracefullyShutdown] Sending WebSocket Closing message to server")
	err := c.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Time{})
	if err != nil {
		log.Error().Err(err).Msg("[consumer:gracefullyShutdown] WebSocket Close")
		return
	}
	c.Close()
}

// close closes the underlying network connection without
// sending or waiting for a close frame.
func (c *wsConn) close() {
	if c == nil {
		log.Error().Msg("[consumer:close] nil connection")
		return
	}
	if c.get() != nil {
		c.mu.Lock()
		c.Conn.Close()
		c.mu.Unlock()
	}

	c.setIsConnected(false)
}
