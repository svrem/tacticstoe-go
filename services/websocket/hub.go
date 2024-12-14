package websocket_service

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
			println("Registering client")
			h.clients[client] = true
			client.hub = h
			queue.register <- client

		case client := <-h.unregister:
			println("Unregistering client")
			if _, ok := h.clients[client]; ok {
				gamePool.unregister <- client
				println("Unregistering client from game pool")
				queue.unregister <- client
				println("Unregistering client from queue")

				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}
