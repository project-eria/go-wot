package producer

import (
	"sync"

	"github.com/project-eria/go-wot/thing"
	zlog "github.com/rs/zerolog/log"
)

// Producer is an protocols server (http, ws, ...) made for webthings model
type Producer struct {
	servers []ProtocolServer
	things  map[string]ExposedThing
	//	wsHandlers []*thingWSHandler
	_wait *sync.WaitGroup
	mu    sync.RWMutex
}

// New constructs the server
func New(wait *sync.WaitGroup) *Producer {
	producer := &Producer{
		things:  map[string]ExposedThing{},
		servers: []ProtocolServer{},
		_wait:   wait,
	}

	return producer
}

type ProtocolServer interface {
	Expose(string, ExposedThing)
	Start()
	Stop()
}

func (p *Producer) AddServer(server ProtocolServer) {
	if p == nil {
		zlog.Error().Msg("[producer:AddServer] nil Producer")
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.servers = append(p.servers, server)
}

// New constructs the http server, and register the router
func (p *Producer) Produce(ref string, td *thing.Thing) ExposedThing {
	if p == nil {
		zlog.Error().Msg("[producer:Produce] nil Producer")
		return nil
	}
	exposedThing := NewExposedThing(td, ref, p._wait)

	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.things[ref]; exists {
		zlog.Error().Msg("[producer:Produce] thing ref already exists")
		return nil
	}
	p.things[ref] = exposedThing
	return exposedThing
}

// Produce constructs
func (p *Producer) Expose() {
	if p == nil {
		zlog.Error().Msg("[producer:Expose] nil Producer")
		return
	}

	if len(p.servers) == 0 {
		zlog.Fatal().Msg("[producer:Expose] no servers to expose Things")
		return
	}
	for _, s := range p.servers {
		for ref, t := range p.things {
			s.Expose(ref, t)
		}
	}
}

// Launch servers
func (p *Producer) Start() {
	if p == nil {
		zlog.Error().Msg("[producer:Start] nil Producer")
		return
	}
	if len(p.servers) == 0 {
		zlog.Fatal().Msg("[producer:Start] no servers to start")
		return
	}
	zlog.Info().Msg("[producer:Start] Starting...")
	for _, s := range p.servers {
		s.Start()
	}
}

func (p *Producer) Stop() {
	if p == nil {
		zlog.Error().Msg("[producer:Stop] nil server")
		return
	}
	zlog.Info().Msg("[producer:Stop] Stopping...")
	for _, s := range p.servers {
		s.Stop()
	}
}
