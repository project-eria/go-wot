package protocolWebSocket

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/rs/zerolog/log"
)

var (
	propertiesObservers = map[string]map[string]map[string]*wsConnection{} // [thing][property][wsKey]
	eventSubscriptions  = map[string]map[string]map[string]*wsConnection{} // [thing][property][wsKey]
	mu                  sync.RWMutex
)

type WsServer struct {
	*protocolHttp.HttpServer
}

func NewServer(httpServer *protocolHttp.HttpServer) *WsServer {
	return &WsServer{
		HttpServer: httpServer,
	}
}

func (s *WsServer) Expose(ref string, thing *producer.ExposedThing) {
	prefix := ""
	if ref != "" {
		prefix = "/" + ref
	}
	g := s.Group(prefix)

	addEndPoints(g, s.ExposedAddr, prefix, thing)
	propertyChangeChan := thing.GetPropertyChangeChannel()
	eventChan := thing.GetEventChannel()
	go monitorPropertyObserver(propertyChangeChan)
	go monitorEvent(eventChan)
}

func addEndPoints(g fiber.Router, exposedAddr string, prefix string, t *producer.ExposedThing) {
	for _, property := range t.Td.Properties {
		if property.Observable {
			addPropertyEndPoints(g, exposedAddr, prefix, t, property)
		}
	}

	for _, event := range t.Td.Events {
		addEventEndPoints(g, exposedAddr, prefix, t, event)
	}
}

func addPropertyEndPoints(g fiber.Router, exposedAddr string, prefix string, t *producer.ExposedThing, property *interaction.Property) {
	// TODO https://w3c.github.io/wot-thing-description/#form-uriVariables
	form := &interaction.Form{
		ContentType: "application/json",
		Supplement:  map[string]interface{}{},
		Op:          []string{"observeproperty", "unobserveproperty"},
		UrlBuilder: func(host string, secure bool) string {
			protocol := "ws"
			if secure {
				protocol = "wss"
			}
			if exposedAddr != "" { // force exposed host
				host = exposedAddr
			}
			return fmt.Sprintf("%s://%s%s/%s", protocol, host, prefix, property.Key)
		},
	}
	g.Use("/"+property.Key, propertyObserverHandler(t, property))

	property.Forms = append(property.Forms, form)
	if _, in := propertiesObservers[t.Ref]; !in {
		propertiesObservers[t.Ref] = map[string]map[string]*wsConnection{}
	}
	propertiesObservers[t.Ref][property.Key] = map[string]*wsConnection{}
}

func addEventEndPoints(g fiber.Router, exposedAddr string, prefix string, t *producer.ExposedThing, event *interaction.Event) {
	// TODO https://w3c.github.io/wot-thing-description/#form-uriVariables
	form := &interaction.Form{
		ContentType: "application/json",
		Supplement:  map[string]interface{}{},
		Op:          []string{"subscribeevent"},
		UrlBuilder: func(host string, secure bool) string {
			protocol := "ws"
			if secure {
				protocol = "wss"
			}
			if exposedAddr != "" { // force exposed host
				host = exposedAddr
			}
			return fmt.Sprintf("%s://%s%s/%s", protocol, host, prefix, event.Key)
		},
	}
	g.Get("/"+event.Key, eventHandler(t, event))

	event.Forms = append(event.Forms, form)
	if _, in := eventSubscriptions[t.Ref]; !in {
		eventSubscriptions[t.Ref] = map[string]map[string]*wsConnection{}
	}
	eventSubscriptions[t.Ref][event.Key] = map[string]*wsConnection{}
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

func monitorPropertyObserver(c <-chan producer.PropertyChange) {
	for {
		propertyChange, ok := <-c
		if !ok {
			log.Trace().Str("ThingRef", propertyChange.ThingRef).Str("property", propertyChange.Name).Msg("[protocolWebSocket:monitorPropertyObserver] channel closed")
			break
		}
		if observers, ok := propertiesObservers[propertyChange.ThingRef][propertyChange.Name]; ok {
			log.Trace().Str("ThingRef", propertyChange.ThingRef).Str("property", propertyChange.Name).Msg("[protocolWebSocket:monitorPropertyObserver] Sending property change")
			for _, wsConn := range observers {
				err := wsConn.jsonWSRenderer(propertyChange.Value)
				if err != nil {
					log.Error().Err(err).Str("ThingRef", propertyChange.ThingRef).Str("property", propertyChange.Name).Msg("[protocolWebSocket:monitorPropertyObserver]")
				}
			}
		}
	}
}

func monitorEvent(c <-chan producer.Event) {
	for {
		event, ok := <-c
		if !ok {
			log.Trace().Str("ThingRef", event.ThingRef).Str("property", event.Name).Msg("[protocolWebSocket:monitorEvent] channel closed")
			break
		}
		if subscribers, ok := eventSubscriptions[event.ThingRef][event.Name]; ok {
			log.Trace().Str("ThingRef", event.ThingRef).Str("event", event.Name).Msg("[protocolWebSocket:monitorEvent] Sending event")
			for _, wsConn := range subscribers {
				err := wsConn.jsonWSRenderer(event.Value)
				if err != nil {
					log.Error().Err(err).Str("ThingRef", event.ThingRef).Str("property", event.Name).Msg("[protocolWebSocket:monitorEvent]")
				}
			}
		}
	}
}

// TODO processRxMsg processes incoming messages
// func (h *affordanceHandler) processRxMsg(wsConn *wsConnection, message *wsMessage) {
// 	log.Trace().Str("key", message.key).Str("type", message.MessageType).Msg("[producer:processRxMsg] Processing WS request")
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
