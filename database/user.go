package database

import (
	"aurora-agent/database/model"
)

func init() {
	db, _ = DBConnect()
}

func GetAllUsers() ([]model.User, error) {
	var users []model.User
	result := db.Model(&model.User{}).Find(&users).Error
	return users, result
}
