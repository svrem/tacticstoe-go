package db

type User struct {
	ID string `gorm:"primarykey"`

	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`

	ELO_Rating int `json:"elo_rating" gorm:"default:1000"`

	Provider   string
	ProviderId string
}
