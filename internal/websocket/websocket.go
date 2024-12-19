package websocket

import (
	"log/slog"
	"net/http"
	"tacticstoe/internal/auth"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type JoinData struct {
	Id string `json:"id"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, database *gorm.DB) {
	slog.Info("wsHandler: ServeWs")
	user := auth.AutherizeUser(w, r, database)

	if user == nil {
		slog.Info("wsHandler: Unauthorized")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("wsHandler: " + err.Error())
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),

		user: user,
	}
	hub.register <- client

	joinData := JoinData{
		Id: client.user.ID,
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
