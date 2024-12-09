package websocket_service

import (
	"encoding/json"
	"log/slog"
)

type GameStart struct {
	Starting_player string `json:"starting_player"`
}

type GameUpdate struct {
	X            int    `json:"x"`
	Y            int    `json:"y"`
	State        int    `json:"state"`
	ActivePlayer string `json:"active_player"`
}

type GameAction struct {
	player *Client
	x      int
	y      int
}

type Game struct {
	player1 *Client
	player2 *Client

	board [4][4]int

	active_player *Client

	make_action chan GameAction
	unregister  chan *Client
}

func newGame(player1 *Client, player2 *Client) Game {
	game_join_data_message := DataMessage[GameStart]{
		Type: "game_start",
		Data: GameStart{
			Starting_player: player1.id,
		},
	}
	data, err := json.Marshal(game_join_data_message)
	if err != nil {
		slog.Error("Error marshalling game start data.")
		return Game{}
	}

	player1.send <- []byte(data)
	player2.send <- []byte(data)

	player1.queue = nil
	player2.queue = nil

	return Game{
		player1: player1,
		player2: player2,

		board: [4][4]int{},

		active_player: player1,

		make_action: make(chan GameAction),
		unregister:  make(chan *Client),
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

		case action := <-g.make_action:
			if action.player != g.active_player {
				continue
			}

			if g.board[action.x][action.y] != 0 {
				continue
			}

			is_active_player_player1 := action.player == g.player1

			var new_state int
			if is_active_player_player1 {
				new_state = 1
			} else {
				new_state = 2
			}

			g.board[action.x][action.y] = new_state

			// Switch active player
			if is_active_player_player1 {
				g.active_player = g.player2
			} else {
				g.active_player = g.player1
			}

			gameUpdateMessage := DataMessage[GameUpdate]{
				Type: "game_update",
				Data: GameUpdate{
					X:            action.x,
					Y:            action.y,
					State:        new_state,
					ActivePlayer: g.active_player.id,
				}}

			dataUpdateMessageString, err := json.Marshal(gameUpdateMessage)

			if err != nil {
				slog.Error("Error marshalling game update data.")
				continue
			}

			g.player1.send <- []byte(dataUpdateMessageString)
			g.player2.send <- []byte(dataUpdateMessageString)

		}
	}
}
