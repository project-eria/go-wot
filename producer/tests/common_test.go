package test

import (
	"os"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/project-eria/go-wot/securityScheme"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var base_url = "http://127.0.0.1:8888"

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func checkGet(t *testing.T, client *resty.Client, path string, code int, body string) {
	resp, _ := client.R().
		EnableTrace().
		Get(base_url + path)
	assert.Equal(t, code, resp.StatusCode())
	assert.Equal(t, body, resp.String())
}

func checkPut(t *testing.T, client *resty.Client, path string, value string, code int, body string) {
	resp, _ := client.R().
		EnableTrace().
		SetBody(value).
		Put(base_url + path)
	assert.Equal(t, code, resp.StatusCode())
	assert.Equal(t, body, resp.String())
}

func checkPost(t *testing.T, client *resty.Client, path string, value string, code int, body string) {
	resp, _ := client.R().
		EnableTrace().
		SetBody(value).
		Post(base_url + path)
	assert.Equal(t, code, resp.StatusCode())
	assert.Equal(t, body, resp.String())
}

func createClient() *resty.Client {
	client := resty.New()
	client.
		// Set retry count to non zero to enable retries
		SetRetryCount(3).
		// You can override initial retry wait time.
		// Default is 100 milliseconds.
		SetRetryWaitTime(1*time.Second).
		// MaxWaitTime can be overridden as well.
		// Default is 2 seconds.
		SetRetryMaxWaitTime(20*time.Second).
		SetHeader("Content-Type", "application/json")
	return client
}

func createThing() *thing.Thing {
	mything, err := thing.New(
		"dev:ops:my-actuator-1234",
		"Actuator1 Example",
		"An actuator 1st example",
		[]string{},
	)

	if err != nil {
		os.Exit(1)
	}
	// Add Security
	noSecurityScheme := securityScheme.NewNoSecurity()
	mything.AddSecurity("no_sec", noSecurityScheme)

	return mything
}
