package db

import (
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, username string, profilePicture string, provider string, providerId string) (*User, error) {
	userId := uuid.New().String()

	existingUserWithProvider := User{}

	if err := db.Where("provider = ? AND provider_id = ?", provider, providerId).First(&existingUserWithProvider).Error; err == nil {
		slog.Info("User already exists, returning existing user")
		return &existingUserWithProvider, nil
	}

	user := User{
		ID: userId,

		Username:       username,
		ProfilePicture: profilePicture,

		Provider:   provider,
		ProviderId: providerId,
	}

	err := db.Create(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(db *gorm.DB, userId string) *User {
	user := User{
		ID: userId,
	}
	if err := db.First(&user).Error; err != nil {
		return nil
	}

	return &user
}

func UpdateUser(db *gorm.DB, user *User) error {
	res := db.Save(user)

	if res.Error != nil {
		return res.Error
	}

	return nil
}
