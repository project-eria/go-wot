package protocolWebSocket

import (
	"fmt"

	"github.com/gofiber/websocket/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
)

func eventHandler(t *producer.ExposedThing, tdEvent *interaction.Event) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		log.Trace().Str("event", tdEvent.Key).Msg("[protocolWebSocket:propertyEventHandler] Received Thing event WS request")

		// TODO Handle Origin for debug plugins
		// upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		// conn, err := upgrader.Upgrade(w, r, http.Header{"Sec-WebSocket-Protocol": []string{r.Header.Get("Sec-WebSocket-Protocol")}})
		// if err != nil {
		// 	log.Warn().Str("uri", c.Path()).Err(err).Msg("[protocolWebSocket:WSGet] WebSocket Upgrade")
		// 	return
		// }
		key := c.Locals("key").(string)
		wsConn := &wsConnection{Conn: c}

		if err := addEventSubscription(t, tdEvent.Key, key, wsConn); err != nil {
			wsConn.errorWSRenderer(err.Error())
			wsConn.Close()
			return
		}

		for {
			var data interface{}
			err := wsConn.ReadJSON(&data)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Trace().Str("key", key).Msg("[protocolWebSocket:handleEventSubscription] WebSocket Normal Closure")
				removeEventSubscription(t, tdEvent.Key, key)
				return
			}
			if err != nil {
				log.Error().Str("key", key).Err(err).Msg("[protocolWebSocket:handleEventSubscription] WebSocket error")
				removeEventSubscription(t, tdEvent.Key, key)
				return
			}
			log.Trace().Str("key", key).Msgf("[protocolWebSocket:handleEventSubscription] Received WebSocket message: %#v", data)
		}
	}
}

func addEventSubscription(t *producer.ExposedThing, name string, key string, wsConn *wsConnection) error {
	mu.Lock()
	defer mu.Unlock()
	log.Trace().Str("key", key).Msg("[protocolWebSocket:addEventSubscription] Register WS event subscription connection")
	eventSubscriptions[name][key] = wsConn
	// TODO t._wait.Add(1)
	return nil
}

func removeEventSubscription(t *producer.ExposedThing, name string, key string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := t.ExposedEvents[name]; ok {
		log.Trace().Str("key", key).Msg("[protocolWebSocket:removePropertyObserver] Unregister WS Connection")
		if _, ok := eventSubscriptions[name]; ok {
			//		conn.Close() // don't close the websocket.Conn or ReadJSON returns a "use of closed network connection" error
			delete(eventSubscriptions[name], key)
			// TODO t._wait.Done()
		}
		return nil
	}
	log.Trace().Str("event", name).Msg("[protocolWebSocket:removePropertyObserver] event not found")
	return fmt.Errorf("event %s not found", name)
}
