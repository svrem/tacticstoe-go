package rating

import "math"

const (
	K = 30
)

func probability(rating1 int, rating2 int) float64 {
	return 1.0 / (1.0 + math.Pow(10, ((float64(rating2)-float64(rating1))/400)))
}

// if outcome == 1, player1 wins, if outcome == 0.5, draw, if outcome == 0, player2 wins
func CalcNewRating(rating1 int, rating2 int, outcome float32) (int, int) {
	expectedScore1 := probability(rating1, rating2)
	expectedScore2 := probability(rating2, rating1)

	newRating1 := rating1 + int(K*(float64(outcome)-expectedScore1))
	newRating2 := rating2 + int(K*(float64(1-outcome)-expectedScore2))

	return newRating1, newRating2
}
