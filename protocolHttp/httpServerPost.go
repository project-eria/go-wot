package protocolHttp

import (
	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
)

// post handle the POST request method for a thing action
// https://w3c.github.io/wot-scripting-api/#handling-action-requests
func actionHandler(t *producer.ExposedThing, tdAction *interaction.Action) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		log.Trace().Str("uri", c.Path()).Msg("[protocolHttp:actionHandler] Received Thing action POST request")
		action := t.ExposedActions[tdAction.Key]
		handler := action.GetHandler()
		if handler != nil {
			var data interface{}
			if len(c.Body()) > 0 {
				if err := c.BodyParser(&data); err != nil {
					log.Trace().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] Incorrect JSON value")
					return c.Status(EncodingError.httpStatus).JSON(fiber.Map{
						"error": "Incorrect JSON value",
						"type":  EncodingError.errorType,
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
					return c.Status(DataError.httpStatus).JSON(fiber.Map{
						"error": message,
						"type":  DataError.errorType,
					})
				}
			}
			// Execute the action requests
			output, err := handler(data)
			if err != nil {
				log.Error().Str("action", tdAction.Key).Err(err).Msg("[protocolHttp:actionHandler]")
				return c.Status(UnknownError.httpStatus).JSON(fiber.Map{
					"error": err.Error(),
					"type":  UnknownError.errorType,
				})
			}

			// Check the output data
			if action.Output != nil {
				if err := action.Output.Check(output); err != nil {
					log.Error().Str("action", tdAction.Key).Err(err).Msg("[protocolHttp:actionHandler] incorrect handler returned value")
					return c.Status(UnknownError.httpStatus).JSON(fiber.Map{
						"error": "Incorrect handler returned value",
						"type":  UnknownError.errorType,
					})
				}
				log.Trace().Interface("response", output).Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] JSON Response to Thing action POST request")
				return c.JSON(fiber.Map{"ok": true})
			}

			log.Trace().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] OK Response to Thing action POST request")
			return c.JSON(fiber.Map{"ok": true})
		} else {
			log.Warn().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] no handler function for the action")
			return c.Status(NotSupportedError.httpStatus).JSON(fiber.Map{
				"error": "Not Implemented",
				"type":  NotSupportedError.errorType,
			})
		}
	}
}
