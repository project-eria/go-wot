package protocolHttp

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
)

// get handle the GET method for thing single property
// https://w3c.github.io/wot-scripting-api/#handling-requests-for-reading-a-property
func propertyReadHandler(t producer.ExposedThing, tdProperty *interaction.Property) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		options := c.AllParams()
		log.Trace().Str("uri", c.Path()).Interface("options", options).Msg("[protocolHttp:propertyReadHandler] Received Thing property GET request")
		if tdProperty.WriteOnly {
			log.Trace().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyReadHandler] Access to WriteOnly property")
			return c.Status(NotAllowedError.HttpStatus).JSON(fiber.Map{
				"error": "Write Only property",
				"type":  NotAllowedError.ErrorType,
			})
		} else {
			property, err := t.ExposedProperty(tdProperty.Key)
			if err != nil {
				log.Error().Err(err).Str("property", tdProperty.Key).Msg("[protocolHttp:propertyReadHandler]")
				return c.Status(UnknownError.HttpStatus).JSON(fiber.Map{
					"error": fmt.Sprintf("ExposedProperty `%s` not found", tdProperty.Key),
					"type":  UnknownError.ErrorType,
				})
			} else {
				handler := property.GetReadHandler()
				if handler != nil {
					// Check the options (uriVariables) data
					if err := property.CheckUriVariables(options); err != nil {
						return c.Status(DataError.HttpStatus).JSON(fiber.Map{
							"error": err.Error(),
							"type":  DataError.ErrorType,
						})
					}
					// Call the function that handle the property read
					content, err := handler(t, tdProperty.Key, options)
					if err != nil {
						log.Error().Str("uri", c.Path()).Err(err).Msg("[protocolHttp:propertyReadHandler]")
						return c.Status(UnknownError.HttpStatus).JSON(fiber.Map{
							"error": err.Error(),
							"type":  UnknownError.ErrorType,
						})
					}
					log.Trace().Interface("response", content).Str("property", tdProperty.Key).Msg("[protocolHttp:propertyReadHandler] Response to Thing property GET request")
					return c.JSON(content)
				} else {
					log.Warn().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyReadHandler] Not Implemented")
					return c.Status(NotSupportedError.HttpStatus).JSON(fiber.Map{
						"error": "Not Implemented",
						"type":  NotSupportedError.ErrorType,
					})
				}
			}
		}
	}
}
