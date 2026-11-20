package test

import (
	"net/http"
	"testing"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
)

var jpeg = []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 'J', 'F', 'I', 'F'}

// Properties and actions declaring a non-JSON contentMediaType are served raw
func TestRawContent(t *testing.T) {
	mything := getThing()
	imageData, _ := dataSchema.NewString(
		dataSchema.StringContentMediaType("image/jpeg"),
		dataSchema.StringContentEncoding("base64"),
	)
	mything.AddProperty(interaction.NewProperty("image", "Image", "", imageData, interaction.PropertyReadOnly(true)))
	mything.AddProperty(interaction.NewProperty("imageRW", "Writable image", "", imageData))
	htmlData, _ := dataSchema.NewString(dataSchema.StringContentMediaType("text/html; charset=utf-8"))
	mything.AddProperty(interaction.NewProperty("page", "Page", "", htmlData, interaction.PropertyReadOnly(true)))
	inputData, _ := dataSchema.NewObject()
	mything.AddAction(interaction.NewAction("snapshot", "Snapshot", "",
		interaction.ActionInput(&inputData),
		interaction.ActionOutput(&imageData),
	))
	mything.AddAction(interaction.NewAction("broken", "Broken", "", interaction.ActionOutput(&imageData)))

	httpServer, exposedThing := getProducer(mything)
	image := func(producer.ExposedThing, string, map[string]interface{}) (interface{}, error) {
		return jpeg, nil
	}
	exposedThing.SetPropertyReadHandler("image", image)
	exposedThing.SetPropertyReadHandler("imageRW", image)
	exposedThing.SetPropertyReadHandler("page", func(producer.ExposedThing, string, map[string]interface{}) (interface{}, error) {
		return "<p>Bonjour</p>", nil
	})
	exposedThing.SetActionHandler("snapshot", func(interface{}, map[string]interface{}) (interface{}, error) {
		return jpeg, nil
	})
	exposedThing.SetActionHandler("broken", func(interface{}, map[string]interface{}) (interface{}, error) {
		return 42, nil
	})

	e := fastHTTPTester(t, httpServer)

	t.Run("binary property", func(t *testing.T) {
		resp := e.GET("/image").Expect().Status(http.StatusOK)
		resp.Header("Content-Type").IsEqual("image/jpeg")
		resp.Body().IsEqual(string(jpeg))
	})

	t.Run("html property", func(t *testing.T) {
		resp := e.GET("/page").Expect().Status(http.StatusOK)
		resp.Header("Content-Type").IsEqual("text/html; charset=utf-8")
		resp.Body().IsEqual("<p>Bonjour</p>")
	})

	t.Run("binary action output", func(t *testing.T) {
		resp := e.POST("/snapshot").WithHeader("Content-type", "application/json").
			WithBytes([]byte(`{"stream":"sub"}`)).
			Expect().Status(http.StatusOK)
		resp.Header("Content-Type").IsEqual("image/jpeg")
		resp.Body().IsEqual(string(jpeg))
	})

	t.Run("unservable handler value", func(t *testing.T) {
		e.POST("/broken").WithHeader("Content-type", "application/json").
			Expect().Status(http.StatusInternalServerError).
			JSON().Object().HasValue("type", "UnknownError")
	})

	t.Run("thing description", func(t *testing.T) {
		td := e.GET("/").Expect().Status(http.StatusOK).JSON().Object()
		image := td.Path("$.properties.image").Object()
		image.HasValue("contentMediaType", "image/jpeg")
		image.HasValue("contentEncoding", "base64")
		image.Path("$.forms[0].contentType").IsEqual("image/jpeg")
		// Read raw, write JSON (base64 string)
		imageRW := td.Path("$.properties.imageRW.forms").Array()
		imageRW.Value(0).Object().HasValue("contentType", "image/jpeg").HasValue("op", []string{"readproperty"})
		imageRW.Value(1).Object().HasValue("contentType", "application/json").HasValue("op", []string{"writeproperty"})
		// JSON input, raw output
		snapshot := td.Path("$.actions.snapshot.forms[0]").Object()
		snapshot.HasValue("contentType", "application/json")
		snapshot.Path("$.response.contentType").IsEqual("image/jpeg")
	})
}
