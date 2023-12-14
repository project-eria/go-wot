package consumer

import (
	"fmt"

	"github.com/project-eria/go-wot/dataSchema"
)

func getUriVariables(uriVariables map[string]dataSchema.Data, dataVariables map[string]interface{}) (map[string]interface{}, error) {
	// check if all uri variables are set
	data := make(map[string]interface{})
	if len(uriVariables) > 0 {
		if len(dataVariables) == 0 {
			return nil, fmt.Errorf("uri variables not set")
		}
		for k, v := range uriVariables {
			if value, ok := dataVariables[k]; ok {
				data[k] = value
			} else if v.Default != "" {
				data[k] = v.Default
			} else {
				return nil, fmt.Errorf("uri variable %s not set", k)
			}
		}
	}
	return data, nil
}
