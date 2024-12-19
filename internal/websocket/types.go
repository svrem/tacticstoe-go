package websocket

import "encoding/json"

type DataMessage[T any] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
}

func (dm *DataMessage[any]) Marshal() ([]byte, error) {
	return json.Marshal(dm)
}
