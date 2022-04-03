package protocolHttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
)

type HttpServer struct {
	router     *httprouter.Router
	Host       string
	Port       uint
	thingChain []Middleware
	getChain   []Middleware
	putChain   []Middleware
	postChain  []Middleware
	*http.Server
}

type key int

const (
	keyDecodedJSON key = iota
)

func NewServer(host string, port uint) *HttpServer {
	address := fmt.Sprintf("%s:%d", host, port)

	router := httprouter.New()
	h := &HttpServer{
		router: router,
		Host:   host,
		Port:   port,
		thingChain: []Middleware{
			injectThing,
			corsHeader,
		},
		getChain: []Middleware{
			injectThing,
			corsHeader,
		},
		putChain: []Middleware{
			injectThing,
			corsHeader,
			decodeJSONResponse,
		},
		postChain: []Middleware{
			injectThing,
			corsHeader,
			decodeJSONResponse,
		},
		Server: &http.Server{
			Addr: address,
		},
	}
	return h
}

func (s *HttpServer) Expose(ref string, thing *producer.ExposedThing) {
	prefix := ""
	if ref != "" {
		prefix = "/" + ref
	}
	s.router.GET("/"+ref, buildChain(thing, HTTPGetThing, s.thingChain...))
	s.router.GET(prefix+"/:name", buildChain(thing, HTTPGet, s.getChain...))
	s.router.PUT(prefix+"/:name", buildChain(thing, HTTPPut, s.putChain...))
	s.router.POST(prefix+"/:name", buildChain(thing, HTTPPost, s.postChain...))
	addEndPoints(s.Host, s.Port, ref, thing)
}

// Produce constructs and launch an http server
func (s *HttpServer) Start() {
	s.router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header.Set("Allow", "GET,POST,PUT,DELETE,OPTIONS")
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})

	s.Server.Handler = s.router

	s.RegisterOnShutdown(func() {
		log.Debug().Msg("[protocolHttp:Start] Gracefully shutdown all websocket connections")
		// Wait for Gracefully shutdown all active websocket connections, for all things
		// for _, wsHandler := range s.wsHandlers {
		// 	wsHandler.gracefullWSShutdown()
		// }
	})

	go func() {
		log.Info().Msg("[protocolHttp:Start] Server listening")
		err := s.ListenAndServe()
		// always returns error. ErrServerClosed on graceful close
		if err == http.ErrServerClosed {
			log.Info().Msg("[protocolHttp:Start] Server closed")
		} else {
			// unexpected error. port in use?
			log.Error().Err(err).Msg("[protocolHttp:Start]")
		}
	}()
}

// stop Gracefully the server and all connections
func (s *HttpServer) Stop() {
	// if p == nil {
	// 	log.Error().Msg("[protocolHttp:GracefullyShutdown] nil server")
	// 	return
	// }
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// // https://golang.org/pkg/net/http/#Server.Shutdown
	// if err := p.Shutdown(ctx); err != nil {
	// 	log.Info().Msg("[thing:Shutdown] Shutdown error")
	// } else {
	// 	log.Debug().Msg("[thing:Shutdown] Wait on websocket connections shutdown")
	// 	// Wait for Gracefully shutdown all active websocket connections, for all things
	// 	// for _, wsHandler := range s.wsHandlers {
	// 	// 	wsHandler.waitWebSocket.Wait()
	// 	// }
	// 	log.Info().Msg("[thing:Shutdown] Gracefully stopped")
	// }
}

