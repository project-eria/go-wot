package protocolHttp

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

// put handle the PUT method for thing single property
// https://w3c.github.io/wot-scripting-api/#handling-requests-for-writing-a-property
func propertyWriteHandler(t producer.ExposedThing, tdProperty *interaction.Property) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		optionsStr := c.AllParams()
		zlog.Trace().Str("uri", c.Path()).Interface("options", optionsStr).Msg("[protocolHttp:propertyWriteHandler] Received Thing property PUT request")
		if tdProperty.ReadOnly {
			zlog.Trace().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyWriteHandler] Access to ReadOnly property")
			return c.Status(NotAllowedError.HttpStatus).JSON(fiber.Map{
				"error": "Read Only property",
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

			var data interface{}
			if len(c.Body()) > 0 {
				if err := c.BodyParser(&data); err != nil {
					fmt.Println(err)
					return c.Status(EncodingError.HttpStatus).JSON(fiber.Map{
						"error": "Incorrect JSON value",
						"type":  EncodingError.ErrorType,
					})
				}
			}
			err = property.Write(t, tdProperty.Key, data, optionsStr)
			if err != nil {
				zlog.Error().Err(err).Str("property", tdProperty.Key).Msg("[protocolHttp:propertyWriteHandler]")

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
			zlog.Trace().Interface("response", "ok").Str("property", tdProperty.Key).Msg("[protocolHttp:propertyWriteHandler] Response to Thing property PUT request")

			return c.JSON(fiber.Map{"ok": true})
		}
	}
}
