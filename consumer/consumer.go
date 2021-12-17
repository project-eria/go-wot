package consumer

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

type Consumer struct {
	things []*ConsumedThing
}

func New() *Consumer {
	consumer := &Consumer{
		things: []*ConsumedThing{},
	}

	return consumer
}

func (c *Consumer) ConsumeURL(url string) (*ConsumedThing, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error().Str("status", resp.Status).Str("url", url).Msg("[consumer:ConsumeURL] incorrect response")
		return nil, errors.New("incorrect HTTP response")
	}

	var td thing.Thing
	if err := json.NewDecoder(resp.Body).Decode(&td); err != nil {
		log.Error().Str("url", url).Err(err).Msg("[consumer:ConsumeURL]")
	}
	return c.Consume(&td), nil
}

func (c *Consumer) Consume(td *thing.Thing) *ConsumedThing {
	consumedThing := &ConsumedThing{
		td: td,
	}
	c.things = append(c.things, consumedThing)
	return consumedThing
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
