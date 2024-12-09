package websocket_service

import "log/slog"

type GamePool struct {
	games []Game

	register chan [2]*Client

	close_game chan *Game
}

func (gp *GamePool) Run() {
	for {
		select {
		case clients := <-gp.register:
			if len(clients) != 2 {
				slog.Error("Invalid number of clients to register to the game pool.")
				continue
			}

			slog.Info("Registering new Game")

			new_game := newGame(clients[0], clients[1])
			gp.games = append(gp.games, new_game)

			clients[0].game = &new_game
			clients[1].game = &new_game

			go new_game.Run(gp)

		case game := <-gp.close_game:
			for i, g := range gp.games {
				if g.id == game.id {
					slog.Info("Closing game.")

					g.player1.game = nil
					g.player2.game = nil

					// g.player1.conn.Close()
					// g.player2.conn.Close()

					gp.games = append(gp.games[:i], gp.games[i+1:]...)
					break
				}
			}
		}
	}
}

func NewGamePool() *GamePool {
	return &GamePool{
		games: make([]Game, 0),

		register: make(chan [2]*Client),

		close_game: make(chan *Game),
	}
}
