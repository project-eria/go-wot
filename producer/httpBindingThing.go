package producer

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

// thingHandler is a set of handler for things root
type thingHandler struct {
	*ExposedThing
}

// get handle the GET method for single thing root
// @param {Object} w The response object
// @param {Object} r The request object
func (h *thingHandler) get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// if r.Header.Get("Upgrade") == "websocket" {
	// 	h.webSocket(w, r)
	// 	return
	// }
	td := h.GetThingDescription()
	content, err := json.Marshal(td)
	if err != nil {
		log.Error().Err(err).Msg("[producer:GetThingDescription]")
		errorHTTPRenderer(w, EncodingError, err.Error())
		return
	}
	jsonHTTPRenderer(w, string(content))
}
