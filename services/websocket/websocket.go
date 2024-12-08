package websocket_service

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(queue Queue, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("wsHandler: " + err.Error())
		return
	}

	client := &Client{
		conn:       conn,
		elo_rating: 1000,
		send:       make(chan []byte, 256),
	}
	queue.register <- client

	go client.writePump()
	go client.readPump()

	// for {
	// 	messageType, data, err := conn.ReadMessage()
	// 	if err != nil {
	// 		slog.Error("wsHandler: " + err.Error())
	// 		return
	// 	}

	// 	err = conn.WriteMessage(messageType, data)
	// 	if err != nil {
	// 		slog.Error("wsHandler: " + err.ha
	// 	}
	// }
}
