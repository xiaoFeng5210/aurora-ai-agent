package database

import (
	"aurora-agent/database/model"
	"errors"
)

func init() {
	db, _ = DBConnect()
}

func CreateUser(user model.User) error {
	result := db.Model(&model.User{}).Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAllUsers() ([]model.User, error) {
	var users []model.User
	result := db.Model(&model.User{}).Find(&users).Error
	return users, result
}

func GetUserByEmail(email string) (model.User, error) {
	var user model.User
	result := db.Model(&model.User{}).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	if result.RowsAffected == 0 {
		return user, errors.New("user not found")
	}
	return user, nil
}
