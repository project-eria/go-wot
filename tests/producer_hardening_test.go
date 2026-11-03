package test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
)

// SECURITY_AUDIT #11: an internal error never reaches the client; a data
// error, meant for it, does.
func TestInternalErrorHidden(t *testing.T) {
	mything := getThing()
	mything.AddAction(interaction.NewAction("broken", "", ""))
	mything.AddAction(interaction.NewAction("refused", "", ""))
	httpServer, exposed := getProducer(mything)
	exposed.SetActionHandler("broken", func(any, map[string]any) (any, error) {
		return nil, errors.New("open /var/lib/eria/secret.db: permission denied")
	})
	exposed.SetActionHandler("refused", func(any, map[string]any) (any, error) {
		return nil, &producer.DataError{Message: "unknown node \"x\""}
	})

	e := JSONTester(t, httpServer)
	e.POST("/broken").WithHeader("Content-type", "application/json").Expect().
		Status(http.StatusInternalServerError).JSON().Object().HasValue("error", "internal error")
	e.POST("/refused").WithHeader("Content-type", "application/json").Expect().
		Status(http.StatusBadRequest).JSON().Object().HasValue("error", "unknown node \"x\"")
}
