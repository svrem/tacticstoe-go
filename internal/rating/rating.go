package rating

import (
	"math"
	db "tacticstoe/internal/database"

	"gorm.io/gorm"
)

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

	println(expectedScore1, expectedScore2)

	newRating1 := rating1 + int(K*(float64(outcome)-expectedScore1))
	newRating2 := rating2 + int(K*(float64(1-outcome)-expectedScore2))

	return newRating1, newRating2
}

// if outcome == 1, player1 wins, if outcome == 0.5, draw, if outcome == 0, player2 wins
func UpdateRatings(database *gorm.DB, player1 *db.User, player2 *db.User, outcome float32) (int, int) {
	newRating1, newRating2 := CalcNewRating(player1.ELO_Rating, player2.ELO_Rating, outcome)

	player1.ELO_Rating = newRating1
	player2.ELO_Rating = newRating2

	db.UpdateUser(database, player1)
	db.UpdateUser(database, player2)

	return newRating1, newRating2
}
