package protocolWebSocket

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/protocolHttp"
	zlog "github.com/rs/zerolog/log"
)

func eventHandler(t producer.ExposedThing, tdEvent *interaction.Event) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		event, err := t.ExposedEvent(tdEvent.Key)
		if err != nil {
			zlog.Error().Str("event", tdEvent.Key).Msg("[protocolWebSocket:eventHandler] ExposedEvent not found")
			return c.Status(protocolHttp.UnknownError.HttpStatus).JSON(fiber.Map{
				"error": fmt.Sprintf("ExposedEvent `%s` not found", tdEvent.Key),
				"type":  protocolHttp.UnknownError.ErrorType,
			})
		} else {
			if websocket.IsWebSocketUpgrade(c) {
				key := c.Get("Sec-Websocket-Key")
				c.Locals("key", key)

				// Check the params (uriVariables) data
				parametersStr := c.AllParams()
				parameters, err := event.CheckUriVariables(parametersStr)
				if err != nil {
					return c.Status(protocolHttp.DataError.HttpStatus).JSON(fiber.Map{
						"error": err.Error(),
						"type":  protocolHttp.DataError.ErrorType,
					})
				}

				return websocket.New(eventWSHandler(t, tdEvent, parameters))(c)
			}
			return c.Next()
		}
	}
}

func eventWSHandler(t producer.ExposedThing, tdEvent *interaction.Event, parameters map[string]interface{}) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		zlog.Trace().Str("event", tdEvent.Key).Interface("parameters", parameters).Msg("[protocolWebSocket:propertyEventHandler] Received Thing event WS request")

		// TODO Handle Origin for debug plugins
		// upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		// conn, err := upgrader.Upgrade(w, r, http.Header{"Sec-WebSocket-Protocol": []string{r.Header.Get("Sec-WebSocket-Protocol")}})
		// if err != nil {
		// 	zlog.Warn().Str("uri", c.Path()).Err(err).Msg("[protocolWebSocket:WSGet] WebSocket Upgrade")
		// 	return
		// }
		key := c.Locals("key").(string)

		// Deep clone the listener parameters
		parametersCopy := make(map[string]interface{})
		parametersJSON, _ := json.Marshal(parameters)   // Marshalling the map to JSON
		json.Unmarshal(parametersJSON, &parametersCopy) // Unmarshalling JSON to a new map (Deep Copy)

		wsConn := &wsConnection{Conn: c, listenerParameters: parametersCopy}

		if err := addEventSubscription(t, tdEvent.Key, key, wsConn); err != nil {
			wsConn.errorWSRenderer(err.Error())
			wsConn.Close()
			return
		}

		for {
			var data interface{}
			err := wsConn.ReadJSON(&data)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				zlog.Trace().Str("key", key).Msg("[protocolWebSocket:handleEventSubscription] WebSocket Normal Closure")
				removeEventSubscription(t, tdEvent.Key, key)
				return
			}
			if err != nil {
				zlog.Error().Str("key", key).Err(err).Msg("[protocolWebSocket:handleEventSubscription] WebSocket error")
				removeEventSubscription(t, tdEvent.Key, key)
				return
			}
			zlog.Trace().Str("key", key).Msgf("[protocolWebSocket:handleEventSubscription] Received WebSocket message: %#v", data)
		}
	}
}

func addEventSubscription(t producer.ExposedThing, name string, key string, wsConn *wsConnection) error {
	mu.Lock()
	defer mu.Unlock()
	zlog.Trace().Str("ThingRef", t.Ref()).Str("event", name).Str("key", key).Msg("[protocolWebSocket:addEventSubscription] Register WS event subscription connection")
	eventSubscriptions[t.Ref()][name][key] = wsConn
	// TODO t._wait.Add(1)
	return nil
}

func removeEventSubscription(t producer.ExposedThing, name string, key string) error {
	mu.Lock()
	defer mu.Unlock()
	_, err := t.ExposedEvent(name)
	if err != nil {
		zlog.Trace().Str("ThingRef", t.Ref()).Str("event", name).Msg("[protocolWebSocket:removePropertyObserver] event not found")
		return fmt.Errorf("event %s/%s not found", t.Ref(), name)
	} else {
		zlog.Trace().Str("ThingRef", t.Ref()).Str("event", name).Str("key", key).Msg("[protocolWebSocket:removePropertyObserver] Unregister WS Connection")
		if _, ok := eventSubscriptions[t.Ref()][name]; ok {
			//		conn.Close() // don't close the websocket.Conn or ReadJSON returns a "use of closed network connection" error
			delete(eventSubscriptions[t.Ref()][name], key)
			// TODO t._wait.Done()
		}
		return nil
	}
}
