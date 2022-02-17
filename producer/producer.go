package producer

import (
	"sync"

	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

// Producer is an protocols server (http, ws, ...) made for webthings model
type Producer struct {
	servers []ProtocolServer
	things  []*ExposedThing
	//	wsHandlers []*thingWSHandler
	_wait *sync.WaitGroup
	mu    sync.RWMutex
}

// New constructs the server
func New(wait *sync.WaitGroup) *Producer {
	producer := &Producer{
		things:  []*ExposedThing{},
		servers: []ProtocolServer{},
		_wait:   wait,
	}

	return producer
}

type ProtocolServer interface {
	Expose(*ExposedThing)
	Start()
	Stop()
}

func (p *Producer) AddServer(server ProtocolServer) {
	if p == nil {
		log.Error().Msg("[producer:AddServer] nil Producer")
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.servers = append(p.servers, server)
}

// New constructs the http server, and register the router
func (p *Producer) Produce(td *thing.Thing) *ExposedThing {
	if p == nil {
		log.Error().Msg("[producer:Produce] nil Producer")
		return nil
	}
	exposedThing := NewExposedThing(td, p._wait)

	// addFormHttp(exposedThing, host, p.secure)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.things = append(p.things, exposedThing)
	return exposedThing
}

// Produce constructs and launch an http server
func (p *Producer) Expose() {
	if p == nil {
		log.Error().Msg("[producer:Expose] nil Producer")
		return
	}
	log.Info().Msg("[producer:Expose] Starting...")

	if len(p.servers) == 0 {
		log.Fatal().Msg("[producer:Expose] no servers to expose Things")
		return
	}
	for _, s := range p.servers {
		for _, t := range p.things {
			s.Expose(t)
		}
		s.Start()
	}
}

func (p *Producer) Stop() {
	if p == nil {
		log.Error().Msg("[producer:Stop] nil server")
		return
	}
	log.Info().Msg("[producer:Stop] Stopping...")
	for _, s := range p.servers {
		s.Stop()
	}
}
