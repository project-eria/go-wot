package producer

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/project-eria/go-wot/form"
	"github.com/rs/zerolog/log"
)

type key int

const (
	keyDecodedJSON key = iota
)

func registerRoutes(prefix string, router *httprouter.Router, thing *ExposedThing) {
	var rootPrefix string
	if prefix != "" {
		rootPrefix = "/" + prefix
	}

	thingHandler := &thingHandler{thing}
	propertyHandler := &propertyHandler{thing}
	actionHandler := &actionHandler{thing}
	router.GET("/"+prefix, corsHeader(thingHandler.get))
	router.GET(rootPrefix+"/:propertyName", corsHeader(propertyHandler.get))
	router.PUT(rootPrefix+"/:propertyName", corsHeader(decodeJSON(propertyHandler.put)))
	router.POST(rootPrefix+"/:actionName", corsHeader(decodeJSON(actionHandler.post)))
}

// Produce constructs and launch an http server
func (p *Producer) exposeHttp() {
	router := httprouter.New()
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	if len(p.things) == 1 {
		thing := p.things[0]
		registerRoutes("", router, thing)
		//		server.wsHandlers = append(server.wsHandlers, thing.thingWSHandler)
		// } else {
		// thingsHandler := &thingsHandler{things}
		// router.GET("/", corsHeader(thingsHandler.get))
		// for _, thing := range things {
		// 	registerRoutes(thing.ref(), router, thing)
		// 	server.wsHandlers = append(server.wsHandlers, thing.thingWSHandler)
		// }
	}

	p.Server.Handler = router

	p.RegisterOnShutdown(func() {
		log.Debug().Msg("[thing:Shutdown] Gracefully shutdown all websocket connections")
		// Wait for Gracefully shutdown all active websocket connections, for all things
		// for _, wsHandler := range s.wsHandlers {
		// 	wsHandler.gracefullWSShutdown()
		// }
	})

	go func() {
		log.Info().Msg("[thingServer:Start] Server listening")
		err := p.ListenAndServe()
		// always returns error. ErrServerClosed on graceful close
		if err == http.ErrServerClosed {
			log.Info().Msg("[thingServer:Start] Server closed")
		} else {
			// unexpected error. port in use?
			log.Error().Err(err).Msg("[thingServer:Start]")
		}
	}()
}

func addFormHttp(e *ExposedThing, host string, secure bool) {
	http_url := "http://" + host
	ws_url := "ws://" + host
	if secure {
		http_url = "https://" + host
		ws_url = "wss://" + host
	}

	for _, property := range e.Td.Properties {
		op := []string{}
		if !property.ReadOnly {
			op = append(op, "writeproperty")
		}
		if !property.WriteOnly {
			op = append(op, "readproperty")
		}
		property.Interaction.AddForm(http_url,
			form.Form{
				ContentType: "application/json",
				Op:          op,
			},
		)
		if property.Observable {
			property.Interaction.AddForm(ws_url,
				form.Form{
					ContentType: "application/json",
					Op:          []string{"observeproperty", "unobserveproperty"},
				},
			)
		}
	}
	for _, action := range e.Td.Actions {
		action.Interaction.AddForm(http_url,
			form.Form{
				ContentType: "application/json",
				Op:          []string{"invokeaction"},
			},
		)
	}
}

// GracefullyShutdown Gracefully the server and all connections
func (p *Producer) GracefullyShutdown() {
	if p == nil {
		log.Error().Msg("[thingServer:GracefullyShutdown] nil server")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// https://golang.org/pkg/net/http/#Server.Shutdown
	if err := p.Shutdown(ctx); err != nil {
		log.Info().Msg("[thing:Shutdown] Shutdown error")
	} else {
		log.Debug().Msg("[thing:Shutdown] Wait on websocket connections shutdown")
		// Wait for Gracefully shutdown all active websocket connections, for all things
		// for _, wsHandler := range s.wsHandlers {
		// 	wsHandler.waitWebSocket.Wait()
		// }
		log.Info().Msg("[thing:Shutdown] Gracefully stopped")
	}
}

//jsonHTTPRenderer Add header and write response as json string
func jsonHTTPRenderer(w http.ResponseWriter, content interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if body, ok := content.(string); ok {
		io.WriteString(w, body)
	} else {
		body, err := json.Marshal(content)
		if err != nil {
			log.Error().Err(err).Msg("[thing:jsonHTTPRenderer]")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.WriteString(w, string(body))
	}
}

//okHTTPRenderer Add header and write response as ok: true
func okHTTPRenderer(w http.ResponseWriter) {
	response := map[string]interface{}{"ok": true}
	jsonHTTPRenderer(w, response)
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

func decodeJSON(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if next != nil {
			if r.Header.Get("Content-type") != "application/json" {
				log.Error().Msg("[thing:decodeJSON] Request without Content-type='application/json'")
				// break here instead of continuing the chain
				errorHTTPRenderer(w, EncodingError, "Content-type 'application/json' is required")
				return
			}
			var value interface{}
			err := json.NewDecoder(r.Body).Decode(&value)
			if err == io.EOF { // no JSON
				next(w, r, p)
				return
			}
			if err != nil {
				log.Error().Err(err).Msg("[thing:decodeJSON]")
				errorHTTPRenderer(w, EncodingError, "Incorrect JSON value")
				return
			}

			ctx := context.WithValue(r.Context(), keyDecodedJSON, value)
			next(w, r.WithContext(ctx), p)
		}
	}
}

func corsHeader(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if next != nil {
			header := w.Header()
			header.Set("Allow", "GET,POST,PUT,DELETE,OPTIONS")
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

			next(w, r, p)
		}
	}
}
