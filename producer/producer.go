package producer

import (
	"fmt"
	"net/http"
	"sync"

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
	_wait *sync.WaitGroup
	*http.Server
}

// New constructs the server
func New(ip string, port int, secure bool, wait *sync.WaitGroup) *Producer {
	address := fmt.Sprintf(":%d", port)
	log.Info().Str("url", address).Msg("[producer:New] Server setup")

	producer := &Producer{
		ip:     ip,
		port:   port,
		secure: secure,
		things: []*ExposedThing{},
		_wait:  wait,
		Server: &http.Server{
			Addr: address,
		},
	}

	return producer
}

// New constructs the http server, and register the router
func (p *Producer) Produce(td *thing.Thing) *ExposedThing {
	exposedThing := NewExposedThing(td, p._wait)
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
	log.Info().Msg("[producer:Expose] Starting...")

	for _, t := range p.things {
		t.Expose()
	}

	p.exposeHttp()
}

func (p *Producer) Stop() {
	if p == nil {
		log.Error().Msg("[producer:Stop] nil server")
	}
	log.Info().Msg("[producer:Stop] Stopping...")
	for _, t := range p.things {
		t.gracefullWSShutdown()
		t.Destroy()
	}
}
