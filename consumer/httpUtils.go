package consumer

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
)

// GetPropertyValue returns the value of a remote thing property
// func (t *ThingConnection) GetPropertyValue(property string) interface{} {
// 	if t == nil {
// 		log.Error().Msg("[consumer:GetPropertyValue] nil connection")
// 		return nil
// 	}

// 	urlHTTP := t.getHTTPURL("/properties/" + property)
// 	data, err := getHTTPJSON(urlHTTP)
// 	if err != nil {
// 		log.Error().Err(err).Msg("[consumer:GetPropertyValue]")
// 		return nil
// 	}
// 	value, ok := data[property]
// 	if !ok {
// 		log.Error().Msg("[consumer:GetPropertyValue] property name not found in JSON response")
// 		return nil
// 	}
// 	return value
// }

// PostActionRequest send an action request to a remote thing action
// func (t *ThingConnection) PostActionRequest(action string, input map[string]interface{}) interface{} {
// 	if t == nil {
// 		log.Error().Msg("[consumer:PostActionRequest] nil connection")
// 		return nil
// 	}

// 	urlHTTP := t.getHTTPURL("/actions/" + action)
// 	parameters := make(map[string]interface{})

// 	if len(input) > 0 {
// 		parameters["input"] = input
// 	}

// 	payload, err := json.Marshal(map[string]interface{}{
// 		action: parameters,
// 	})

// 	if err != nil {
// 		log.Error().Err(err).Msg("[consumer:PostActionRequest]")
// 		return nil
// 	}

// 	data, err := postHTTPJSON(urlHTTP, payload)
// 	if err != nil {
// 		log.Error().Err(err).Msg("[consumer:PostActionRequest]")
// 		return nil
// 	}

// 	value, ok := data[action]
// 	if !ok {
// 		log.Error().Msg("[consumer:PostActionRequest] action name not found in JSON response")
// 		return nil
// 	}
// 	return value
// }

// getHTTPJSON get a JSON data from HTTP GET request
func getHTTPJSON(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error().Str("status", resp.Status).Str("url", url).Msg("[consumer:getHTTPJSON] incorrect response")
		return nil, errors.New("incorrect HTTP response")
	}
	data, err := decodeJSON(resp)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// postHTTPJSON send JSON data using HTTP POST request
// func postHTTPJSON(url string, payload []byte) (map[string]interface{}, error) {
// 	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusCreated {
// 		log.Error().Str("status", resp.Status).Str("url", url).Msg("[consumer:postHTTPJSON] incorrect response")
// 		return nil, errors.New("Incorrect HTTP response")
// 	}

// 	data, err := decodeJSON(resp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }

func decodeJSON(r *http.Response) (interface{}, error) {
	var object interface{}
	err := json.NewDecoder(r.Body).Decode(&object)
	if err != nil {
		return nil, err
	}
	return object, nil
}
