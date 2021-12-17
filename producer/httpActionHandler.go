package producer

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

// actionHandler Handle a action requests to /actions/<name>
type actionHandler struct {
	*ExposedThing
}

// post handle the POST request method for a thing action
// https://w3c.github.io/wot-scripting-api/#handling-action-requests
// @param {Object} w The response object
// @param {Object} r The request object
// @param {Object} params The url parmeters
func (h *actionHandler) post(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("actionName")
	log.Debug().Str("uri", r.RequestURI).Str("action", name).Msg("[action:POST] Received Thing action POST request")

	if action, ok := h.td.Actions[name]; ok {
		if handler, ok := h.actionHandlers[name]; ok {
			input := r.Context().Value(keyDecodedJSON)
			// Check the input data
			if action.Input != nil {
				if err := action.Input.Check(input); err != nil {
					message := "incorrect input value: " + err.Error()
					log.Trace().Str("uri", r.RequestURI).Str("action", name).Msg("[action:POST] " + message)
					errorHTTPRenderer(w, DataError, message)
					return
				}
			}
			// Execute the action requests
			output, err := handler(input)
			if err != nil {
				log.Error().Str("uri", r.RequestURI).Str("action", name).Err(err).Msg("[action:POST]")
				errorHTTPRenderer(w, UnknownError, err.Error())
				return
			}

			// Check the output data
			if action.Output != nil {
				if err := action.Output.Check(output); err != nil {
					log.Error().Str("uri", r.RequestURI).Str("action", name).Err(err).Msg("[action:POST] incorrect handler returned value")
					errorHTTPRenderer(w, UnknownError, "Incorrect handler returned value")
					return
				}
				log.Trace().Str("uri", r.RequestURI).Interface("response", output).Str("action", name).Msg("[action:POST] JSON Response to Thing action POST request")
				jsonHTTPRenderer(w, output)
				return
			}

			log.Trace().Str("uri", r.RequestURI).Str("action", name).Msg("[action:POST] OK Response to Thing action POST request")
			okHTTPRenderer(w)
			return
		} else {
			log.Warn().Str("uri", r.RequestURI).Str("action", name).Msg("[action:POST] no handler function for the action")
			errorHTTPRenderer(w, NotSupportedError, "Not Implemented")
		}
	}
	log.Debug().Str("uri", r.RequestURI).Msgf("[thing:post] action /%s not found", name)
	errorHTTPRenderer(w, NotFoundError, "Action not found")
}
