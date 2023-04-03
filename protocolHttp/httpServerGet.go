package protocolHttp

import (
	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
)

// get handle the GET method for thing single property
// https://w3c.github.io/wot-scripting-api/#handling-requests-for-reading-a-property
func propertyReadHandler(t *producer.ExposedThing, tdProperty *interaction.Property) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Locals("websocket") == true {
			return c.Next() // delegate to the next handler (websocket)
		}
		log.Trace().Str("uri", c.Path()).Msg("[protocolHttp:propertyReadHandler] Received Thing property GET request")
		if tdProperty.WriteOnly {
			log.Trace().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyReadHandler] Access to WriteOnly property")
			return c.Status(NotAllowedError.httpStatus).JSON(fiber.Map{
				"error": "Write Only property",
				"type":  NotAllowedError.errorType,
			})
		} else {
			property := t.ExposedProperties[tdProperty.Key]
			handler := property.GetReadHandler()
			if handler != nil {
				content, err := handler(t, tdProperty.Key)
				if err != nil {
					log.Error().Str("uri", c.Path()).Err(err).Msg("[protocolHttp:propertyReadHandler]")
					return c.Status(UnknownError.httpStatus).JSON(fiber.Map{
						"error": err.Error(),
						"type":  UnknownError.errorType,
					})
				}
				log.Trace().Interface("response", content).Str("property", tdProperty.Key).Msg("[protocolHttp:propertyReadHandler] Response to Thing property GET request")
				return c.JSON(content)
			} else {
				log.Warn().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyReadHandler] Not Implemented")
				return c.Status(NotSupportedError.httpStatus).JSON(fiber.Map{
					"error": "Not Implemented",
					"type":  NotSupportedError.errorType,
				})
			}
		}
	}
}
