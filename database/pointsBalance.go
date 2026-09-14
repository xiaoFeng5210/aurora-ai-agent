package database

import "aurora-agent/database/model"

func CreatePointsBalance(userId int) error {
	data := &model.PointsBalance{
		UserId: userId,
	}

	return db.Model(&model.PointsBalance{}).Create(data).Error
}
