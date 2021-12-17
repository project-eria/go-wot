package producer

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

// propertyHandler Handle a request to single property /{name}.
type propertyHandler struct {
	*ExposedThing
}

// get handle the GET method for thing single property
// https://w3c.github.io/wot-scripting-api/#handling-requests-for-reading-a-property
// @param {Object} w The response object
// @param {Object} r The request object
// @param {Object} params The url parmeters
func (h *propertyHandler) get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("propertyName")
	log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[property:GET] Received Thing property GET request")
	if property, ok := h.td.Properties[name]; ok {
		if property.WriteOnly {
			log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[property:GET] Access to WriteOnly property")
			errorHTTPRenderer(w, NotAllowedError, "Write Only property")
		} else {
			if handler, ok := h.propertiesReadHandlers[name]; ok {
				content, err := handler()
				if err != nil {
					log.Error().Str("uri", r.RequestURI).Err(err).Msg("[property:GET]")
					errorHTTPRenderer(w, UnknownError, err.Error())
					return
				}
				log.Trace().Interface("response", content).Str("property", name).Msg("[property:GET] Response to Thing property GET request")
				jsonHTTPRenderer(w, content)
			} else {
				log.Warn().Str("uri", r.RequestURI).Str("property", name).Msg("[property:GET] Not Implemented")
				errorHTTPRenderer(w, NotSupportedError, "Not Implemented")
			}
		}
		return
	}
	log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[property:GET] property not found")
	errorHTTPRenderer(w, NotFoundError, "Property not found")
}

// put handle the PUT method for thing single property
// https://w3c.github.io/wot-scripting-api/#handling-requests-for-writing-a-property
// @param {Object} w The response object
// @param {Object} r The request object
// @param {Object} params The url parmeters
func (h *propertyHandler) put(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("propertyName")
	log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[property:PUT] Received Thing property PUT request")
	if property, ok := h.td.Properties[name]; ok {
		if property.ReadOnly {
			log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[property:PUT] Access to ReadOnly property")
			errorHTTPRenderer(w, NotAllowedError, "Read Only property")
		} else {
			if handler, ok := h.propertiesWriteHandlers[name]; ok {
				data := r.Context().Value(keyDecodedJSON)
				err := handler(data)
				if err != nil {
					log.Error().Str("uri", r.RequestURI).Err(err).Msg("[property:PUT]")
					errorHTTPRenderer(w, UnknownError, err.Error())
					return
				}
				log.Trace().Interface("response", "ok").Str("property", name).Msg("[property:PUT] Response to Thing property PUT request")
				okHTTPRenderer(w)
			} else {
				log.Warn().Str("uri", r.RequestURI).Str("property", name).Msg("[property:PUT] Not Implemented")
				errorHTTPRenderer(w, NotSupportedError, "Not Implemented")
			}
		}
		return
	}
	log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[property:PUT] property not found")
	errorHTTPRenderer(w, NotFoundError, "Property not found")
}
