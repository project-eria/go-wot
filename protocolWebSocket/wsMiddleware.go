package protocolWebSocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func checkUpgrade() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			key := c.Get("Sec-Websocket-Key")
			c.Locals("key", key)
			c.Locals("websocket", true) // Add a 'websocket' flag for the http handlers
		}
		return c.Next()
	}
}
