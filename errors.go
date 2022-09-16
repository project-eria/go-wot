package wot

import "errors"

// ErrNotConnected is returned when the application read/writes
// a message and the connection is closed
var (
	ErrWSNotConnected = errors.New("websocket not connected")
)
