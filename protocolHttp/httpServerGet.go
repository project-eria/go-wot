package protocolHttp

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

// get handle the GET method for thing single property
// https://w3c.github.io/wot-scripting-api/#handling-requests-for-reading-a-property
func propertyReadHandler(t producer.ExposedThing, tdProperty *interaction.Property) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		optionsStr := c.AllParams()
		zlog.Trace().Str("uri", c.Path()).Interface("options", optionsStr).Msg("[protocolHttp:propertyReadHandler] Received Thing property GET request")
		if tdProperty.WriteOnly {
			zlog.Trace().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyReadHandler] Access to WriteOnly property")
			return c.Status(NotAllowedError.HttpStatus).JSON(fiber.Map{
				"error": "Write Only property",
				"type":  NotAllowedError.ErrorType,
			})
		} else {
			property, err := t.ExposedProperty(tdProperty.Key)
			if err != nil {
				zlog.Error().Err(err).Str("property", tdProperty.Key).Msg("[protocolHttp:propertyReadHandler]")
				return c.Status(UnknownError.HttpStatus).JSON(fiber.Map{
					"error": fmt.Sprintf("ExposedProperty `%s` not found", tdProperty.Key),
					"type":  UnknownError.ErrorType,
				})
			}
			output, err := property.Read(t, tdProperty.Key, optionsStr)
			if err != nil {
				zlog.Error().Err(err).Str("property", tdProperty.Key).Msg("[protocolHttp:propertyReadHandler]")

				if _, ok := err.(*producer.DataError); ok {
					return c.Status(DataError.HttpStatus).JSON(fiber.Map{
						"error": err.Error(),
						"type":  DataError.ErrorType,
					})
				} else if _, ok := err.(*producer.NotImplementedError); ok {
					return c.Status(NotSupportedError.HttpStatus).JSON(fiber.Map{
						"error": err.Error(),
						"type":  NotSupportedError.ErrorType,
					})
				} else {
					return c.Status(UnknownError.HttpStatus).JSON(fiber.Map{
						"error": err.Error(),
						"type":  UnknownError.ErrorType,
					})
				}
			}

			zlog.Trace().Interface("response", output).Str("property", tdProperty.Key).Msg("[protocolHttp:propertyReadHandler] Response to Thing property GET request")
			return c.JSON(output)
		}
	}
}
