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

func propertyObserverHandler(t producer.ExposedThing, tdProperty *interaction.Property) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		property, err := t.ExposedProperty(tdProperty.Key)
		if err != nil {
			zlog.Error().Str("property", tdProperty.Key).Msg("[protocolWebSocket:propertyObserverHandler] ExposedProperty not found")
			return c.Status(protocolHttp.UnknownError.HttpStatus).JSON(fiber.Map{
				"error": fmt.Sprintf("ExposedProperty `%s` not found", tdProperty.Key),
				"type":  protocolHttp.UnknownError.ErrorType,
			})
		} else {
			if websocket.IsWebSocketUpgrade(c) {
				key := c.Get("Sec-Websocket-Key")
				c.Locals("key", key)

				// Check the params (uriVariables) data
				parametersStr := c.AllParams()
				parameters, err := property.CheckUriVariables(parametersStr)
				if err != nil {
					return c.Status(protocolHttp.DataError.HttpStatus).JSON(fiber.Map{
						"error": err.Error(),
						"type":  protocolHttp.DataError.ErrorType,
					})
				}

				return websocket.New(propertyObserverWSHandler(t, tdProperty, parameters))(c)
			}
			return c.Next()
		}
	}
}

func propertyObserverWSHandler(t producer.ExposedThing, tdProperty *interaction.Property, parameters map[string]interface{}) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		zlog.Trace().Str("ThingRef", t.Ref()).Str("property", tdProperty.Key).Interface("parameters", parameters).Msg("[protocolWebSocket:propertyObserverHandler] Received Thing property WS request")

		// TODO Handle Origin for debug plugins
		// upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		// conn, err := upgrader.Upgrade(w, r, http.Header{"Sec-WebSocket-Protocol": []string{r.Header.Get("Sec-WebSocket-Protocol")}})
		// if err != nil {
		// 	zlog.Warn().Str("uri", c.Path()).Err(err).Msg("[protocolWebSocket:propertyObserverHandler] WebSocket Upgrade")
		// 	return
		// }

		key := c.Locals("key").(string)

		// Deep clone the listener parameters
		parametersCopy := make(map[string]interface{})
		parametersJSON, _ := json.Marshal(parameters)   // Marshalling the map to JSON
		json.Unmarshal(parametersJSON, &parametersCopy) // Unmarshalling JSON to a new map (Deep Copy)

		wsConn := &wsConnection{Conn: c, listenerParameters: parametersCopy}
		if err := addPropertyObserver(t, tdProperty.Key, key, wsConn); err != nil {
			wsConn.errorWSRenderer(err.Error())
			wsConn.Close()
			return
		}

		for {
			var data interface{}
			err := wsConn.ReadJSON(&data)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				zlog.Trace().Str("key", key).Msg("[protocolWebSocket:propertyObserverHandler] WebSocket Normal Closure")
				removePropertyObserver(t, tdProperty.Key, key)
				return
			}
			if err != nil {
				zlog.Error().Str("key", key).Err(err).Msg("[protocolWebSocket:propertyObserverHandler] WebSocket error")
				removePropertyObserver(t, tdProperty.Key, key)
				return
			}
			zlog.Trace().Str("key", key).Msgf("[protocolWebSocket:propertyObserverHandler] Received WebSocket message: %#v", data)
			// TODO
			// h.processRxMsg(wsConn, &message)
		}

	}
}

func addPropertyObserver(t producer.ExposedThing, name string, key string, wsConn *wsConnection) error {
	mu.Lock()
	defer mu.Unlock()
	zlog.Trace().Str("ThingRef", t.Ref()).Str("property", name).Str("key", key).Msg("[protocolWebSocket:addPropertyObserver] Register WS property observer connection")
	propertiesObservers[t.Ref()][name][key] = wsConn
	// TODO t._wait.Add(1)
	return nil
}

func removePropertyObserver(t producer.ExposedThing, name string, key string) error {
	mu.Lock()
	defer mu.Unlock()
	property, err := t.ExposedProperty(name)
	if err != nil {
		zlog.Trace().Str("ThingRef", t.Ref()).Str("property", name).Msg("[protocolWebSocket:removePropertyObserver] property not found")
		return fmt.Errorf("property %s/%s not found", t.Ref(), name)
	} else {
		if property.IsObservable() {
			zlog.Trace().Str("ThingRef", t.Ref()).Str("property", name).Str("key", key).Msg("[protocolWebSocket:removePropertyObserver] Unregister WS Connection")
			if _, ok := propertiesObservers[t.Ref()][name]; ok {
				// conn.Close() // don't close the websocket.Conn or ReadJSON returns a "use of closed network connection" error
				delete(propertiesObservers[t.Ref()][name], key)
				// TODO t._wait.Done()
			}
			return nil
		}
		zlog.Trace().Str("ThingRef", t.Ref()).Str("property", name).Msg("[protocolWebSocket:removePropertyObserver] property not observable")
		return fmt.Errorf("property %s/%s not observable", t.Ref(), name)
	}
}
