package protocolHttp

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
)

// get handle the GET method for thing single property
// https://w3c.github.io/wot-scripting-api/#handling-requests-for-reading-a-property
// @param {Object} w The response object
// @param {Object} r The request object
// @param {Object} params The url parmeters
func HTTPGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	t := r.Context().Value("thing").(*producer.ExposedThing)
	name := params.ByName("name")
	log.Trace().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:GET] Received Thing property GET request")
	if property, ok := t.Td.Properties[name]; ok {
		if property.WriteOnly {
			log.Trace().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:GET] Access to WriteOnly property")
			errorHTTPRenderer(w, NotAllowedError, "Write Only property")
		} else {
			property := t.ExposedProperties[name]
			handler := property.GetReadHandler()
			if handler != nil {
				content, err := handler(t, name)
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
	log.Trace().Str("uri", r.RequestURI).Str("property", name).Msg("[propertyHandler:GET] property not found")
	errorHTTPRenderer(w, NotFoundError, "Property not found")
}
