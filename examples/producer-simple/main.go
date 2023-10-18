package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/protocolWebSocket"

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
	mything.AddContext("schema", "https://schema.org/")
	mything.AddVersion("schema:softwareVersion", "1.1.1")

	// Add Security
	noSecurityScheme := securityScheme.NewNoSecurity()
	mything.AddSecurity("no_sec", noSecurityScheme)

	// Properties
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

	// 	propertyRW := interaction.NewProperty(
	// 		"boolRW",
	// 		"RW bool",
	// 		"Readable/Writable/Not Observable boolean",
	// 		false,
	// 		false,
	// 		false,
	// 		booleanData,
	// 	)
	// 	mything.AddProperty(propertyRW)

	// 	propertyR := interaction.NewProperty(
	// 		"boolR",
	// 		"R bool",
	// 		"Readable only/Not Observable boolean",
	// 		true,
	// 		false,
	// 		false,
	// 		booleanData,
	// 	)
	// 	mything.AddProperty(propertyR)

	// 	propertyW := interaction.NewProperty(
	// 		"boolW",
	// 		"W bool",
	// 		"Writable only/Not Observable boolean",
	// 		false,
	// 		true,
	// 		false,
	// 		booleanData,
	// 	)
	// 	mything.AddProperty(propertyW)

	// 	aAction := interaction.NewAction(
	// 		"a",
	// 		"No Input, No Output",
	// 		"",
	// 		nil,
	// 		nil,
	// 	)
	// 	mything.AddAction(aAction)

	// 	stringInput := dataSchema.NewString("", 0, 0, "")
	// 	bAction := interaction.NewAction(
	// 		"b",
	// 		"String Input, No Output",
	// 		"",
	// 		&stringInput,
	// 		nil,
	// 	)
	// 	mything.AddAction(bAction)
	// 	stringOutput := dataSchema.NewString("", 0, 0, "")
	// 	cAction := interaction.NewAction(
	// 		"c",
	// 		"String Input, String Output",
	// 		"",
	// 		&stringInput,
	// 		&stringOutput,
	// 	)
	// 	mything.AddAction(cAction)

	// Run Server
	var wait sync.WaitGroup
	myProducer := producer.New(&wait)
	//exposedThing :=
	myProducer.Produce("", mything)
	// 	exposedThing.SetActionHandler("a", handlerA)
	// 	exposedThing.SetActionHandler("b", handlerB)
	// 	exposedThing.SetActionHandler("c", handlerC)
	httpServer := protocolHttp.NewServer(":8888", "", "My App", "My App v0.0.0")
	myProducer.AddServer(httpServer)
	wsServer := protocolWebSocket.NewServer(httpServer)
	myProducer.AddServer(wsServer)

	myProducer.Expose()

	// 	for {
	// 		time.Sleep(10 * time.Second)
	// 		exposedThing.ExposedProperties["boolRWO"].Value = !(exposedThing.ExposedProperties["boolRWO"].Value.(bool))
	// 		exposedThing.EmitPropertyChange("boolRWO")
	// 	}

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

// func handlerA(interface{}) (interface{}, error) {
// 	println("a action")
// 	return nil, nil
// }

// func handlerB(value interface{}) (interface{}, error) {
// 	println("b action: " + value.(string))
// 	return nil, nil
// }

// func handlerC(value interface{}) (interface{}, error) {
// 	v := value.(string)
// 	if v != "c" {
// 		return nil, errors.New("the input string should be 'c'")
// 	}
// 	println("c action: " + v)
// 	return "ok", nil
// }
