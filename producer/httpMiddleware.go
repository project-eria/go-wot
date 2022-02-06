package producer

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

type Middleware func(*ExposedThing, httprouter.Handle) httprouter.Handle

var thingChain = []Middleware{
	corsHeader,
}

var getChain = []Middleware{
	corsHeader,
	updateWebsocket,
}

var putChain = []Middleware{
	corsHeader,
	decodeJSON,
}

var postChain = []Middleware{
	corsHeader,
	decodeJSON,
}

// buildChain builds the middlware chain recursively, functions are first class
func buildChain(thing *ExposedThing, f httprouter.Handle, m ...Middleware) httprouter.Handle {
	// if our chain is done, use the original handlerfunc
	if len(m) == 0 {
		return f
	}
	// otherwise nest the handlerfuncs
	return m[0](thing, buildChain(thing, f, m[1:cap(m)]...))
}

func corsHeader(thing *ExposedThing, next httprouter.Handle) httprouter.Handle {
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

func decodeJSON(thing *ExposedThing, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if next != nil {
			if r.Header.Get("Content-type") != "application/json" {
				log.Error().Msg("[HTTPMiddleWare:decodeJSON] Request without Content-type='application/json'")
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
				log.Error().Err(err).Msg("[HTTPMiddleWare:decodeJSON]")
				errorHTTPRenderer(w, EncodingError, "Incorrect JSON value")
				return
			}

			ctx := context.WithValue(r.Context(), keyDecodedJSON, value)
			next(w, r.WithContext(ctx), p)
		}
	}
}

func updateWebsocket(thing *ExposedThing, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if next != nil {
			if r.Header.Get("Upgrade") == "websocket" {

				// // webthing SubProtocol do not exists
				// if r.Header.Get("Sec-Websocket-Protocol") != "webthing" {
				// 	log.Error().Msg("[producer:webSocket] Connection not using webthing protocol")
				// 	w.WriteHeader(http.StatusBadRequest)
				// 	io.WriteString(w, "Connection not using webthing protocol")
				// 	return
				// }

				thing.WSGet(w, r, p)
				return
			} else {
				next(w, r, p)
			}
		}
	}
}
