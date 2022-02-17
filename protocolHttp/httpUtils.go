package protocolHttp

import (
	"encoding/json"
	"io"
)

func decodeJSON(body io.ReadCloser) (interface{}, error) {
	var object interface{}
	err := json.NewDecoder(body).Decode(&object)
	if err != nil {
		return nil, err
	}
	return object, nil
}
