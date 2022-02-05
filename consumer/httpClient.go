package consumer

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// getHTTPJSON get a JSON data from HTTP GET request
func getHTTPJSON(url string) (interface{}, error) {
	data, err := sendHTTPJSON(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// putHTTPJSON send JSON data using HTTP PUT request
func putHTTPJSON(url string, payload interface{}) (interface{}, error) {
	data, err := sendHTTPJSON(url, http.MethodPut, payload)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// postHTTPJSON send JSON data using HTTP POST request
func postHTTPJSON(url string, payload interface{}) (interface{}, error) {
	data, err := sendHTTPJSON(url, http.MethodPost, payload)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func sendHTTPJSON(url string, method string, payload interface{}) (interface{}, error) {
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
	resp, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := decodeJSON(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		dataMap := data.(map[string]interface{})
		if msg, hasError := dataMap["error"]; hasError {
			log.Error().Str("status", resp.Status).Str("url", url).Str("msg", msg.(string)).Msg("[consumer:sendHTTPJSON]")
			return data, errors.New(msg.(string))
		}
		log.Error().Str("status", resp.Status).Str("url", url).Msg("[consumer:sendHTTPJSON] request returned error")
		return data, errors.New("request returned error")
	}
	return data, nil
}

func decodeJSON(r *http.Response) (interface{}, error) {
	var object interface{}
	err := json.NewDecoder(r.Body).Decode(&object)
	if err != nil {
		return nil, err
	}
	return object, nil
}
