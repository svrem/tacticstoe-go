package main

import (
	"net/http"
	ws_service "tacticstoe/services/websocket"
	"testing"

	"github.com/gorilla/websocket"
)

func TestLoadTestGame(t *testing.T) {
	go startWebsocketServer()

	for i := 0; i < 1000; i++ {
		c := connectToWebsocket()

		if c == nil {
			t.Fail()
		}
	}
}

func connectToWebsocket() *websocket.Conn {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)

	if err != nil {
		return nil
	}

	return c
}

func startWebsocketServer() {
	game_pool := ws_service.NewGamePool()
	go game_pool.Run()

	queue := ws_service.NewQueue()
	go queue.Run(game_pool)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws_service.ServeWs(*queue, w, r)
	})

	http.ListenAndServe(":8080", nil)
}
