// Minimal example: serve a Thing over HTTP and advertise it via mDNS / DNS-SD
// under the `_wot._tcp` service type.
//
// Run alongside `dns-sd -B _wot._tcp local.` (macOS) or `avahi-browse -r _wot._tcp` (Linux).
package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/protocolDnssd"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/securityScheme"
	"github.com/project-eria/go-wot/thing"
	zlog "github.com/rs/zerolog/log"
)

func main() {
	mything, err := thing.New(
		"dev:ops:dnssd-demo-1",
		"0.0.0",
		"DnssdDemo",
		"A Thing advertised via mDNS",
		[]string{},
	)
	if err != nil {
		zlog.Fatal().Err(err).Msg("[main]")
	}
	mything.AddSecurity("no_sec", securityScheme.NewNoSecurity())

	var wait sync.WaitGroup
	p := producer.New(&wait)
	p.Produce("demo", mything)

	httpServer := protocolHttp.NewServer(":8888", "", "DnssdDemo", "DnssdDemo v0.0.0")
	p.AddServer(httpServer)

	advertiser, err := protocolDnssd.NewAdvertiser(protocolDnssd.Options{
		Port:   8888,
		Scheme: protocolDnssd.SchemeHTTP,
	})
	if err != nil {
		zlog.Fatal().Err(err).Msg("[main] advertiser")
	}
	p.AddServer(advertiser)

	p.Expose()
	p.Start()

	zlog.Info().Msg("[main] running — Ctrl+C to stop")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	p.Stop()
	wait.Wait()
}
