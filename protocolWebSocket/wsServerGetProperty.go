package protocolWebSocket

import (
	"fmt"

	"github.com/gofiber/websocket/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
)

func propertyObserverHandler(t *producer.ExposedThing, tdProperty *interaction.Property) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		log.Trace().Str("ThingRef", t.Ref).Str("property", tdProperty.Key).Msg("[protocolWebSocket:propertyObserverHandler] Received Thing property WS request")

		// TODO Handle Origin for debug plugins
		// upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		// conn, err := upgrader.Upgrade(w, r, http.Header{"Sec-WebSocket-Protocol": []string{r.Header.Get("Sec-WebSocket-Protocol")}})
		// if err != nil {
		// 	log.Warn().Str("uri", c.Path()).Err(err).Msg("[protocolWebSocket:propertyObserverHandler] WebSocket Upgrade")
		// 	return
		// }

		key := c.Locals("key").(string)
		wsConn := &wsConnection{Conn: c}

		if err := addPropertyObserver(t, tdProperty.Key, key, wsConn); err != nil {
			wsConn.errorWSRenderer(err.Error())
			wsConn.Close()
			return
		}

		for {
			var data interface{}
			err := wsConn.ReadJSON(&data)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Trace().Str("key", key).Msg("[protocolWebSocket:propertyObserverHandler] WebSocket Normal Closure")
				removePropertyObserver(t, tdProperty.Key, key)
				return
			}
			if err != nil {
				log.Error().Str("key", key).Err(err).Msg("[protocolWebSocket:propertyObserverHandler] WebSocket error")
				removePropertyObserver(t, tdProperty.Key, key)
				return
			}
			log.Trace().Str("key", key).Msgf("[protocolWebSocket:propertyObserverHandler] Received WebSocket message: %#v", data)
			// TODO
			// h.processRxMsg(wsConn, &message)
		}
	}
}

func addPropertyObserver(t *producer.ExposedThing, name string, key string, wsConn *wsConnection) error {
	mu.Lock()
	defer mu.Unlock()
	log.Trace().Str("ThingRef", t.Ref).Str("property", name).Str("key", key).Msg("[protocolWebSocket:addPropertyObserver] Register WS property observer connection")
	propertiesObservers[t.Ref][name][key] = wsConn
	// TODO t._wait.Add(1)
	return nil
}

func removePropertyObserver(t *producer.ExposedThing, name string, key string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := t.ExposedProperties[name]; ok {
		if t.ExposedProperties[name].Observable {
			log.Trace().Str("ThingRef", t.Ref).Str("property", name).Str("key", key).Msg("[protocolWebSocket:removePropertyObserver] Unregister WS Connection")
			if _, ok := propertiesObservers[t.Ref][name]; ok {
				//		conn.Close() // don't close the websocket.Conn or ReadJSON returns a "use of closed network connection" error
				delete(propertiesObservers[t.Ref][name], key)
				// TODO t._wait.Done()
			}
			return nil
		}
		log.Trace().Str("ThingRef", t.Ref).Str("property", name).Msg("[protocolWebSocket:removePropertyObserver] property not observable")
		return fmt.Errorf("property %s/%s not observable", t.Ref, name)
	}
	log.Trace().Str("ThingRef", t.Ref).Str("property", name).Msg("[protocolWebSocket:removePropertyObserver] property not found")
	return fmt.Errorf("property %s/%s not found", t.Ref, name)
}
