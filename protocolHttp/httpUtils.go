package protocolHttp

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/project-eria/go-wot/interaction"
)

func decodeJSON(body io.ReadCloser) (interface{}, error) {
	var object interface{}
	err := json.NewDecoder(body).Decode(&object)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func getUri(form *interaction.Form, dataVariables map[string]interface{}) string {
	uri := form.Href
	if len(dataVariables) > 0 {
		// Insert uri variables
		for k, v := range dataVariables {
			uri = strings.Replace(uri, "{"+k+"}", v.(string), -1)
		}
	}
	return uri
}
