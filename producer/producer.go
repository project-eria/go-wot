package producer

import (
	"fmt"
	"net/http"

	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

// Producer is an protocols server (http, ws, ...) made for webthings model
type Producer struct {
	ip     string
	port   int
	secure bool
	things []*ExposedThing
	//	wsHandlers []*thingWSHandler
	*http.Server
}

// New constructs the server
func New(ip string, port int, secure bool) *Producer {
	address := fmt.Sprintf(":%d", port)
	log.Info().Str("url", address).Msg("[producer:New] Server setup")

	producer := &Producer{
		ip:     ip,
		port:   port,
		secure: secure,
		things: []*ExposedThing{},
		Server: &http.Server{
			Addr: address,
		},
	}

	return producer
}

// New constructs the http server, and register the router
func (p *Producer) Produce(td *thing.Thing) *ExposedThing {
	exposedThing := NewExposedThing(td)
	host := fmt.Sprintf("%s:%d", p.ip, p.port)
	addFormHttp(exposedThing, host, p.secure)

	p.things = append(p.things, exposedThing)
	return exposedThing
}

// Produce constructs and launch an http server
func (p *Producer) Expose() {
	if p == nil {
		log.Error().Msg("[producer:Expose] nil server")
	}
	p.exposeHttp()
}
