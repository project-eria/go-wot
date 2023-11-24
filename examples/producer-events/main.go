package main

import (
	"os"
	"sync"
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

	mything, err := thing.New(
		"dev:ops:my-actuator-1234",
		"0.0.0-dev",
		"Actuator1 Example",
		"An actuator 1st example",
		[]string{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("[main]")
	}
	// Add Security
	noSecurityScheme := securityScheme.NewNoSecurity()
	mything.AddSecurity("no_sec", noSecurityScheme)

	booleanData := dataSchema.NewBoolean(false)
	propertyRWO := interaction.NewProperty(
		"boolRWO",
		"RWO bool",
		"Readable/Writable/Observable boolean",
		false,
		false,
		true,
		nil,
		booleanData,
	)
	mything.AddProperty(propertyRWO)

	stringEvent := dataSchema.NewString("", 0, 0, "")
	dEvent := interaction.NewEvent("d", "d Event", "", &stringEvent)
	mything.AddEvent(dEvent)

	// Run Server
	var wait sync.WaitGroup
	myProducer := producer.New(&wait)
	exposedThing := myProducer.Produce("", mything)
	exposedThing.SetEventHandler("d", handlerD)
	httpServer := protocolHttp.NewServer(":8888", "", "My App", "My App v0.0.0")
	myProducer.AddServer(httpServer)
	wsServer := protocolWebSocket.NewServer(httpServer)
	myProducer.AddServer(wsServer)
	myProducer.Expose()

	for {
		time.Sleep(10 * time.Second)
		exposedThing.EmitPropertyChange("boolRWO", nil, nil)
		exposedThing.EmitEvent("d", nil)
	}
}

func handlerD() (interface{}, error) {
	return nil, nil
}
