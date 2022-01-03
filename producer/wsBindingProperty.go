package producer

import (
	"net/http"
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
	conn *websocket.Conn
	mu   sync.RWMutex
}

func (h *propertyHandler) webSocket(name string, w http.ResponseWriter, r *http.Request) {
	// // webthing SubProtocol do not exists
	// if r.Header.Get("Sec-Websocket-Protocol") != "webthing" {
	// 	log.Error().Msg("[producer:webSocket] Connection not using webthing protocol")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	io.WriteString(w, "Connection not using webthing protocol")
	// 	return
	// }

	// TODO Handle Origin for debug plugins
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, http.Header{"Sec-WebSocket-Protocol": []string{r.Header.Get("Sec-WebSocket-Protocol")}})
	if err != nil {
		log.Warn().Str("uri", r.RequestURI).Err(err).Msg("[producer:webSocket] WebSocket Upgrade")
		return
	}
	r.Header.Get("Sec-Websocket-Key")
	key := r.Header.Get("Sec-Websocket-Key")
	wsConn := &wsConnection{conn: conn}
	if err := h.AddWSPropertyObserver(name, key, wsConn); err != nil {
		wsConn.errorWSRenderer(err.Error())
		wsConn.Close()
		return
	}

	// TODO
	// wsConn := h.addWSConnection(key, conn)
	// for {
	// 	message := wsMessage{key: key, thing: h.Td}
	// 	err := conn.ReadJSON(&message)
	// 	if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
	// 		log.Debug().Str("key", key).Msg("[producer:webSocket] WebSocket Normal Closure")
	// 		h.removeWSConnection(key)
	// 		return
	// 	}
	// 	if err != nil {
	// 		log.Error().Str("key", key).Err(err).Msg("[producer:webSocket] WebSocket error")
	// 		h.removeWSConnection(key)
	// 		return
	// 	}
	// 	log.Trace().Str("key", key).Msgf("[producer:webSocket] Received WebSocket message: %#v", message)
	// 	h.processRxMsg(wsConn, &message)
	// }
}

// processRxMsg processes incoming messages
// func (h *propertyHandler) processRxMsg(wsConn *wsConnection, message *wsMessage) {
// 	log.Debug().Str("key", message.key).Str("type", message.MessageType).Msg("[producer:processRxMsg] Processing WS request")
// 	switch message.MessageType {
// 	case "setProperty":
// 		content, err := message.thing.processSetProperties(message.Data)
// 		if err != nil {
// 			wsConn.errorWSRenderer(err.Error())
// 			log.Error().Str("key", message.key).Err(err).Msg("[producer:processRxMsg] SetProperty request")
// 			break
// 		}
// 		wsConn.jsonWSRenderer(content)
// 		break
// 	default:
// 		wsConn.errorWSRenderer("Unsupported RX request type")
// 		log.Error().Str("key", message.key).Str("type", message.MessageType).Msg("[producer:processRxMsg] Unsupported RX request type")
// 	}
// }

// processTxMsg processes messages to be send
func (p *ExposedProperty) WSProcessTxMsg(message interface{}) {
	log.Debug().Msg("[ExposedProperty:processTxMsg] Processing WS request")
	// Send the message to all ws connections
	for _, wsConn := range p.observersProperties {
		err := wsConn.jsonWSRenderer(message)
		if err != nil {
			log.Error().Err(err).Msg("[ExposedProperty:processTxMsg] Sending propertyStatus")
		}
	}
}

func (c *wsConnection) jsonWSRenderer(content interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(content)
}

func (c *wsConnection) errorWSRenderer(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn.WriteJSON(map[string]string{"error": message})
}

func (c *wsConnection) Close() error {
	closeNormalClosure := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	if err := c.conn.WriteControl(websocket.CloseMessage, closeNormalClosure, time.Now().Add(time.Second)); err != nil {
		return err
	}
	c.conn.Close()
	return nil
}

// func (h *propertyHandler) removeWSConnection(key string) {
// 	if h == nil {
// 		log.Error().Msg("[producer:removeWSConnection] nil server")
// 	}
// 	h.mu.Lock()
// 	defer h.mu.Unlock()

// 	if _, ok := h.webSocketConnections[key]; ok {
// 		//		conn.Close() // don't close the websocket.Conn or ReadJSON returns a "use of closed network connection" error
// 		delete(h.webSocketConnections, key)
// 		h.waitWebSocket.Done()
// 	}
// }

// func (h *thingWSHandler) gracefullWSShutdown() {
// 	if h == nil {
// 		log.Error().Msg("[webSocket:gracefullWSShutdown] nil server")
// 	}
// 	h.mu.RLock()
// 	conns := h.webSocketConnections
// 	h.mu.RUnlock()

// 	for key, wsConn := range conns {
// 		log.Trace().Str("key", key).Msg("[webSocket:gracefullWSShutdown] Send Close message")
// 		err := wsConn.conn.WriteControl(websocket.CloseMessage,
// 			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
// 			time.Time{})
// 		if err != nil {
// 			log.Error().Str("key", key).Err(err).Msg("[webSocket:gracefullWSShutdown] Sending error")
// 		}
// 		h.removeWSConnection(key)
// 	}
// }
