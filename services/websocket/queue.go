package websocket_service

import (
	"log/slog"
	"strconv"
)

type Queue struct {
	clients map[*Client]bool

	register   chan *Client
	unregister chan *Client
}

func NewQueue() *Queue {
	return &Queue{
		clients: make(map[*Client]bool),

		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (q *Queue) Run(gp *GamePool) {
	for {
		select {
		case client := <-q.register:
			slog.Info("New Client was registered to the queue.")
			q.clients[client] = true
			client.queue = q

			slog.Info("Number of clients in the queue: " + strconv.Itoa(len(q.clients)))

			if len(q.clients) < 2 {
				continue
			}

			slog.Info("Two clients are in the queue. Registering them to the game pool.")

			// get two players from the clients map
			var players []*Client
			for client := range q.clients {
				players = append(players, client)
				if len(players) == 2 {
					break
				}
			}

			delete(q.clients, players[0])
			delete(q.clients, players[1])

			pplayers := [2]*Client{players[0], players[1]}
			gp.register <- pplayers

		case client := <-q.unregister:
			slog.Info("Received unregister request.")
			if _, ok := q.clients[client]; ok {
				delete(q.clients, client)
				slog.Info("Client was unregistered from the queue.")
			}
		}
	}
}
