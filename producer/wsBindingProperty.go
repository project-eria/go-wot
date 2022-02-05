package producer

import (
	"fmt"
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
	mu sync.RWMutex
	*websocket.Conn
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
	wsConn := &wsConnection{Conn: conn}
	if err := h.AddWSPropertyObserver(name, key, wsConn); err != nil {
		wsConn.errorWSRenderer(err.Error())
		wsConn.Close()
		return
	}

	for {
		var data interface{}
		err := wsConn.ReadJSON(&data)
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			log.Debug().Str("key", key).Msg("[producer:webSocket] WebSocket Normal Closure")
			h.RemoveWSPropertyObserver(name, key)
			return
		}
		if err != nil {
			log.Error().Str("key", key).Err(err).Msg("[producer:webSocket] WebSocket error")
			h.RemoveWSPropertyObserver(name, key)
			return
		}
		log.Trace().Str("key", key).Msgf("[producer:webSocket] Received WebSocket message: %#v", data)
		// TODO
		// h.processRxMsg(wsConn, &message)
	}
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

func (t *ExposedThing) AddWSPropertyObserver(name string, key string, wsConn *wsConnection) error {
	if _, ok := t.Td.Properties[name]; ok {
		if t.Td.Properties[name].Observable {
			p := t.exposedProperties[name]
			log.Debug().Str("key", key).Msg("[ExposedThing:AddWSPropertyObserver] Register WS Connection")
			p.mu.Lock()
			defer p.mu.Unlock()
			p.observersProperties[key] = wsConn
			t._wait.Add(1)
			return nil
		}
		log.Debug().Str("property", name).Msg("[ExposedThing:AddWSPropertyObserver] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Debug().Str("property", name).Msg("[ExposedThing:AddWSPropertyObserver] property not found")
	return fmt.Errorf("property %s not found", name)
}

func (t *ExposedThing) RemoveWSPropertyObserver(name string, key string) error {
	if _, ok := t.Td.Properties[name]; ok {
		if t.Td.Properties[name].Observable {
			p := t.exposedProperties[name]
			log.Debug().Str("key", key).Msg("[ExposedThing:RemoveWSPropertyObserver] Unregister WS Connection")
			p.mu.Lock()
			defer p.mu.Unlock()

			if _, ok := p.observersProperties[key]; ok {
				//		conn.Close() // don't close the websocket.Conn or ReadJSON returns a "use of closed network connection" error
				delete(p.observersProperties, key)
				t._wait.Done()
			}
			return nil
		}
		log.Debug().Str("property", name).Msg("[ExposedThing:RemoveWSPropertyObserver] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Debug().Str("property", name).Msg("[ExposedThing:RemoveWSPropertyObserver] property not found")
	return fmt.Errorf("property %s not found", name)
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
