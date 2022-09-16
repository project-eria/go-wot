package protocolWebSocket

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/interaction"
	"github.com/rs/zerolog/log"
)

type WsClient struct {
	schemes []string
	wait    sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewClient() *WsClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &WsClient{
		schemes: []string{"ws"},
		ctx:     ctx,
		cancel:  cancel,
	}
}

type wsConn struct {
	mu          sync.RWMutex
	wsURL       string
	isConnected bool
	dialErr     error
	connWait    connWait
	sub         *consumer.Subscription
	listener    consumer.Listener
	*websocket.Conn
}

func (c *WsClient) GetSchemes() []string {
	return c.schemes
}

// ReadResource get a JSON data from HTTP GET request
func (c *WsClient) ReadResource(form *interaction.Form) (interface{}, error) {
	return nil, errors.New("not implemented")
}

// WriteResource send JSON data using HTTP PUT request
func (c *WsClient) WriteResource(form *interaction.Form, value interface{}) (interface{}, error) {
	return nil, errors.New("not implemented")
}

// InvokeResource send JSON data using HTTP POST request
func (c *WsClient) InvokeResource(form *interaction.Form, value interface{}) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c *WsClient) SubscribeResource(form *interaction.Form, sub *consumer.Subscription, listener consumer.Listener) error {
	c.wait.Add(1)
	go func() {
		c.connectWebSocket(form.Href, sub, listener)
		c.wait.Done()
	}()
	return nil
}

func (c *WsClient) Stop() {
	c.cancel()
	// Wait for the child goroutine to finish, which will only occur when
	// the child process has stopped and the call to cmd.Wait has returned.
	// This prevents main() exiting prematurely.
	c.wait.Wait()
}

// connectWebSocket Connect the thing using the WebSocket access
func (c *WsClient) connectWebSocket(wsURL string, sub *consumer.Subscription, listener consumer.Listener) {
	wsc := &wsConn{
		wsURL:    wsURL,
		connWait: newConnWait(),
		sub:      sub,
		listener: listener,
	}
	for {
		select {
		case <-c.ctx.Done():
			log.Warn().Str("url", wsURL).Msg("[protocolWebSocket:ConnectWebSocket] Connecting interrupted by user")
			return
		case <-wsc.connect():
			if wsc.IsConnected() { // Should come here connected
				select {
				case <-c.ctx.Done():
					wsc.gracefullyShutdown()
					return
				case err := <-wsc.read():
					if err != nil {
						if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
							log.Info().Str("url", wsURL).Msg("[protocolWebSocket:ConnectWebSocket] Normal Closure")
						} else {
							log.Error().Err(err).Str("url", wsURL).Msg("[protocolWebSocket:ConnectWebSocket]")
						}
						wsc.close()
					}
					break // Break the read message loop to go back on connect() to try to reconnect
				}
			}
			// If we go out from the listen() loop we try to reconnect
			log.Info().Str("url", wsURL).Msg("[protocolWebSocket:ConnectWebSocket] Trying to reconnect...")
		}
	}
}

func (c *wsConn) connect() <-chan bool {
	if c == nil {
		log.Error().Msg("[protocolWebSocket:connect] nil connection")
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
				log.Info().Str("url", c.wsURL).Msg("[protocolWebSocket:connect] WebSocket connection successfully established")
				c.connWait.reset()
				success <- true
				return
			}
			// If err != nil
			nextDuration := c.connWait.nextDuration()
			log.Error().Err(err).Str("url", c.wsURL).Msgf("[protocolWebSocket:connect] WebSocket connect will try again in %s.", nextDuration.String())

			time.Sleep(nextDuration)
		}
	}()
	return success
}

// read monitor the WebSocket Messages
func (c *wsConn) read() <-chan error {
	if c == nil {
		log.Error().Msg("[protocolWebSocket:read] nil connection")
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

			log.Trace().Interface("message", data).Msgf("[protocolWebSocket:read] Received from WebSocket")
			if c.listener != nil {
				go c.listener(data, nil)
			}
		}
	}()
	return result
}

func (c *wsConn) get() *websocket.Conn {
	if c == nil {
		log.Error().Msg("[protocolWebSocket:get] nil connection")
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Conn
}

// setIsConnected sets state for isConnected
func (c *wsConn) setIsConnected(state bool) {
	if c == nil {
		log.Error().Msg("[protocolWebSocket:setIsConnected] nil connection")
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isConnected = state
}

// IsConnected returns the WebSocket connection state
func (c *wsConn) IsConnected() bool {
	if c == nil {
		log.Error().Msg("[protocolWebSocket:IsConnected] nil connection")
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
		log.Error().Msg("[protocolWebSocket:gracefullyShutdown] nil connection")
		return
	}
	log.Info().Msg("[protocolWebSocket:gracefullyShutdown] Sending WebSocket Closing message to server")
	err := c.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Time{})
	if err != nil {
		log.Error().Err(err).Msg("[protocolWebSocket:gracefullyShutdown] WebSocket Close")
		return
	}
	c.Close()
}

// close closes the underlying network connection without
// sending or waiting for a close frame.
func (c *wsConn) close() {
	if c == nil {
		log.Error().Msg("[protocolWebSocket:close] nil connection")
		return
	}
	if c.get() != nil {
		c.mu.Lock()
		c.Conn.Close()
		c.mu.Unlock()
	}

	c.setIsConnected(false)
}
