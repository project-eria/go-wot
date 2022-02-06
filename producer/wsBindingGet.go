package producer

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

func (t *ExposedThing) WSGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	name := params.ByName("name")
	log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[HTTPMiddleWare:updateWebsocket] Received Thing property WS request")

	// TODO Handle Origin for debug plugins
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, http.Header{"Sec-WebSocket-Protocol": []string{r.Header.Get("Sec-WebSocket-Protocol")}})
	if err != nil {
		log.Warn().Str("uri", r.RequestURI).Err(err).Msg("[affordanceHandler:webSocket] WebSocket Upgrade")
		return
	}
	r.Header.Get("Sec-Websocket-Key")
	key := r.Header.Get("Sec-Websocket-Key")
	wsConn := &wsConnection{Conn: conn}
	if err := t.AddWSPropertyObserver(name, key, wsConn); err != nil {
		wsConn.errorWSRenderer(err.Error())
		wsConn.Close()
		return
	}

	for {
		var data interface{}
		err := wsConn.ReadJSON(&data)
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			log.Debug().Str("key", key).Msg("[affordanceHandler:webSocket] WebSocket Normal Closure")
			t.RemoveWSPropertyObserver(name, key)
			return
		}
		if err != nil {
			log.Error().Str("key", key).Err(err).Msg("[affordanceHandler:webSocket] WebSocket error")
			t.RemoveWSPropertyObserver(name, key)
			return
		}
		log.Trace().Str("key", key).Msgf("[affordanceHandler:webSocket] Received WebSocket message: %#v", data)
		// TODO
		// h.processRxMsg(wsConn, &message)
	}
}

// TODO processRxMsg processes incoming messages
// func (h *affordanceHandler) processRxMsg(wsConn *wsConnection, message *wsMessage) {
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

// WSProcessTxMsg processes messages to be send
func (p *ExposedProperty) WSProcessTxMsg(message interface{}) {
	log.Debug().Msg("[ExposedProperty:WSProcessTxMsg] Processing WS request")
	// Send the message to all ws connections
	for _, wsConn := range p.observersProperties {
		err := wsConn.jsonWSRenderer(message)
		if err != nil {
			log.Error().Err(err).Msg("[ExposedProperty:WSProcessTxMsg] Sending propertyStatus")
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
