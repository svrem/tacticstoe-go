package websocket_service

import "log/slog"

type GamePool struct {
	games []Game

	register chan [2]*Client

	closeGame chan *Game
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

			newGame := newGame(clients[0], clients[1])
			gp.games = append(gp.games, newGame)

			clients[0].game = &newGame
			clients[1].game = &newGame

			go newGame.Run(gp)

		case game := <-gp.closeGame:
			for i, g := range gp.games {
				if g.id == game.id {
					slog.Info("Closing game.")

					game.player1.game = nil
					game.player2.game = nil

					close(game.player1.send)
					close(game.player2.send)

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

		closeGame: make(chan *Game),
	}
}
