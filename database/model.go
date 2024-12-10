package db

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`

	Provider   string
	ProviderId string
}
