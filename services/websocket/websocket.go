package websocket_service

import (
	"log/slog"
	"net/http"
	auth_service "tacticstoe/services/auth"

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
	user := auth_service.AutherizeUser(w, r, database)

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
		conn:      conn,
		eloRating: 1000,
		send:      make(chan []byte, 256),
		id:        user.ID,
	}
	hub.register <- client

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
