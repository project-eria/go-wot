package protocolHttp

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

// post handle the POST request method for a thing action
// https://w3c.github.io/wot-scripting-api/#handling-action-requests
func actionHandler(t producer.ExposedThing, tdAction *interaction.Action) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		options := c.AllParams()
		zlog.Trace().Str("uri", c.Path()).Interface("options", options).Msg("[protocolHttp:actionHandler] Received Thing action POST request")
		action, err := t.ExposedAction(tdAction.Key)
		if err != nil {
			zlog.Error().Err(err).Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler]")
			return c.Status(UnknownError.HttpStatus).JSON(fiber.Map{
				"error": fmt.Sprintf("ExposedAction `%s` not found", tdAction.Key),
				"type":  UnknownError.ErrorType,
			})
		} else {
			handler := action.GetHandler()
			if handler != nil {
				// Check the params (uriVariables) data
				if err := action.CheckUriVariables(options); err != nil {
					return c.Status(DataError.HttpStatus).JSON(fiber.Map{
						"error": err.Error(),
						"type":  DataError.ErrorType,
					})
				}

				var data interface{}
				if len(c.Body()) > 0 {
					if err := c.BodyParser(&data); err != nil {
						zlog.Trace().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] Incorrect JSON value")
						return c.Status(EncodingError.HttpStatus).JSON(fiber.Map{
							"error": "Incorrect JSON value",
							"type":  EncodingError.ErrorType,
						})
					}
				} else {
					zlog.Trace().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] No body data")
				}

				// Check the input data
				if action.Input() != nil {
					if err := action.Input().Check(data); err != nil {
						message := "incorrect input value: " + err.Error()
						zlog.Trace().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] " + message)
						return c.Status(DataError.HttpStatus).JSON(fiber.Map{
							"error": message,
							"type":  DataError.ErrorType,
						})
					}
				}
				// Execute the action requests
				output, err := handler(data, options)
				if err != nil {
					zlog.Error().Str("action", tdAction.Key).Err(err).Msg("[protocolHttp:actionHandler]")
					return c.Status(UnknownError.HttpStatus).JSON(fiber.Map{
						"error": err.Error(),
						"type":  UnknownError.ErrorType,
					})
				}

				// Check the output data
				if action.Output() != nil {
					if err := action.Output().Check(output); err != nil {
						zlog.Error().Str("action", tdAction.Key).Err(err).Msg("[protocolHttp:actionHandler] incorrect handler returned value")
						return c.Status(UnknownError.HttpStatus).JSON(fiber.Map{
							"error": "Incorrect handler returned value",
							"type":  UnknownError.ErrorType,
						})
					}
					zlog.Trace().Interface("response", output).Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] JSON Response to Thing action POST request")
					return c.JSON(output)
				}

				zlog.Trace().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] OK Response to Thing action POST request")
				return c.JSON(fiber.Map{"ok": true})
			} else {
				zlog.Warn().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] no handler function for the action")
				return c.Status(NotSupportedError.HttpStatus).JSON(fiber.Map{
					"error": "Not Implemented",
					"type":  NotSupportedError.ErrorType,
				})
			}
		}
	}
}
