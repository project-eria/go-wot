package protocolWebSocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/rs/zerolog/log"
)

var (
	propertiesObservers = map[string]map[string]*wsConnection{} // TODO handle multiple things
	eventSubscriptions  = map[string]map[string]*wsConnection{} // TODO handle multiple things
	mu                  sync.RWMutex
)

type WsServer struct {
	httpServer *protocolHttp.HttpServer
}

func NewServer(httpServer *protocolHttp.HttpServer) *WsServer {
	httpServer.AddGetMiddleware(upgrade)
	return &WsServer{
		httpServer: httpServer,
	}
}

func (s *WsServer) Expose(thing *producer.ExposedThing) {
	url := fmt.Sprintf("ws://%s:%d", s.httpServer.Host, s.httpServer.Port)
	// if secure {
	// 	url = "wss://" + host
	// }
	s.addEndPoints(url, thing)
	go monitorPropertyObserver(thing.PropertyChangeChan)
	go monitorEvent(thing.EventChan)
}

func (s *WsServer) addEndPoints(base string, t *producer.ExposedThing) {
	if t == nil {
		log.Error().Msg("[protocolWebSocket:GracefullyShutdown] nil thing")
		return
	}
	var (
		href = base + "/"
	)

	for _, property := range t.Td.Properties {
		if property.Observable {
			form := interaction.Form{
				Href:        href + property.Key,
				ContentType: "application/json",
				Supplement:  map[string]interface{}{},
				Op:          []string{"observeproperty", "unobserveproperty"},
			}
			property.Forms = append(property.Forms, form)
			propertiesObservers[property.Key] = map[string]*wsConnection{}
		}
	}

	for _, event := range t.Td.Events {
		form := interaction.Form{
			Href:        href + event.Key,
			ContentType: "application/json",
			Supplement:  map[string]interface{}{},
			Op:          []string{"subscribeevent"},
		}
		event.Forms = append(event.Forms, form)
		eventSubscriptions[event.Key] = map[string]*wsConnection{}

	}
}

func (s *WsServer) Start() {

}

func (s *WsServer) Stop() {
	// TODO
	// Stop Chan monitoring routines
}

// func (t *ExposedThing) gracefullWSShutdown() {
// 	for _, p := range t.exposedProperties {
// 		p.mu.RLock()
// 		conns := p.observersProperties
// 		p.mu.RUnlock()
// 		for key, wsConn := range conns {
// 			log.Trace().Str("key", key).Msg("[ExposedProperty:gracefullWSShutdown] Send Close message")
// 			err := wsConn.WriteControl(websocket.CloseMessage,
// 				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
// 				time.Time{})
// 			if err != nil {
// 				log.Error().Str("key", key).Err(err).Msg("[ExposedProperty:gracefullWSShutdown] Sending error")
// 			}
// 			delete(p.observersProperties, key)
// 			t._wait.Done()
// 		}
// 	}
// }

func upgrade(thing *producer.ExposedThing, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if next != nil {
			if r.Header.Get("Upgrade") == "websocket" {
				// // webthing SubProtocol do not exists
				// if r.Header.Get("Sec-Websocket-Protocol") != "webthing" {
				// 	log.Error().Msg("[producer:webSocket] Connection not using webthing protocol")
				// 	w.WriteHeader(http.StatusBadRequest)
				// 	io.WriteString(w, "Connection not using webthing protocol")
				// 	return
				// }

				WSGet(w, r, p)
				return
			} else {
				next(w, r, p)
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    []string{"webthing"},
}

type wsConnection struct {
	mu sync.RWMutex
	*websocket.Conn
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

func monitorPropertyObserver(c chan producer.PropertyChange) {
	for {
		propertyChange, ok := <-c
		if !ok {
			log.Debug().Msg("[protocolWebSocket:monitorPropertyObserver] channel closed")
			break
		}
		if observers, ok := propertiesObservers[propertyChange.Name]; ok {
			log.Debug().Str("property", propertyChange.Name).Msg("[protocolWebSocket:monitorPropertyObserver] Sending property change")
			for _, wsConn := range observers {
				err := wsConn.jsonWSRenderer(propertyChange.Value)
				if err != nil {
					log.Error().Err(err).Msg("[protocolWebSocket:monitorPropertyObserver]")
				}
			}
		}
	}
}

func monitorEvent(c chan producer.Event) {
	for {
		event, ok := <-c
		if !ok {
			log.Debug().Msg("[protocolWebSocket:monitorEvent] channel closed")
			break
		}
		if subscribers, ok := eventSubscriptions[event.Name]; ok {
			log.Debug().Str("event", event.Name).Msg("[protocolWebSocket:monitorEvent] Sending event")
			for _, wsConn := range subscribers {
				err := wsConn.jsonWSRenderer(event.Value)
				if err != nil {
					log.Error().Err(err).Msg("[protocolWebSocket:monitorEvent]")
				}
			}
		}
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
