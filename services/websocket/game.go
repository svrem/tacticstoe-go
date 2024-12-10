package websocket_service

import (
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
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

type GameEnd struct {
	Winner string    `json:"winner"`
	Coords [3][2]int `json:"coords"`
}

type Game struct {
	id string

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
		id: uuid.New().String(),

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
			defer func() {
				recover()
			}()

			if client == nil {
			}

			g.player1.send <- []byte("Hello")
			g.player2.send <- []byte("Hello")

			gp.close_game <- g

		case action := <-g.make_action:
			if action.player != g.active_player {
				continue
			}

			is_active_player_player1 := action.player == g.player1

			var new_state int
			if is_active_player_player1 {
				new_state = 1
			} else {
				new_state = 2
			}

			next_board := g.board
			next_board[action.x][action.y] = new_state
			winnerData := checkBoardForWin(next_board)

			isValidMove := checkForMoveValidity(winnerData, action, g.board, next_board)

			if !isValidMove {
				continue
			}

			g.board = next_board

			// Switch active player
			if is_active_player_player1 {
				g.active_player = g.player2
			} else {
				g.active_player = g.player1
			}

			if err := sendGameUpdate(action, new_state, g); err != nil {
				slog.Error("Error sending game update.")
				continue
			}

			if winnerData != nil {
				var winner string
				if winnerData.player == 1 {
					winner = g.player1.id
				} else {
					winner = g.player2.id
				}

				gameEndMessage := DataMessage[GameEnd]{
					Type: "game_end",
					Data: GameEnd{
						Winner: winner,
						Coords: winnerData.coords,
					},
				}

				dataEndMessageString, err := json.Marshal(gameEndMessage)

				if err != nil {
					slog.Error("Error marshalling game end data.")
					continue
				}

				g.player1.send <- []byte(dataEndMessageString)
				g.player2.send <- []byte(dataEndMessageString)

				gp.close_game <- g
			}

		}
	}
}

func checkForMoveValidity(winnerData *BoardWinData, action GameAction, board [4][4]int, next_board [4][4]int) bool {

	if board[action.x][action.y] != 0 {
		return false
	}

	if winnerData == nil {
		return true
	}

	var opponent int
	if winnerData.player == 1 {
		opponent = 2
	} else {
		opponent = 1
	}

	// Check if the move is valid,
	// by checking if the opponent has placed their tick in a square directly adjacent
	// to the square the player is trying to place their tick in
	// e.g. if the player is trying to place their tick in (1, 1),
	// and the opponent has placed their tick in (0, 1), (2, 1), (1, 0), or (1, 2),
	// then the move is valid
	if (action.x == 0 || next_board[action.x-1][action.y] != opponent) &&
		(action.x == 3 || next_board[action.x+1][action.y] != opponent) &&
		(action.y == 0 || next_board[action.x][action.y-1] != opponent) &&
		(action.y == 3 || next_board[action.x][action.y+1] != opponent) {
		return false
	}

	return true
}

func sendGameUpdate(action GameAction, new_state int, g *Game) error {
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
		return err
	}

	g.player1.send <- []byte(dataUpdateMessageString)
	g.player2.send <- []byte(dataUpdateMessageString)

	return nil
}

type BoardWinData struct {
	player int
	coords [3][2]int
}

func checkBoardForWin(board [4][4]int) *BoardWinData {
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			if board[x][y] == 0 {
				continue
			}

			if x > 0 && x < 3 && board[x][y] == board[x-1][y] && board[x][y] == board[x+1][y] {
				return &BoardWinData{player: board[x][y], coords: [3][2]int{{x - 1, y}, {x, y}, {x + 1, y}}}
			}

			if y > 0 && y < 3 && board[x][y] == board[x][y-1] && board[x][y] == board[x][y+1] {
				return &BoardWinData{player: board[x][y], coords: [3][2]int{{x, y - 1}, {x, y}, {x, y + 1}}}
			}

			if x > 0 && x < 3 && y > 0 && y < 3 && board[x][y] == board[x-1][y-1] && board[x][y] == board[x+1][y+1] {
				return &BoardWinData{player: board[x][y], coords: [3][2]int{{x - 1, y - 1}, {x, y}, {x + 1, y + 1}}}
			}

			if x > 0 && x < 3 && y > 0 && y < 3 && board[x][y] == board[x-1][y+1] && board[x][y] == board[x+1][y-1] {
				return &BoardWinData{player: board[x][y], coords: [3][2]int{{x - 1, y + 1}, {x, y}, {x + 1, y - 1}}}
			}
		}
	}

	return nil
}
