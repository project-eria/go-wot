package action

import (
	"sync"
)

// Action represents an individual action on a thing.
type Action struct {
	ID     string
	ATtype string
	//	Thing   *Thing
	//	execute func(*ActionRequest)
	mu sync.RWMutex
	//	ActionMeta
}

// // ActionMeta represents additional attribute for the action
// type ActionMeta struct {
// 	Title       string
// 	Description string
// 	Input       ActionInputMeta
// }

// // ActionInputMeta represents additional attribute for the action inputs
// type ActionInputMeta struct {
// 	Type       string
// 	Properties map[string]ActionInputPropertyMeta
// }

// // ActionInputPropertyMeta represents additional attribute for the action inputs properties
// type ActionInputPropertyMeta struct {
// 	Type    string
// 	Minimum int
// 	Maximum int
// 	Unit    string
// }

// // AddAction add action to a thing
// func (t *Thing) AddAction(id string, meta ActionMeta) *Action {
// 	if t == nil {
// 		log.Error().Msg("[Thing:AddAction] nil thing")
// 		return nil
// 	}
// 	action := &Action{
// 		ID:         id,
// 		Thing:      t,
// 		ActionMeta: meta,
// 	}

// 	t.mu.Lock()
// 	defer t.mu.Unlock()
// 	t.actions[id] = action
// 	return action
// }

// // GetAction get an existing thing action
// func (t *Thing) GetAction(id string) *Action {
// 	if t == nil {
// 		log.Error().Msg("[Thing:GetAction] nil thing")
// 		return nil
// 	}
// 	t.mu.RLock()
// 	defer t.mu.RUnlock()

// 	if action, ok := t.actions[id]; ok {
// 		return action
// 	}
// 	log.Error().Msg("[Thing:GetAction] action not found")
// 	return nil
// }

// // AddExecute add a handler function that is executed on action request
// func (a *Action) AddExecute(execute func(*ActionRequest)) {
// 	a.mu.RLock()
// 	defer a.mu.RUnlock()
// 	a.execute = execute
// }

// // description returns the action details
// func (a *Action) description() map[string]interface{} {
// 	if a == nil {
// 		log.Error().Msg("[thingAction:description] nil action")
// 		return nil
// 	}
// 	result := make(map[string]interface{})
// 	a.mu.Lock()
// 	result["title"] = a.Title
// 	result["description"] = a.Description
// 	a.mu.Unlock()
// 	result["forms"] = [1]map[string]string{{"href": a.href()}}
// 	return result
// }

// func (a *Action) href() string {
// 	if a == nil {
// 		log.Error().Msg("[thingAction:description] nil action")
// 		return ""
// 	}
// 	a.mu.RLock()
// 	defer a.mu.RUnlock()
// 	return a.Thing.href() + "actions/" + a.ID + "/"
// }
