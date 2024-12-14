package db

type User struct {
	ID string `gorm:"primarykey"`

	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`

	Provider   string
	ProviderId string
}
