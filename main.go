package main

import (
	"log/slog"
	"net/http"

	"tacticstoe/internal/auth"
	db "tacticstoe/internal/database"
	"tacticstoe/internal/websocket"

	"github.com/joho/godotenv"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/index.html")
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Warn("No .env file found")
	}

	// Utils
	database := db.OpenDatabase()
	db.MigrateModel(database)

	hub := websocket.NewHub()
	go hub.Run(database)

	// Routes
	http.HandleFunc("/", getRoot)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r, database)
	})

	http.HandleFunc("GET /auth/me/", func(w http.ResponseWriter, r *http.Request) {
		auth.MeHandler(w, r, database)
	})
	http.HandleFunc("GET /auth/login/{provider}/", auth.LoginHandler)
	http.HandleFunc("GET /auth/logout/", auth.LogoutHandler)
	http.HandleFunc("GET /auth/callback/{provider}/", func(w http.ResponseWriter, r *http.Request) {
		auth.CallbackHandler(w, r, database)
	})

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	slog.Info("Server started at port 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		slog.Error("ListenAndServe: \"" + err.Error() + "\"")
	}

	slog.Info("Server stopped")
}
