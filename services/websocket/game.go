package websocket_service

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type GameStart struct {
	StartingPlayer string `json:"starting_player"`
}

type GameAbort struct {
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
	Winner string     `json:"winner"`
	Coords *[3][2]int `json:"coords"`
}

type Game struct {
	id string

	player1 *Client
	player2 *Client

	board [4][4]int

	activePlayer *Client

	makeAction chan GameAction
	unregister chan *Client
}

func newGame(player1 *Client, player2 *Client) Game {
	gameJoinDataMessage := DataMessage[GameStart]{
		Type: "game_start",
		Data: GameStart{
			StartingPlayer: player1.id,
		},
	}
	data, err := gameJoinDataMessage.Marshal()
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

		activePlayer: player1,

		makeAction: make(chan GameAction),
		unregister: make(chan *Client),
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

			gameAbortedMessage := DataMessage[GameAbort]{
				Type: "game_abort",
				Data: GameAbort{},
			}
			gameAbortedMessageString, err := gameAbortedMessage.Marshal()

			if err != nil {
				slog.Error("Error marshalling game abort data.")
				continue
			}

			g.player1.send <- []byte(gameAbortedMessageString)
			g.player2.send <- []byte(gameAbortedMessageString)

			gp.unregister <- client

		case action := <-g.makeAction:
			if action.player != g.activePlayer {
				continue
			}

			isActivePlayerPlayer1 := action.player == g.player1

			var newState int
			if isActivePlayerPlayer1 {
				newState = 1
			} else {
				newState = 2
			}

			nextBoard := g.board
			nextBoard[action.x][action.y] = newState
			winnerData := checkBoardForWin(nextBoard)

			isValidMove := checkForMoveValidity(winnerData, action, g.board, nextBoard)

			if !isValidMove {
				continue
			}

			g.board = nextBoard

			if checkDraw(g.board) {
				gameEndMessage := DataMessage[GameEnd]{
					Type: "game_end",
					Data: GameEnd{
						Winner: "draw",
						Coords: nil,
					},
				}

				gameEndMessageString, err := gameEndMessage.Marshal()

				if err != nil {

					slog.Error("Error marshalling game end data.")
					continue
				}

				g.player1.send <- []byte(gameEndMessageString)
				g.player2.send <- []byte(gameEndMessageString)
			}
			// Switch active player
			if isActivePlayerPlayer1 {
				g.activePlayer = g.player2
			} else {
				g.activePlayer = g.player1
			}

			if err := sendGameUpdate(action, newState, g); err != nil {
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
						Coords: &winnerData.coords,
					},
				}

				gameEndMessageString, err := gameEndMessage.Marshal()

				if err != nil {
					slog.Error("Error marshalling game end data.")
					continue
				}

				g.player1.send <- []byte(gameEndMessageString)
				g.player2.send <- []byte(gameEndMessageString)

				time.Sleep(500 * time.Millisecond)

				g.player1.hub.unregister <- g.player1
			}

		}
	}
}

func checkDraw(board [4][4]int) bool {
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			if board[x][y] == 0 {
				return false
			}
		}
	}

	return true
}

func checkForMoveValidity(winnerData *BoardWinData, action GameAction, board [4][4]int, nextBoard [4][4]int) bool {

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
	if (action.x == 0 || nextBoard[action.x-1][action.y] != opponent) &&
		(action.x == 3 || nextBoard[action.x+1][action.y] != opponent) &&
		(action.y == 0 || nextBoard[action.x][action.y-1] != opponent) &&
		(action.y == 3 || nextBoard[action.x][action.y+1] != opponent) {
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
			ActivePlayer: g.activePlayer.id,
		}}

	gameUpdateMessageString, err := gameUpdateMessage.Marshal()

	if err != nil {
		slog.Error("Error marshalling game update data.")
		return err
	}

	g.player1.send <- []byte(gameUpdateMessageString)
	g.player2.send <- []byte(gameUpdateMessageString)

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
