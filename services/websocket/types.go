package websocket_service

type DataMessage[T any] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
}
