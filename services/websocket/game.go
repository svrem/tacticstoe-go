package websocket_service

import "log/slog"

type Game struct {
	player1 *Client
	player2 *Client

	unregister chan *Client
}

func newGame(player1 *Client, player2 *Client) Game {
	player1.send <- []byte("Game started.")
	player2.send <- []byte("Game started.")

	player1.queue = nil
	player2.queue = nil

	return Game{
		player1: player1,
		player2: player2,

		unregister: make(chan *Client),
	}
}

func (g *Game) Run(gp *GamePool) {
	for {
		select {
		case client := <-g.unregister:
			slog.Info("Unregistering client from the game.")
			if client == g.player1 {
				g.player2.game = nil
				g.player2.send <- []byte("Opponent left the game.")
			} else {
				g.player1.game = nil
				g.player1.send <- []byte("Opponent left the game.")
			}

			gp.close_game <- g
		}
	}
}
