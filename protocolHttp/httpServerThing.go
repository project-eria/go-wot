package protocolHttp

import (
	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog/log"
)

// get handle the GET method for single thing root
func thingHandler(t producer.ExposedThing) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Locals("websocket") == true {
			return c.Next() // delegate to the next handler (websocket)
		}
		log.Trace().Str("uri", c.Path()).Msg("[protocolHttp:thingHandler] Received Thing GET request")
		td := t.GetThingDescription()

		// Dynamically build href
		for _, property := range td.Properties {
			property := property // Copy https://go.dev/doc/faq#closures_and_goroutines
			for _, form := range property.Forms {
				form := form // Copy https://go.dev/doc/faq#closures_and_goroutines
				if form.UrlBuilder != nil {
					form.Href = form.UrlBuilder(c.Hostname(), c.Secure())
				}
			}
		}
		for _, action := range td.Actions {
			action := action // Copy https://go.dev/doc/faq#closures_and_goroutines
			for _, form := range action.Forms {
				form := form // Copy https://go.dev/doc/faq#closures_and_goroutines
				if form.UrlBuilder != nil {
					form.Href = form.UrlBuilder(c.Hostname(), c.Secure())
				}
			}
		}
		for _, event := range td.Events {
			event := event // Copy https://go.dev/doc/faq#closures_and_goroutines
			for _, form := range event.Forms {
				form := form // Copy https://go.dev/doc/faq#closures_and_goroutines
				if form.UrlBuilder != nil {
					form.Href = form.UrlBuilder(c.Hostname(), c.Secure())
				}
			}
		}
		return c.JSON(td)
	}
}
