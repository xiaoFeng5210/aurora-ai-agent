package database

import (
	"aurora-agent/database/model"
	"aurora-agent/handler/dto"
	"errors"
	"time"
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

// 根据ID获取用户信息
func GetUserById(id int) (model.User, error) {
	var user model.User
	result := db.Model(&model.User{}).Where("id = ?", id).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	if result.RowsAffected == 0 {
		return user, errors.New("user not found")
	}
	return user, nil
}

// 根据Username获取用户信息
func GetUserByUsername(username string) (model.User, error) {
	var user model.User
	result := db.Model(&model.User{}).Where("username = ?", username).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	if result.RowsAffected == 0 {
		return user, errors.New("user not found")
	}
	return user, nil
}

// 根据用户名模糊查询
func GetUsersByUsername(username string) ([]model.User, error) {
	var users []model.User
	result := db.Model(&model.User{}).Where("username LIKE ?", "%"+username+"%").Find(&users)
	if result.Error != nil {
		return users, result.Error
	}
	return users, nil
}

// 筛选年龄大于多少岁的用户
func GetUsersByBirthdayMoreThanAge(birthday time.Time, queryAge int) ([]model.User, error) {
	// 先把年龄计算出来根据生日
	trueAge := time.Now().Year() - birthday.Year()

	if trueAge > queryAge {
		return nil, errors.New("age is too young")
	}

	queryBirthday := time.Now().AddDate(-queryAge, 0, 0)

	var users []model.User
	result := db.Model(&model.User{}).Where("birthday >= ?", queryBirthday).Find(&users)
	if result.Error != nil {
		return users, result.Error
	}
	return users, nil
}

func QueryUser(filter dto.QueryUserDTO) ([]model.User, error) {
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

	if !filter.Birthday.IsZero() {
		queryDB = queryDB.Where("birthday = ?", filter.Birthday)
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		queryDB = queryDB.Offset((filter.Page - 1) * filter.PageSize).Limit(filter.PageSize)
	} else {
		queryDB = queryDB.Limit(10)
	}

	var users []model.User
	result := queryDB.Find(&users)
	if result.Error != nil {
		return users, result.Error
	}
	return users, nil
}
