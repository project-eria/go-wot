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
	log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[protocolWebSocket:WSGet] Received WS request")

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
	if err := addPropertyObserver(t, name, key, wsConn); err != nil {
		wsConn.errorWSRenderer(err.Error())
		wsConn.Close()
		return
	}

	for {
		var data interface{}
		err := wsConn.ReadJSON(&data)
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			log.Debug().Str("key", key).Msg("[protocolWebSocket:WSGet] WebSocket Normal Closure")
			removePropertyObserver(t, name, key)
			return
		}
		if err != nil {
			log.Error().Str("key", key).Err(err).Msg("[protocolWebSocket:WSGet] WebSocket error")
			removePropertyObserver(t, name, key)
			return
		}
		log.Trace().Str("key", key).Msgf("[protocolWebSocket:WSGet] Received WebSocket message: %#v", data)
		// TODO
		// h.processRxMsg(wsConn, &message)
	}
}

func addPropertyObserver(t *producer.ExposedThing, name string, key string, wsConn *wsConnection) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := t.Td.Properties[name]; ok {
		if t.Td.Properties[name].Observable {
			log.Debug().Str("key", key).Msg("[protocolWebSocket:addWSObserver] Register WS property observer connection")
			propertiesObservers[name][key] = wsConn
			// TODO t._wait.Add(1)
			return nil
		}
		log.Debug().Str("property", name).Msg("[protocolWebSocket:addWSObserver] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Debug().Str("property", name).Msg("[protocolWebSocket:addWSObserver] property not found")
	return fmt.Errorf("property %s not found", name)
}

func removePropertyObserver(t *producer.ExposedThing, name string, key string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := t.Td.Properties[name]; ok {
		if t.Td.Properties[name].Observable {
			log.Debug().Str("key", key).Msg("[ExposedThing:RemoveWSPropertyObserver] Unregister WS Connection")
			if _, ok := propertiesObservers[name]; ok {
				//		conn.Close() // don't close the websocket.Conn or ReadJSON returns a "use of closed network connection" error
				delete(propertiesObservers[name], key)
				// TODO t._wait.Done()
			}
			return nil
		}
		log.Debug().Str("property", name).Msg("[ExposedThing:RemoveWSPropertyObserver] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Debug().Str("property", name).Msg("[ExposedThing:RemoveWSPropertyObserver] property not found")
	return fmt.Errorf("property %s not found", name)
}
