package protocolHttp

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	zlog "github.com/rs/zerolog/log"
)

func checkContentType() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if (c.Method() == fiber.MethodPut || c.Method() == fiber.MethodPost) && c.Get("Content-Type") != fiber.MIMEApplicationJSON {
			zlog.Error().Msg("[protocolHttp:checkContentType] Request without Content-type='application/json'")
			// break here instead of continuing the chain
			return c.Status(EncodingError.HttpStatus).JSON(fiber.Map{
				"error": "Content-type 'application/json' is required",
				"type":  EncodingError.ErrorType,
			})
		}
		return c.Next()
	}
}

func corsHeader() func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowHeaders: "Origin, X-Requested-With, Content-Type, Accept",
	})
}
