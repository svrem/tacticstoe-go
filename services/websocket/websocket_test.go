package websocket_service

import "testing"

func TestCheckForWin(t *testing.T) {
	boards := [][4][4]int{
		{
			{1, 0, 0, 0},
			{1, 0, 0, 2},
			{1, 0, 2, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0},
			{0, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{1, 1, 1, 0},
		},
		{
			{2, 0, 0, 0},
			{2, 0, 0, 0},
			{2, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{2, 0, 0, 0},
			{0, 2, 0, 0},
			{0, 0, 2, 0},
			{0, 0, 0, 0},
		}, {
			{0, 0, 2, 0},
			{0, 2, 0, 0},
			{2, 0, 0, 0},
			{0, 0, 0, 0},
		}, {
			{0, 0, 0, 0},
			{0, 2, 0, 0},
			{0, 0, 2, 0},
			{0, 2, 0, 0},
		},
	}
	var expecting = []int{
		1,
		1,
		1,
		1,
		2,
		2,
		2,
		0,
	}

	for index, board := range boards {
		winnerData := checkBoardForWin(board)
		if (winnerData == nil && expecting[index] != 0) || (winnerData != nil && winnerData.player != expecting[index]) {
			t.Errorf("Expected false, got true.")
		}
	}
}
