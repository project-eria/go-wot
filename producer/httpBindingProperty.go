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
	if r.Header.Get("Upgrade") == "websocket" {
		log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:GET] Received Thing property WS request")
		h.webSocket(name, w, r)
		return
	}
	log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:GET] Received Thing property GET request")
	if property, ok := h.Td.Properties[name]; ok {
		if property.WriteOnly {
			log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:GET] Access to WriteOnly property")
			errorHTTPRenderer(w, NotAllowedError, "Write Only property")
		} else {
			property := h.exposedProperties[name]
			handler := property.GetReadHandler()
			if handler != nil {
				content, err := handler(h.ExposedThing, name)
				if err != nil {
					log.Error().Str("uri", r.RequestURI).Err(err).Msg("[propertyHandler:GET]")
					errorHTTPRenderer(w, UnknownError, err.Error())
					return
				}
				log.Trace().Interface("response", content).Str("property", name).Msg("[propertyHandler:GET] Response to Thing property GET request")
				jsonHTTPRenderer(w, content, http.StatusOK)
			} else {
				log.Warn().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:GET] Not Implemented")
				errorHTTPRenderer(w, NotSupportedError, "Not Implemented")
			}
		}
		return
	}
	log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:GET] property not found")
	errorHTTPRenderer(w, NotFoundError, "Property not found")
}

// put handle the PUT method for thing single property
// https://w3c.github.io/wot-scripting-api/#handling-requests-for-writing-a-property
// @param {Object} w The response object
// @param {Object} r The request object
// @param {Object} params The url parmeters
func (h *propertyHandler) put(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("propertyName")
	log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:PUT] Received Thing property PUT request")
	if property, ok := h.Td.Properties[name]; ok {
		if property.ReadOnly {
			log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:PUT] Access to ReadOnly property")
			errorHTTPRenderer(w, NotAllowedError, "Read Only property")
		} else {
			property := h.exposedProperties[name]
			handler := property.GetWriteHandler()
			if handler != nil {
				data := r.Context().Value(keyDecodedJSON)
				if data == nil {
					log.Warn().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:PUT] No Data")
					errorHTTPRenderer(w, DataError, "No data provided")
					return
				}
				err := handler(h.ExposedThing, name, data)
				if err != nil {
					log.Error().Str("uri", r.RequestURI).Err(err).Msg("[propertyHandler:PUT]")
					errorHTTPRenderer(w, UnknownError, err.Error())
					return
				}
				log.Trace().Interface("response", "ok").Str("property", name).Msg("[propertyHandler:PUT] Response to Thing property PUT request")
				okHTTPRenderer(w, http.StatusOK)
			} else {
				log.Warn().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:PUT] Not Implemented")
				errorHTTPRenderer(w, NotSupportedError, "Not Implemented")
			}
		}
		return
	}
	log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:PUT] property not found")
	errorHTTPRenderer(w, NotFoundError, "Property not found")
}
