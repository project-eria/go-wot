package wot

import (
	"context"
	"sync"
)

var (
	_showVersion *bool
	_logLevel    *string
	_configPath  *string
	_version     string
	_appName     string
	_cancel      context.CancelFunc
	_ctx         context.Context
	_wait        sync.WaitGroup
)

// // RunSingleThing lauch a server for HTTP and WS requests
// func RunSingleThing(t *thing.Thing, port int) {
// 	_wait.Add(1)
// 	server := server.New([]*thing.Thing{t}, port)
// 	server.Start()
// 	go func() {
// 		<-_ctx.Done()
// 		server.GracefullyShutdown()
// 		_wait.Done()
// 	}()
// }

// // RunMultipleThing lauch a server for HTTP and WS requests
// func RunMultipleThing(t []*thing.Thing, port int) {
// 	_wait.Add(1)
// 	server := server.New(t, port)
// 	server.Start()
// 	go func() {
// 		<-_ctx.Done()
// 		server.GracefullyShutdown()
// 		_wait.Done()
// 	}()
// }

// // ConsumeThing connect a remote thing WS server
// func ConsumeThing(u string, wsRequired bool) (*client.ThingConnection, error) {
// 	thingURL, err := url.Parse(u)
// 	if err != nil {
// 		log.Fatal().Str("url", u).Err(err).Msg("[eria] Can't parse thing url. The url format should be `//<host>:<port>/<thing>`")
// 		return nil, errors.New("Can't parse thing url")
// 	}

// 	conn, err := client.New(*thingURL, wsRequired)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if wsRequired {
// 		_wait.Add(1)
// 		go func() {
// 			conn.ConnectWebSocket(_ctx)
// 			_wait.Done()
// 		}()
// 	}
// 	return conn, nil
// }
