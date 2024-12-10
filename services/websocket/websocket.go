package websocket_service

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type JoinData struct {
	Id string `json:"id"`
}

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
		conn:      conn,
		eloRating: 1000,
		send:      make(chan []byte, 256),
		id:        uuid.New().String(),
	}
	queue.register <- client

	joinData := JoinData{
		Id: client.id,
	}
	joinMessage := DataMessage[JoinData]{
		Type: "join",
		Data: joinData,
	}

	joinMessageString, err := joinMessage.Marshal()

	if err != nil {
		slog.Error("wsHandler: " + err.Error())
		return
	}

	client.send <- joinMessageString

	go client.writePump()
	go client.readPump()
}
