package protocolWebSocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
)

func WSGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	t := r.Context().Value("thing").(*producer.ExposedThing)
	name := params.ByName("name")
	log.Trace().Str("uri", r.RequestURI).Str("property", name).Msg("[protocolWebSocket:WSGet] Received WS request")

	// TODO Handle Origin for debug plugins
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, http.Header{"Sec-WebSocket-Protocol": []string{r.Header.Get("Sec-WebSocket-Protocol")}})
	if err != nil {
		log.Warn().Str("uri", r.RequestURI).Err(err).Msg("[protocolWebSocket:WSGet] WebSocket Upgrade")
		return
	}
	r.Header.Get("Sec-Websocket-Key")
	key := r.Header.Get("Sec-Websocket-Key")
	wsConn := &wsConnection{Conn: conn}

	if property, ok := t.ExposedProperties[name]; ok && property.Observable {
		handlePropertyObserver(t, name, key, wsConn)
	} else if _, ok := t.ExposedEvents[name]; ok {
		handleEventSubscription(t, name, key, wsConn)
	} else {
		log.Warn().Str("uri", r.RequestURI).Str("name", name).Msg("[protocolWebSocket:WSGet] No Observable Property/Even")
	}
}

func handlePropertyObserver(t *producer.ExposedThing, name string, key string, wsConn *wsConnection) {
	if err := addPropertyObserver(t, name, key, wsConn); err != nil {
		wsConn.errorWSRenderer(err.Error())
		wsConn.Close()
		return
	}

	for {
		var data interface{}
		err := wsConn.ReadJSON(&data)
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			log.Trace().Str("key", key).Msg("[protocolWebSocket:handlePropertyObserver] WebSocket Normal Closure")
			removePropertyObserver(t, name, key)
			return
		}
		if err != nil {
			log.Error().Str("key", key).Err(err).Msg("[protocolWebSocket:handlePropertyObserver] WebSocket error")
			removePropertyObserver(t, name, key)
			return
		}
		log.Trace().Str("key", key).Msgf("[protocolWebSocket:handlePropertyObserver] Received WebSocket message: %#v", data)
		// TODO
		// h.processRxMsg(wsConn, &message)
	}
}

func handleEventSubscription(t *producer.ExposedThing, name string, key string, wsConn *wsConnection) {
	if err := addEventSubscription(t, name, key, wsConn); err != nil {
		wsConn.errorWSRenderer(err.Error())
		wsConn.Close()
		return
	}

	for {
		var data interface{}
		err := wsConn.ReadJSON(&data)
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			log.Trace().Str("key", key).Msg("[protocolWebSocket:handleEventSubscription] WebSocket Normal Closure")
			removeEventSubscription(t, name, key)
			return
		}
		if err != nil {
			log.Error().Str("key", key).Err(err).Msg("[protocolWebSocket:handleEventSubscription] WebSocket error")
			removeEventSubscription(t, name, key)
			return
		}
		log.Trace().Str("key", key).Msgf("[protocolWebSocket:handleEventSubscription] Received WebSocket message: %#v", data)
	}
}

func addPropertyObserver(t *producer.ExposedThing, name string, key string, wsConn *wsConnection) error {
	mu.Lock()
	defer mu.Unlock()
	log.Trace().Str("key", key).Msg("[protocolWebSocket:addPropertyObserver] Register WS property observer connection")
	propertiesObservers[name][key] = wsConn
	// TODO t._wait.Add(1)
	return nil
}

func addEventSubscription(t *producer.ExposedThing, name string, key string, wsConn *wsConnection) error {
	mu.Lock()
	defer mu.Unlock()
	log.Trace().Str("key", key).Msg("[protocolWebSocket:addEventSubscription] Register WS event subscription connection")
	eventSubscriptions[name][key] = wsConn
	// TODO t._wait.Add(1)
	return nil
}

func removePropertyObserver(t *producer.ExposedThing, name string, key string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := t.ExposedProperties[name]; ok {
		if t.ExposedProperties[name].Observable {
			log.Trace().Str("key", key).Msg("[protocolWebSocket:removePropertyObserver] Unregister WS Connection")
			if _, ok := propertiesObservers[name]; ok {
				//		conn.Close() // don't close the websocket.Conn or ReadJSON returns a "use of closed network connection" error
				delete(propertiesObservers[name], key)
				// TODO t._wait.Done()
			}
			return nil
		}
		log.Trace().Str("property", name).Msg("[protocolWebSocket:removePropertyObserver] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Trace().Str("property", name).Msg("[protocolWebSocket:removePropertyObserver] property not found")
	return fmt.Errorf("property %s not found", name)
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
