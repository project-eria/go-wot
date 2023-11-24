package protocolHttp

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
)

type HttpServer struct {
	addr        string
	ExposedAddr string
	*fiber.App
}

func NewServer(addr string, exposedAddr string, header string, appName string) *HttpServer {
	router := fiber.New(fiber.Config{
		ServerHeader: header,
		AppName:      appName,
	})

	router.Use(checkContentType())
	router.Use(corsHeader())

	h := &HttpServer{
		addr:        addr,
		ExposedAddr: exposedAddr,
		App:         router,
	}
	return h
}

func (s *HttpServer) Expose(ref string, thing producer.ExposedThing) {
	prefix := ""
	if ref != "" {
		prefix = "/" + ref
	}
	s.Get("/"+ref, thingHandler(thing))
	g := s.Group(prefix)
	addEndPoints(g, s.ExposedAddr, prefix, thing)
	// propertyChangeChan := thing.GetPropertyChangeChannel()
	// eventChan := thing.GetEventChannel()
}

// Produce constructs and launch an http server
func (s *HttpServer) Start() {
	// TODO
	// s.RegisterOnShutdown(func() {
	// 	log.Trace().Msg("[protocolHttp:Start] Gracefully shutdown all websocket connections")
	// 	// Wait for Gracefully shutdown all active websocket connections, for all things
	// 	// for _, wsHandler := range s.wsHandlers {
	// 	// 	wsHandler.gracefullWSShutdown()
	// 	// }
	// })

	go func() {
		log.Info().Msg("[protocolHttp:Start] Server listening")
		err := s.Listen(s.addr)
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
	// 	log.Trace().Msg("[thing:Shutdown] Wait on websocket connections shutdown")
	// 	// Wait for Gracefully shutdown all active websocket connections, for all things
	// 	// for _, wsHandler := range s.wsHandlers {
	// 	// 	wsHandler.waitWebSocket.Wait()
	// 	// }
	// 	log.Info().Msg("[thing:Shutdown] Gracefully stopped")
	// }
}

func addEndPoints(g fiber.Router, exposedAddr string, prefix string, t producer.ExposedThing) {
	// var (
	// allReadOnly   = true
	// allWriteOnly  = true
	// anyProperties = false
	// )

	for _, property := range t.TD().Properties {
		addPropertyEndPoints(g, exposedAddr, prefix, t, property)
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

	for _, action := range t.TD().Actions {
		addActionEndPoints(g, exposedAddr, prefix, t, action)
	}

	// TODO
	// https://github.com/gofiber/fiber/issues/646
	// https://github.com/LdDl/fiber-long-poll
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

func addPropertyEndPoints(g fiber.Router, exposedAddr string, prefix string, t producer.ExposedThing, property *interaction.Property) {
	var uriVars string
	var handlerVars string
	// https://w3c.github.io/wot-thing-description/#form-uriVariables
	// How to decide /{city} or {?unit} format?
	if len(property.UriVariables) > 0 {
		for uriVar := range property.UriVariables {
			uriVars += fmt.Sprintf("/{%s}", uriVar)
			handlerVars += fmt.Sprintf("/:%s", uriVar)
		}
	}
	form := &interaction.Form{
		ContentType: "application/json",
		Supplement:  map[string]interface{}{},
		UrlBuilder: func(host string, secure bool) string {
			protocol := "http"
			if secure {
				protocol = "https"
			}
			if exposedAddr != "" { // force exposed host
				host = exposedAddr
			}
			return fmt.Sprintf("%s://%s%s/%s%s", protocol, host, prefix, property.Key, uriVars)
		},
	}

	// if !property.ReadOnly {
	// 	// allReadOnly = false
	// } else if !property.WriteOnly {
	// 	// allWriteOnly = false
	// }

	if property.ReadOnly {
		form.Op = []string{"readproperty"}
		form.Supplement["htv:methodName"] = "GET"
	} else if property.WriteOnly {
		form.Op = []string{"writeproperty"}
		form.Supplement["htv:methodName"] = "PUT"
	} else {
		form.Op = []string{"readproperty", "writeproperty"}
	}
	g.Get("/"+property.Key+handlerVars, propertyReadHandler(t, property))
	g.Put("/"+property.Key+handlerVars, propertyWriteHandler(t, property))

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

func addActionEndPoints(g fiber.Router, exposedAddr string, prefix string, t producer.ExposedThing, action *interaction.Action) {
	var uriVars string
	var handlerVars string
	// https://w3c.github.io/wot-thing-description/#form-uriVariables
	// How to decide /{city} or {?unit} format?
	if len(action.UriVariables) > 0 {
		for uriVar := range action.UriVariables {
			uriVars += fmt.Sprintf("/{%s}", uriVar)
			handlerVars += fmt.Sprintf("/:%s", uriVar)
		}
	}

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
			if exposedAddr != "" { // force exposed host
				host = exposedAddr
			}
			return fmt.Sprintf("%s://%s%s/%s%s", protocol, host, prefix, action.Key, uriVars)
		},
	}

	g.Post("/"+action.Key+handlerVars, actionHandler(t, action))

	action.Forms = append(action.Forms, form)
}
