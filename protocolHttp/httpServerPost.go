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
		optionsStr := c.AllParams()
		zlog.Trace().Str("uri", c.Path()).Interface("options", optionsStr).Msg("[protocolHttp:actionHandler] Received Thing action POST request")
		action, err := t.ExposedAction(tdAction.Key)
		if err != nil {
			zlog.Error().Err(err).Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler]")
			return c.Status(UnknownError.HttpStatus).JSON(fiber.Map{
				"error": fmt.Sprintf("ExposedAction `%s` not found", tdAction.Key),
				"type":  UnknownError.ErrorType,
			})
		}
		var input interface{}
		if len(c.Body()) > 0 {
			if err := c.BodyParser(&input); err != nil {
				zlog.Trace().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] Incorrect JSON value")
				return c.Status(EncodingError.HttpStatus).JSON(fiber.Map{
					"error": "Incorrect JSON value",
					"type":  EncodingError.ErrorType,
				})
			}
		} else {
			zlog.Trace().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] No body data")
		}
		output, err := action.Run(t, tdAction.Key, input, optionsStr)
		if err != nil {
			zlog.Error().Err(err).Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler]")

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
		if output != nil {
			zlog.Trace().Interface("response", output).Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] JSON Response to Thing action POST request")
			return c.JSON(output)
		}
		zlog.Trace().Str("action", tdAction.Key).Msg("[protocolHttp:actionHandler] OK Response to Thing action POST request")
		return c.JSON(fiber.Map{"ok": true})
	}
}
