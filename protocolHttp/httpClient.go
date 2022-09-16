package protocolHttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/interaction"
	"github.com/rs/zerolog/log"
)

type HttpClient struct {
	Client  HttpClientProcessor
	schemes []string
}

func NewClient() *HttpClient {
	return &HttpClient{
		Client:  &http.Client{},
		schemes: []string{"http"},
	}
}

// HttpClientProcessor interface
type HttpClientProcessor interface {
	Do(req *http.Request) (*http.Response, error)
}

func (c *HttpClient) GetSchemes() []string {
	return c.schemes
}

// ReadResource get a JSON data from HTTP GET request
func (c *HttpClient) ReadResource(form *interaction.Form) (interface{}, error) {
	data, err := c.sendJSON(form.Href, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WriteResource send JSON data using HTTP PUT request
func (c *HttpClient) WriteResource(form *interaction.Form, value interface{}) (interface{}, error) {
	data, err := c.sendJSON(form.Href, http.MethodPut, value)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// InvokeResource send JSON data using HTTP POST request
func (c *HttpClient) InvokeResource(form *interaction.Form, value interface{}) (interface{}, error) {
	data, err := c.sendJSON(form.Href, http.MethodPost, value)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *HttpClient) SubscribeResource(form *interaction.Form, sub *consumer.Subscription, listener consumer.Listener) error {
	return errors.New("not implemented")
}

func (c *HttpClient) Stop() {
	// Nothing to do
}

func (c *HttpClient) sendJSON(url string, method string, payload interface{}) (interface{}, error) {
	var buffer io.Reader
	if payload != nil {
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		buffer = bytes.NewBuffer(jsonBytes)
	}
	req, err := http.NewRequest(method, url, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")

	// Do sends an HTTP request and
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := decodeJSON(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		dataMap := data.(map[string]interface{})
		if msg, hasError := dataMap["error"]; hasError {
			log.Error().Str("status", resp.Status).Str("url", url).Str("msg", msg.(string)).Msg("[consumer:sendJSON]")
			return data, errors.New(msg.(string))
		}
		log.Error().Str("status", resp.Status).Str("url", url).Msg("[consumer:sendJSON] request returned error")
		return data, errors.New("request returned error")
	}
	return data, nil
}
