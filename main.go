package main

import (
	"log/slog"
	"net/http"

	ws_service "tacticstoe/services/websocket"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/index.html")
}

func main() {
	game_pool := ws_service.NewGamePool()
	go game_pool.Run()

	queue := ws_service.NewQueue()
	go queue.Run(game_pool)

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws_service.ServeWs(*queue, w, r)
	})
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	slog.Info("Server started at port 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		slog.Error("ListenAndServe: \"" + err.Error() + "\"")
	}

	slog.Info("Server stopped")
}
