package producer

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

// get handle the GET method for single thing root
// @param {Object} w The response object
// @param {Object} r The request object
func (t *ExposedThing) HTTPGetThing(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Debug().Str("uri", r.RequestURI).Msg("[thingHandler:GET] Received Thing GET request")

	td := t.GetThingDescription()
	content, err := json.Marshal(td)
	if err != nil {
		log.Error().Err(err).Msg("[producer:GetThingDescription]")
		errorHTTPRenderer(w, EncodingError, err.Error())
		return
	}
	jsonHTTPRenderer(w, string(content), http.StatusOK)
}
