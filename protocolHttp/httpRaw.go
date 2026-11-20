package protocolHttp

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/dataSchema"
)

// rawMediaType returns the media type a string schema declares for its
// content (image/jpeg, text/html…), or "" when the value travels as JSON
func rawMediaType(d *dataSchema.Data) string {
	if d == nil {
		return ""
	}
	s, ok := d.DataSchema.(dataSchema.String)
	if !ok || s.ContentMediaType == "" || strings.HasPrefix(s.ContentMediaType, fiber.MIMEApplicationJSON) {
		return ""
	}
	return s.ContentMediaType
}

// sendRaw serves a handler value as is: []byte (binary content) or string
// (text content such as HTML)
func sendRaw(c *fiber.Ctx, mediaType string, value interface{}) error {
	var body []byte
	switch v := value.(type) {
	case []byte:
		body = v
	case string:
		body = []byte(v)
	default:
		return fmt.Errorf("value of type %T can't be served as %s", value, mediaType)
	}
	c.Set(fiber.HeaderContentType, mediaType)
	return c.Send(body)
}
