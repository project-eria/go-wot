package protocolHttp

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
	zlog "github.com/rs/zerolog/log"
)

// post handle the POST request method for a thing action
// https://w3c.github.io/wot-scripting-api/#handling-action-requests
func actionHandler(t *producer.ExposedThing, tdAction *interaction.Action) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		log.Trace().Str("uri", c.Path()).Msg("[protocolHttp:actionHandler] Received Thing action POST request")
		if action, ok := t.ExposedActions[tdAction.Key]; ok {
			handler := action.GetHandler()
			if handler != nil {
				// Check the params (uriVariables) data
				params := c.AllParams()
				if err := action.CheckUriVariables(params); err != nil {
					return c.Status(DataError.HttpStatus).JSON(fiber.Map{
						"error": err.Error(),
						"type":  DataError.ErrorType,
					})
				}

				var data interface{}
				if len(c.Body()) > 0 {
					if err := c.BodyParser(&data); err != nil {
						log.Trace().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] Incorrect JSON value")
						return c.Status(EncodingError.HttpStatus).JSON(fiber.Map{
							"error": "Incorrect JSON value",
							"type":  EncodingError.ErrorType,
						})
					}
				} else {
					log.Trace().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] No body data")
				}

				// Check the input data
				if action.Input != nil {
					if err := action.Input.Check(data); err != nil {
						message := "incorrect input value: " + err.Error()
						log.Trace().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] " + message)
						return c.Status(DataError.HttpStatus).JSON(fiber.Map{
							"error": message,
							"type":  DataError.ErrorType,
						})
					}
				}
				// Execute the action requests
				output, err := handler(data, params)
				if err != nil {
					log.Error().Str("action", tdAction.Key).Err(err).Msg("[protocolHttp:actionHandler]")
					return c.Status(UnknownError.HttpStatus).JSON(fiber.Map{
						"error": err.Error(),
						"type":  UnknownError.ErrorType,
					})
				}

				// Check the output data
				if action.Output != nil {
					if err := action.Output.Check(output); err != nil {
						log.Error().Str("action", tdAction.Key).Err(err).Msg("[protocolHttp:actionHandler] incorrect handler returned value")
						return c.Status(UnknownError.HttpStatus).JSON(fiber.Map{
							"error": "Incorrect handler returned value",
							"type":  UnknownError.ErrorType,
						})
					}
					log.Trace().Interface("response", output).Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] JSON Response to Thing action POST request")
					return c.JSON(fiber.Map{"ok": true})
				}

				log.Trace().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] OK Response to Thing action POST request")
				return c.JSON(fiber.Map{"ok": true})
			} else {
				log.Warn().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] no handler function for the action")
				return c.Status(NotSupportedError.HttpStatus).JSON(fiber.Map{
					"error": "Not Implemented",
					"type":  NotSupportedError.ErrorType,
				})
			}
		} else {
			zlog.Error().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] ExposedAction not found")
			return c.Status(UnknownError.HttpStatus).JSON(fiber.Map{
				"error": fmt.Errorf("ExposedAction `%s` not found", tdAction.Key),
				"type":  UnknownError.ErrorType,
			})
		}
	}
}
