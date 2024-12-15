package websocket_service

import (
	"log/slog"
	"time"
)

type GamePool struct {
	games []*Game

	register   chan [2]*Client
	unregister chan *Client
}

func (gp *GamePool) Run() {
	for {
		select {
		case players := <-gp.register:
			if len(players) != 2 {
				slog.Error("Invalid number of clients to register to the game pool.")
				continue
			}

			slog.Info("Registering new Game")

			newGame := newGame(players[0], players[1])
			gp.games = append(gp.games, newGame)

			players[0].game = newGame
			players[1].game = newGame

			go newGame.Run(gp)

		case client := <-gp.unregister:
			for i, g := range gp.games {
				if g.player1 == client || g.player2 == client {
					if !g.isOver {

						gameAbortedMessageString, err := generateGameAbortedMessageString()

						if err != nil {
							slog.Error("Error generating game aborted message string: " + err.Error())
							continue
						}

						if client != g.player1 {
							g.player1.send <- []byte(gameAbortedMessageString)

							time.Sleep(500 * time.Millisecond)

							g.player1.hub.unregister <- g.player1
						} else {
							g.player2.send <- []byte(gameAbortedMessageString)

							time.Sleep(500 * time.Millisecond)

							g.player2.hub.unregister <- g.player2
						}
					}

					gp.games = append(gp.games[:i], gp.games[i+1:]...)

					break
				}
			}
		}
	}
}

func generateGameAbortedMessageString() ([]byte, error) {
	type GameAbort struct {
	}

	gameAbortedMessage := DataMessage[GameAbort]{
		Type: "game_abort",
		Data: GameAbort{},
	}

	return gameAbortedMessage.Marshal()

}

func NewGamePool() *GamePool {
	return &GamePool{
		games: make([]*Game, 0),

		register:   make(chan [2]*Client),
		unregister: make(chan *Client),
	}
}
