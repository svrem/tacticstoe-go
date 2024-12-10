package db

import "gorm.io/gorm"

func CreateUser(db *gorm.DB, user *User) {
	db.Create(user)
}

func GetUserByID(db *gorm.DB, userId uint) *User {
	var user User
	if err := db.First(&user, userId).Error; err != nil {
		return nil
	}

	return &user
}
