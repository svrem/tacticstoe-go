package websocket_service

import (
	"log/slog"
	"tacticstoe/services/rating"
	"time"

	"gorm.io/gorm"
)

type GamePool struct {
	games []*Game

	register   chan [2]*Client
	unregister chan *Client
}

func (gp *GamePool) Run(database *gorm.DB) {
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

			go newGame.Run(gp, database)

		case client := <-gp.unregister:
			for i, g := range gp.games {
				isPlayer1 := g.player1 == client
				isPlayer2 := g.player2 == client
				if isPlayer1 || isPlayer2 {
					if !g.isOver {
						var newEloRating int
						if isPlayer1 {
							_, newEloRating = rating.UpdateRatings(database, g.player1.user, g.player2.user, 0)
						} else {
							newEloRating, _ = rating.UpdateRatings(database, g.player1.user, g.player2.user, 1)
						}

						gameAbortedMessageString, err := generateGameAbortedMessageString(newEloRating)

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
					} else {
						if client != g.player1 {
							g.player1.hub.unregister <- g.player1
						} else {
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

func generateGameAbortedMessageString(newEloRating int) ([]byte, error) {
	gameAbortedMessage := DataMessage[GameEnd]{
		Type: "game_end",
		Data: GameEnd{
			Winner:       "aborted",
			Coords:       nil,
			NewEloRating: newEloRating,
		},
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