func addEndPoints(sHost string, sPort uint, ref string, t *producer.ExposedThing) {
	if t == nil {
		log.Error().Msg("[protocolHttp:GracefullyShutdown] nil thing")
		return
	}
	prefix := ""
	if ref != "" {
		prefix = "/" + ref
	}
	// var (
	// allReadOnly   = true
	// allWriteOnly  = true
	// anyProperties = false
	// )

	for _, property := range t.Td.Properties {
		// anyProperties = true

		form := &interaction.Form{
			ContentType: "application/json",
			Supplement:  map[string]interface{}{},
			UrlBuilder: func(host string, secure bool) string {
				protocol := "http"
				if secure {
					protocol = "https"
				}
				if sHost != "" { // force host
					return fmt.Sprintf("%s://%s:%d%s/%s", protocol, sHost, sPort, prefix, property.Key)
				} else {
					return fmt.Sprintf("%s://%s%s/%s", protocol, host, prefix, property.Key)
				}
			},
		}

		if !property.ReadOnly {
			// allReadOnly = false
		} else if !property.WriteOnly {
			// allWriteOnly = false
		}

		if property.ReadOnly {
			form.Op = []string{"readproperty"}
			form.Supplement["htv:methodName"] = "GET"
		} else if property.WriteOnly {
			form.Op = []string{"writeproperty"}
			form.Supplement["htv:methodName"] = "PUT"
		} else {
			form.Op = []string{"readproperty", "writeproperty"}
		}

		property.Forms = append(property.Forms, form)

		// if property is observable add an additional form with a observable href
		// if property.Observable {
		// 	form := interaction.Form{
		// 		Href:        href + property.Key,
		// 		ContentType: "application/json",
		// 	}
		// 	form.Op = []string{"observeproperty", "unobserveproperty"}
		// 	form.Subprotocol = "longpoll"
		// 	property.Forms = append(property.Forms, form)
		// }
	}

	// TODO
	// if anyProperties {
	// 	form := interaction.Form{
	// 		Href:        href,
	// 		ContentType: "application/json",
	// 	}
	// 	if allReadOnly {
	// 		form.Op = []string{"readallproperties", "readmultipleproperties"}
	// 	} else if allWriteOnly {
	// 		form.Op = []string{"writeallproperties", "writemultipleproperties"}
	// 	} else {
	// 		form.Op = []string{
	// 			"readallproperties",
	// 			"readmultipleproperties",
	// 			"writeallproperties",
	// 			"writemultipleproperties",
	// 		}
	// 	}
	// 	t.Td.Forms = append(t.Td.Forms, form)
	// }

	for _, action := range t.Td.Actions {
		form := &interaction.Form{
			ContentType: "application/json",
			Op:          []string{"invokeaction"},
			Supplement: map[string]interface{}{
				"htv:methodName": "POST",
			},
			UrlBuilder: func(host string, secure bool) string {
				protocol := "http"
				if secure {
					protocol = "https"
				}
				if sHost != "" { // force host
					return fmt.Sprintf("%s://%s:%d/%s/%s", protocol, sHost, sPort, prefix, action.Key)
				} else {
					return fmt.Sprintf("%s://%s/%s/%s", protocol, host, prefix, action.Key)
				}
			},
		}
		action.Forms = append(action.Forms, form)
	}

	// TODO
	// for _, event := range t.Td.Events {
	// 	form := interaction.Form{
	// 		Href:        href + event.Key,
	// 		ContentType: "application/json",
	// 		Op:          []string{"subscribeevent", "unsubscribeevent"},
	// 		Subprotocol: "longpoll"
	// 	}
	// 	event.Forms = append(event.Forms, form)
	// }
}

//jsonHTTPRenderer Add header and write response as json string
func jsonHTTPRenderer(w http.ResponseWriter, content interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	if body, ok := content.(string); ok {
		w.WriteHeader(status)
		io.WriteString(w, body)
	} else {
		body, err := json.Marshal(content)
		if err != nil {
			log.Error().Err(err).Msg("[thing:jsonHTTPRenderer]")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(status)
		io.WriteString(w, string(body))
	}
}

//okHTTPRenderer Add header and write response as ok: true
func okHTTPRenderer(w http.ResponseWriter, status int) {
	response := map[string]interface{}{"ok": true}
	jsonHTTPRenderer(w, response, status)
}

func errorHTTPRenderer(w http.ResponseWriter, errObj errorReturn, message string) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(errObj.httpStatus)
	body, err := json.Marshal(map[string]interface{}{
		"error": message,
		"type":  errObj.errorType,
	})
	if err != nil {
		log.Error().Err(err).Msg("[thing:errorHTTPRenderer]")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(body))
}
