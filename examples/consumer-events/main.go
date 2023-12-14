package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/protocolWebSocket"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func main() {
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "02/01|15:04:05"})
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local)
	}
	defer func() {
		zlog.Info().Msg("[main] Stopped")
	}()

	url := "http://127.0.0.1:8888/"
	resp, err := http.Get(url)
	if err != nil {
		zlog.Error().Err(err).Msg("[main]")
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		zlog.Error().Str("status", resp.Status).Str("url", url).Msg("[main] incorrect response")
		os.Exit(1)
	}

	var td thing.Thing
	if err := json.NewDecoder(resp.Body).Decode(&td); err != nil {
		zlog.Error().Str("url", url).Err(err).Msg("[main]")
		os.Exit(1)
	}

	myConsumer := consumer.New()
	httpClient := protocolHttp.NewClient()
	myConsumer.AddClient(httpClient)
	wsClient := protocolWebSocket.NewClient()
	myConsumer.AddClient(wsClient)
	consumedThing := myConsumer.Consume(&td)

	value, err := consumedThing.ReadProperty("boolRWO", nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(value)
	}

	consumedThing.ObserveProperty("boolRWO", nil, func(value interface{}, err error) {
		fmt.Println(value)
	})

	// consumedThing.SubscribeEvent("d", func(value interface{}, err error) {
	// 	fmt.Println("d")
	// })

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	// Block until keyboard interrupt is received.
	<-c

	myConsumer.Shutdown()
}
