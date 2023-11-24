package consumer

import (
	"net/url"
	"sync"

	"github.com/project-eria/go-wot/interaction"

	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

type Consumer struct {
	clients map[string]ProtocolClient
	things  []ConsumedThing
	mu      sync.RWMutex
}

func New() *Consumer {
	consumer := &Consumer{
		clients: map[string]ProtocolClient{},
		things:  []ConsumedThing{},
	}

	return consumer
}

type ProtocolClient interface {
	GetSchemes() []string
	ReadResource(*interaction.Form) (interface{}, error)
	WriteResource(*interaction.Form, interface{}) (interface{}, error)
	InvokeResource(*interaction.Form, interface{}) (interface{}, error)
	SubscribeResource(*interaction.Form, *Subscription, Listener) error
	Stop()
}

func (c *Consumer) Consume(td *thing.Thing) ConsumedThing {
	ct := &consumedThing{
		consumer: c,
		td:       td,
	}
	c.things = append(c.things, ct)
	return ct
}

func (c *Consumer) GetClientFor(form *interaction.Form) ProtocolClient {
	u, err := url.Parse(form.Href)
	if err != nil {
		log.Error().Str("href", form.Href).Err(err).Msg("[consumer:getClientFor] href not readable")
		return nil
	}
	if client, found := c.clients[u.Scheme]; found {
		log.Trace().Str("scheme", u.Scheme).Msg("[consumer:getClientFor] got client for scheme")
		return client
	}
	log.Error().Str("scheme", u.Scheme).Msg("[consumer:getClientFor] missing client for scheme")
	return nil
}

func (c *Consumer) AddClient(client ProtocolClient) {
	if c == nil {
		log.Error().Msg("[consumer:AddClient] nil Consumer")
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	schemes := client.GetSchemes()
	for _, scheme := range schemes {
		c.clients[scheme] = client
	}
}

func (c *Consumer) Shutdown() {
	for _, client := range c.clients {
		client.Stop()
	}
}

// // Consumer structure to handle the connection and description information
// type ConsumedThing struct {
// 	url           url.URL
// 	urn           string
// 	title         string
// 	subscriptions []subscription
// 	mu            sync.RWMutex
// 	isConnected   bool
// 	wsRequired    bool
// 	dialErr       error
// 	//	connWait      connWait
// 	*websocket.Conn
// }

// // SingleHandler function handler type for single property
// type SingleHandler func(string, interface{}, ...interface{})

// // AllHandler function handler type for all properties
// type AllHandler func(map[string]interface{}, ...interface{})

// type subscription struct {
// 	property      string
// 	handlerSingle SingleHandler
// 	handlerAll    AllHandler
// 	context       []interface{}
// }

// // New a remote thing, using WebSocket
// func Consume(url url.URL, wsRequired bool) (*ConsumedThing, error) {
// 	url.Scheme = "http"

// 	data, err := getHTTPJSON(url.String())
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Info().Str("url", url.String()).Msg("[consumer:New] Successfully got thing description")
// 	log.Trace().Interface("data", data).Msg("Description raw data")

// 	if _, ok := data["id"]; !ok {
// 		return nil, errors.New("Incorrect JSON thing properties")
// 	}

// 	conn := &Consumer{
// 		urn:           data["id"].(string),
// 		title:         data["title"].(string),
// 		subscriptions: []subscription{},
// 		url:           url,
// 		wsRequired:    wsRequired,
// 		connWait:      newConnWait(),
// 	}

// 	return conn, nil
// }

// // SubscribeSingle to a specific property
// // call the handler when a change event is received
// func (t *Consumer) SubscribeSingle(property string, handler SingleHandler, context ...interface{}) {
// 	if t == nil {
// 		log.Error().Msg("[consumer:SubscribeSingle] nil connection")
// 		return
// 	}
// 	t.mu.Lock()
// 	defer t.mu.Unlock()

// 	if !t.wsRequired {
// 		log.Error().Msg("[consumer:SubscribeSingle] Subscription require WebSocket connection")
// 		return
// 	}
// 	t.subscriptions = append(t.subscriptions, subscription{property: property, handlerSingle: handler, context: context})
// }

// // SubscribeAll to a specific property
// // call the handler when a change event is received
// func (t *Consumer) SubscribeAll(handler AllHandler, context ...interface{}) {
// 	if t == nil {
// 		log.Error().Msg("[consumer:SubscribeAll] Subscription require WebSocket connection")
// 		return
// 	}
// 	t.mu.Lock()
// 	defer t.mu.Unlock()

// 	if !t.wsRequired {
// 		log.Error().Msg("[consumer:SubscribeAll] the WebSocket is not connected")
// 		return
// 	}
// 	t.subscriptions = append(t.subscriptions, subscription{handlerAll: handler, context: context})
// }

// // getWSURL returns current WebSocket connection url
// func (t *Consumer) getWSURL() string {
// 	if t == nil {
// 		log.Error().Msg("[consumer:getWSURL] nil connection")
// 		return ""
// 	}
// 	t.mu.RLock()
// 	defer t.mu.RUnlock()

// 	wsURL := t.url
// 	wsURL.Scheme = "ws"

// 	return wsURL.String()
// }

// // getHTTPURL returns current WebSocket connection url
// func (t *Consumer) getHTTPURL(subpath string) string {
// 	if t == nil {
// 		log.Error().Msg("[consumer:getWSURL] nil connection")
// 		return ""
// 	}
// 	t.mu.RLock()
// 	defer t.mu.RUnlock()

// 	httpURL := t.url
// 	httpURL.Scheme = "http"
// 	httpURL.Path = httpURL.Path + subpath
// 	return httpURL.String()
// }
