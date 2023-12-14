package protocolHttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/interaction"
	zlog "github.com/rs/zerolog/log"
)

type HttpClient struct {
	Client  HttpClientProcessor
	Schemes []string
}

func NewClient() *HttpClient {
	return &HttpClient{
		Client:  &http.Client{},
		Schemes: []string{"http"},
	}
}

// HttpClientProcessor interface
type HttpClientProcessor interface {
	Do(req *http.Request) (*http.Response, error)
}

func (c *HttpClient) GetSchemes() []string {
	return c.Schemes
}

// ReadResource get a JSON data from HTTP GET request
func (c *HttpClient) ReadResource(form *interaction.Form, dataVariables map[string]interface{}) (interface{}, string, error) {
	uri := getUri(form, dataVariables)
	data, err := c.sendJSON(uri, http.MethodGet, nil)
	if err != nil {
		return nil, uri, err
	}
	return data, uri, nil
}

// WriteResource send JSON data using HTTP PUT request
func (c *HttpClient) WriteResource(form *interaction.Form, dataVariables map[string]interface{}, value interface{}) (interface{}, string, error) {
	uri := getUri(form, dataVariables)
	data, err := c.sendJSON(uri, http.MethodPut, value)
	if err != nil {
		return nil, uri, err
	}
	return data, uri, nil
}

// InvokeResource send JSON data using HTTP POST request
func (c *HttpClient) InvokeResource(form *interaction.Form, dataVariables map[string]interface{}, value interface{}) (interface{}, string, error) {
	uri := getUri(form, dataVariables)
	data, err := c.sendJSON(uri, http.MethodPost, value)
	if err != nil {
		return nil, uri, err
	}
	return data, uri, nil
}

func (c *HttpClient) SubscribeResource(form *interaction.Form, dataVariables map[string]interface{}, sub *consumer.Subscription, listener consumer.Listener) (string, error) {
	return form.Href, errors.New("not implemented")
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
			zlog.Error().Str("status", resp.Status).Str("url", url).Str("msg", msg.(string)).Msg("[consumer:sendJSON]")
			return data, errors.New(msg.(string))
		}
		zlog.Error().Str("status", resp.Status).Str("url", url).Msg("[consumer:sendJSON] request returned error")
		return data, errors.New("request returned error")
	}
	return data, nil
}
