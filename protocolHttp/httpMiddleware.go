package protocolHttp

import (
	"context"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
)

type Middleware func(*producer.ExposedThing, httprouter.Handle) httprouter.Handle

// buildChain builds the middlware chain recursively, functions are first class
func buildChain(thing *producer.ExposedThing, f httprouter.Handle, m ...Middleware) httprouter.Handle {
	// if our chain is done, use the original handlerfunc
	if len(m) == 0 {
		return f
	}
	// otherwise nest the handlerfuncs
	return m[0](thing, buildChain(thing, f, m[1:]...))
}

func injectThing(thing *producer.ExposedThing, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if next != nil {
			ctx := context.WithValue(r.Context(), "thing", thing)
			next(w, r.WithContext(ctx), p)
		}
	}
}

func corsHeader(thing *producer.ExposedThing, next httprouter.Handle) httprouter.Handle {
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

func decodeJSONResponse(thing *producer.ExposedThing, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if next != nil {
			if r.Header.Get("Content-type") != "application/json" {
				log.Error().Msg("[HTTPMiddleWare:decodeJSONResponse] Request without Content-type='application/json'")
				// break here instead of continuing the chain
				errorHTTPRenderer(w, EncodingError, "Content-type 'application/json' is required")
				return
			}
			value, err := decodeJSON(r.Body)
			if err == io.EOF { // no JSON
				next(w, r, p)
				return
			}
			if err != nil {
				log.Error().Err(err).Msg("[HTTPMiddleWare:decodeJSONResponse]")
				errorHTTPRenderer(w, EncodingError, "Incorrect JSON value")
				return
			}

			ctx := context.WithValue(r.Context(), keyDecodedJSON, value)
			next(w, r.WithContext(ctx), p)
		}
	}
}

func (s *HttpServer) AddGetMiddleware(m Middleware) {
	s.getChain = append(s.getChain, m)
}

func (s *HttpServer) AddPutMiddleware(m Middleware) {
	s.putChain = append(s.putChain, m)
}

func (s *HttpServer) AddPostMiddleware(m Middleware) {
	s.postChain = append(s.postChain, m)
}
