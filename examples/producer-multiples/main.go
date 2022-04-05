package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/project-eria/go-wot/protocolWebSocket"

	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/protocolHttp"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/securityScheme"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "02/01|15:04:05"})
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local)
	}
	defer func() {
		log.Info().Msg("[main] Stopped")
	}()

	// THING 1
	mything1, err := thing.New(
		"dev:ops:my-actuator-1234",
		"v0.0.0",
		"Actuator1 Example",
		"An actuator 1st example",
		[]string{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("[main]")
	}
	// Add Security
	noSecurityScheme1 := securityScheme.NewNoSecurity()
	mything1.AddSecurity("no_sec", noSecurityScheme1)

	// Properties
	booleanData1 := dataSchema.NewBoolean(false)
	propertyRWO1 := interaction.NewProperty(
		"boolRWO",
		"RWO bool",
		"Readable/Writable/Observable boolean",
		false,
		false,
		true,
		booleanData1,
	)
	mything1.AddProperty(propertyRWO1)

	// THING 2
	mything2, err := thing.New(
		"dev:ops:my-actuator-5678",
		"v0.0.0",
		"Actuator2 Example",
		"An actuator 2nd example",
		[]string{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("[main]")
	}
	// Add Security
	noSecurityScheme2 := securityScheme.NewNoSecurity()
	mything2.AddSecurity("no_sec", noSecurityScheme2)

	// Properties
	booleanData2 := dataSchema.NewBoolean(false)
	propertyRWO2 := interaction.NewProperty(
		"boolRWO",
		"RWO bool",
		"Readable/Writable/Observable boolean",
		false,
		false,
		true,
		booleanData2,
	)
	mything2.AddProperty(propertyRWO2)

	// Run Server
	var wait sync.WaitGroup
	myProducer := producer.New(&wait)
	myProducer.Produce("mything1", mything1)
	myProducer.Produce("mything2", mything2)
	httpServer := protocolHttp.NewServer(":8888", "")
	myProducer.AddServer(httpServer)
	wsServer := protocolWebSocket.NewServer(httpServer)
	myProducer.AddServer(wsServer)
	myProducer.Expose()

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	// Block until keyboard interrupt is received.
	<-c

	myProducer.Stop()
	wait.Wait()
}
