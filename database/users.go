package db

import "gorm.io/gorm"

func CreateUser(db *gorm.DB, user *User) {
	db.Create(user)
}

func GetUserByID(db *gorm.DB, id uint) *User {
	var user User
	db.First(&user, id)
	return &user
}
