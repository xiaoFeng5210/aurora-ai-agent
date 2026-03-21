package database

import (
	"aurora-agent/database/model"
	"time"

	"gorm.io/gorm"
)

type UserQueryFilter struct {
	Email    string
	Username string
	Phone    string
	Birthday *time.Time
	Page     int
	PageSize int
}

func init() {
	db, _ = DBConnect()
}

func CreateUser(user model.User) error {
	return db.Model(&model.User{}).Create(&user).Error
}

func GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := db.Model(&model.User{}).Order("id ASC").Find(&users).Error
	return users, err
}

func GetUserByEmail(email string) (model.User, error) {
	var user model.User
	err := db.Model(&model.User{}).Where("email = ?", email).First(&user).Error
	return user, err
}

func GetUserById(id int) (model.User, error) {
	var user model.User
	err := db.Model(&model.User{}).Where("id = ?", id).First(&user).Error
	return user, err
}

func QueryUsers(filter UserQueryFilter) ([]model.User, error) {
	queryDB := db.Model(&model.User{})
	if filter.Email != "" {
		queryDB = queryDB.Where("email = ?", filter.Email)
	}

	if filter.Username != "" {
		queryDB = queryDB.Where("username = ?", filter.Username)
	}

	if filter.Phone != "" {
		queryDB = queryDB.Where("phone = ?", filter.Phone)
	}

	if filter.Birthday != nil {
		queryDB = queryDB.Where("birthday = ?", *filter.Birthday)
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		queryDB = queryDB.Offset((filter.Page - 1) * filter.PageSize).Limit(filter.PageSize)
	} else if filter.PageSize > 0 {
		queryDB = queryDB.Limit(filter.PageSize)
	}

	var users []model.User
	err := queryDB.Order("id ASC").Find(&users).Error
	return users, err
}

func UsernameExists(username string, excludeID int) (bool, error) {
	var count int64
	queryDB := db.Model(&model.User{}).Where("username = ?", username)
	if excludeID > 0 {
		queryDB = queryDB.Where("id <> ?", excludeID)
	}

	if err := queryDB.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func EmailExists(email string, excludeID int) (bool, error) {
	var count int64
	queryDB := db.Model(&model.User{}).Where("email = ?", email)
	if excludeID > 0 {
		queryDB = queryDB.Where("id <> ?", excludeID)
	}

	if err := queryDB.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdateUserByID(id int, updates map[string]any) error {
	result := db.Model(&model.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func SoftDeleteUser(id int) error {
	result := db.Where("id = ?", id).Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
