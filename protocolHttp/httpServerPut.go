package protocolHttp

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
	zlog "github.com/rs/zerolog/log"
)

// put handle the PUT method for thing single property
// https://w3c.github.io/wot-scripting-api/#handling-requests-for-writing-a-property
func propertyWriteHandler(t *producer.ExposedThing, tdProperty *interaction.Property) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		options := c.AllParams()
		zlog.Trace().Str("uri", c.Path()).Interface("options", options).Msg("[protocolHttp:propertyWriteHandler] Received Thing property PUT request")
		if tdProperty.ReadOnly {
			zlog.Trace().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyWriteHandler] Access to ReadOnly property")
			return c.Status(NotAllowedError.HttpStatus).JSON(fiber.Map{
				"error": "Read Only property",
				"type":  NotAllowedError.ErrorType,
			})
		} else {
			if property, ok := t.ExposedProperties[tdProperty.Key]; ok {
				handler := property.GetWriteHandler()
				if handler != nil {
					// Check the params (uriVariables) data
					if err := property.CheckUriVariables(options); err != nil {
						return c.Status(DataError.HttpStatus).JSON(fiber.Map{
							"error": err.Error(),
							"type":  DataError.ErrorType,
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

					// Check if data has been provided
					if data == nil {
						zlog.Warn().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyWriteHandler] No Data")
						return c.Status(DataError.HttpStatus).JSON(fiber.Map{
							"error": "No data provided",
							"type":  DataError.ErrorType,
						})
					}

					// Check the data sent format
					if err := property.Data.Check(data); err != nil {
						message := "incorrect input value: " + err.Error()
						log.Trace().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyWriteHandler] " + message)
						return c.Status(DataError.HttpStatus).JSON(fiber.Map{
							"error": message,
							"type":  DataError.ErrorType,
						})
					}

					// Call the function that handle the property write
					err := handler(t, tdProperty.Key, data, options)
					if err != nil {
						zlog.Error().Err(err).Msg("[protocolHttp:propertyWriteHandler]")
						return c.Status(UnknownError.HttpStatus).JSON(fiber.Map{
							"error": err.Error(),
							"type":  UnknownError.ErrorType,
						})
					}
					zlog.Trace().Interface("response", "ok").Str("property", tdProperty.Key).Msg("[protocolHttp:propertyWriteHandler] Response to Thing property PUT request")

					// Notify all listeners that the property changed
					if err := t.EmitPropertyChange(tdProperty.Key, data, options); err != nil {
						zlog.Error().Str("property", tdProperty.Key).Interface("value", data).Err(err).Msg("[protocolHttp:propertyWriteHandler]")
						return err
					}

					return c.JSON(fiber.Map{"ok": true})
				} else {
					zlog.Warn().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyWriteHandler] Not Implemented")
					return c.Status(NotSupportedError.HttpStatus).JSON(fiber.Map{
						"error": "Not Implemented",
						"type":  NotSupportedError.ErrorType,
					})
				}
			} else {
				zlog.Error().Str("property", tdProperty.Key).Msg("[protocolHttp:propertyReadHandler] ExposedProperty not found")
				return c.Status(UnknownError.HttpStatus).JSON(fiber.Map{
					"error": fmt.Errorf("ExposedProperty `%s` not found", tdProperty.Key),
					"type":  UnknownError.ErrorType,
				})
			}
		}
	}
}
