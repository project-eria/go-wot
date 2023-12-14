package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/protocolHttp"
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

	url := "http://pool.local:8889/device.0"
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
	consumedThing := myConsumer.Consume(&td)

	fmt.Println(consumedThing.GetThingDescription().Title)
	value, err := consumedThing.ReadProperty("boolRWO", nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(value)
	}

	value, err = consumedThing.ReadProperty("boolRW", nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(value)
	}

	value, err = consumedThing.WriteProperty("boolRW", nil, true)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(value)
	}

	value, err = consumedThing.ReadProperty("boolRW", nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(value)
	}

	value, err = consumedThing.ReadProperty("boolR", nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(value)
	}

	value, err = consumedThing.ReadProperty("boolW", nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(value)
	}

	value, err = consumedThing.InvokeAction("a", nil, nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(value)
	}

	myConsumer.Shutdown()
}
