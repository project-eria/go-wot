package protocolHttp

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
)

// put handle the PUT method for thing single property
// https://w3c.github.io/wot-scripting-api/#handling-requests-for-writing-a-property
func propertyWriteHandler(t *producer.ExposedThing, tdProperty *interaction.Property) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		log.Trace().Str("uri", c.Path()).Msg("[protocolHttp:propertyWriteHandler] Received Thing property PUT request")
		if tdProperty.ReadOnly {
			log.Trace().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyWriteHandler] Access to ReadOnly property")
			return c.Status(NotAllowedError.httpStatus).JSON(fiber.Map{
				"error": "Read Only property",
				"type":  NotAllowedError.errorType,
			})
		} else {
			property := t.ExposedProperties[tdProperty.Key]
			var data interface{}
			if err := c.BodyParser(&data); err != nil {
				fmt.Println(err)
				return c.Status(EncodingError.httpStatus).JSON(fiber.Map{
					"error": "Incorrect JSON value",
					"type":  EncodingError.errorType,
				})
			}
			if data == nil {
				log.Warn().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyWriteHandler] No Data")
				return c.Status(DataError.httpStatus).JSON(fiber.Map{
					"error": "No data provided",
					"type":  DataError.errorType,
				})
			}
			handler := property.GetWriteHandler()
			if handler != nil {
				err := handler(t, tdProperty.Key, data)
				if err != nil {
					log.Error().Err(err).Msg("[protocolHttp:propertyWriteHandler]")
					return c.Status(UnknownError.httpStatus).JSON(fiber.Map{
						"error": err.Error(),
						"type":  UnknownError.errorType,
					})
				}
				log.Trace().Interface("response", "ok").Str("property", tdProperty.Key).Msg("[protocolHttp:propertyWriteHandler] Response to Thing property PUT request")
				return c.JSON(fiber.Map{"ok": true})
			} else {
				log.Warn().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyWriteHandler] Not Implemented")
				return c.Status(NotSupportedError.httpStatus).JSON(fiber.Map{
					"error": "Not Implemented",
					"type":  NotSupportedError.errorType,
				})
			}
		}
	}
}
