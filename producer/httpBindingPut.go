package producer

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

// put handle the PUT method for thing single property
// https://w3c.github.io/wot-scripting-api/#handling-requests-for-writing-a-property
// @param {Object} w The response object
// @param {Object} r The request object
// @param {Object} params The url parmeters
func (t *ExposedThing) HTTPPut(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:PUT] Received Thing property PUT request")
	if property, ok := t.Td.Properties[name]; ok {
		if property.ReadOnly {
			log.Debug().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:PUT] Access to ReadOnly property")
			errorHTTPRenderer(w, NotAllowedError, "Read Only property")
		} else {
			property := t.exposedProperties[name]
			handler := property.GetWriteHandler()
			if handler != nil {
				data := r.Context().Value(keyDecodedJSON)
				if data == nil {
					log.Warn().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:PUT] No Data")
					errorHTTPRenderer(w, DataError, "No data provided")
					return
				}
				err := handler(t, name, data)
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
