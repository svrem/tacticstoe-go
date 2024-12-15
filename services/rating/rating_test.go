package rating_test

import (
	"tacticstoe/services/rating"
	"testing"
)

func TestCalcRatingChange(t *testing.T) {
	winnerRating := 1200
	loserRating := 1000

	newWinnerRating, newLoserRating := rating.CalcNewRating(winnerRating, loserRating, 1)

	println(newWinnerRating, newLoserRating)

	if newWinnerRating != 1207 || newLoserRating != 993 {
		t.Errorf("Expected newWinnerRating to be 1207 and newLoserRating to be 993, got %d and %d", newWinnerRating, newLoserRating)
	}
}
