package websocket_service

import "log/slog"

type Hub struct {
	clients map[*Client]bool

	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),

		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {

	gamePool := NewGamePool()
	go gamePool.Run()

	queue := NewQueue()
	go queue.Run(gamePool)

	for {
		select {
		case client := <-h.register:
			client.hub = h

			clientAlreadyRegistered := false
			for c := range h.clients {
				if c.id == client.id {
					slog.Info("Client already registered, ignoring...")

					close(client.send)

					clientAlreadyRegistered = true
					break
				}
			}
			if clientAlreadyRegistered {
				continue
			}

			h.clients[client] = true
			queue.register <- client

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				gamePool.unregister <- client
				queue.unregister <- client

				delete(h.clients, client)
				close(client.send)

				client.conn.Close()
			}
		}
	}
}