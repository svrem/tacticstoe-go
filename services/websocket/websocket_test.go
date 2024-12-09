package websocket_service

import "testing"

func TestCheckForWin(t *testing.T) {
	boards := [][4][4]int{
		{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
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
		0,
		1,
		1,
		1,
		2,
		2,
		2,
		0,
	}

	for index, board := range boards {
		if checkBoardForWin(board) != expecting[index] {
			t.Errorf("Expected false, got true.")
		}
	}
}
