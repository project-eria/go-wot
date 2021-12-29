package action

//	uuid "github.com/satori/go.uuid"

// // ActionRequest represents an request on a thing action.
// type ActionRequest struct {
// 	id     string
// 	status string
// 	input  map[string]interface{}
// 	*Action
// }

// func (a *Action) newRequest(input map[string]interface{}) *ActionRequest {
// 	id := uuid.NewV4().String()
// 	return &ActionRequest{
// 		id:     id,
// 		input:  input,
// 		status: "pending",
// 		Action: a,
// 	}
// }

// func (r *ActionRequest) href() string {
// 	return r.Action.href() + r.id
// }

// func (r *ActionRequest) description() map[string]interface{} {
// 	return map[string]interface{}{
// 		"input":  r.input,
// 		"status": r.status,
// 		"href":   r.href(),
// 	}
// }

// // GetInputValue returns the value for a specific request input parameter
// func (r *ActionRequest) GetInputValue(name string) (interface{}, error) {
// 	if value, in := r.input[name]; in {
// 		return value, nil
// 	}
// 	return nil, errors.New("input name not found")
// }
