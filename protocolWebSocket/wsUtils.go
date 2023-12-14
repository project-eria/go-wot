package protocolWebSocket

import (
	"strings"

	"github.com/project-eria/go-wot/interaction"
)

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
