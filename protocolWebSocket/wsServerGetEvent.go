package protocolWebSocket

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/rs/zerolog/log"
)

func eventHandler(t *producer.ExposedThing, tdEvent *interaction.Event) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if event, ok := t.ExposedEvents[tdEvent.Key]; ok {
			if websocket.IsWebSocketUpgrade(c) {
				key := c.Get("Sec-Websocket-Key")
				c.Locals("key", key)

				// Check the params (uriVariables) data
				options := c.AllParams()
				if err := event.CheckUriVariables(options); err != nil {
					return c.Status(protocolHttp.DataError.HttpStatus).JSON(fiber.Map{
						"error": err.Error(),
						"type":  protocolHttp.DataError.ErrorType,
					})
				}

				return websocket.New(eventWSHandler(t, tdEvent, options))(c)
			}
			return c.Next()
		} else {
			log.Error().Str("event", tdEvent.Key).Msg("[protocolWebSocket:eventHandler] ExposedEvent not found")
			return c.Status(protocolHttp.UnknownError.HttpStatus).JSON(fiber.Map{
				"error": fmt.Errorf("ExposedEvent `%s` not found", tdEvent.Key),
				"type":  protocolHttp.UnknownError.ErrorType,
			})
		}
	}
}

func eventWSHandler(t *producer.ExposedThing, tdEvent *interaction.Event, options map[string]string) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		log.Trace().Str("event", tdEvent.Key).Interface("options", options).Msg("[protocolWebSocket:propertyEventHandler] Received Thing event WS request")

		// TODO Handle Origin for debug plugins
		// upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		// conn, err := upgrader.Upgrade(w, r, http.Header{"Sec-WebSocket-Protocol": []string{r.Header.Get("Sec-WebSocket-Protocol")}})
		// if err != nil {
		// 	log.Warn().Str("uri", c.Path()).Err(err).Msg("[protocolWebSocket:WSGet] WebSocket Upgrade")
		// 	return
		// }
		key := c.Locals("key").(string)
		// Deep clone the options
		optionsCopy := make(map[string]string)
		for k, v := range options {
			optionsCopy[k] = strings.Clone(v)
		}
		wsConn := &wsConnection{Conn: c, options: optionsCopy}

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
	log.Trace().Str("ThingRef", t.Ref).Str("event", name).Str("key", key).Msg("[protocolWebSocket:addEventSubscription] Register WS event subscription connection")
	eventSubscriptions[t.Ref][name][key] = wsConn
	// TODO t._wait.Add(1)
	return nil
}

func removeEventSubscription(t *producer.ExposedThing, name string, key string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := t.ExposedEvents[name]; ok {
		log.Trace().Str("ThingRef", t.Ref).Str("event", name).Str("key", key).Msg("[protocolWebSocket:removePropertyObserver] Unregister WS Connection")
		if _, ok := eventSubscriptions[t.Ref][name]; ok {
			//		conn.Close() // don't close the websocket.Conn or ReadJSON returns a "use of closed network connection" error
			delete(eventSubscriptions[t.Ref][name], key)
			// TODO t._wait.Done()
		}
		return nil
	}
	log.Trace().Str("ThingRef", t.Ref).Str("event", name).Msg("[protocolWebSocket:removePropertyObserver] event not found")
	return fmt.Errorf("event %s/%s not found", t.Ref, name)
}
