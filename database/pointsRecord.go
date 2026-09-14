package database

import "aurora-agent/database/model"

func CreatePointsRecord(record model.PointsRecord) error {
	return db.Model(&model.PointsRecord{}).Create(&record).Error
}

func GetPointsRecordByUserId(userId int) ([]model.PointsRecord, error) {
	var records []model.PointsRecord
	err := db.Model(&model.PointsRecord{}).Where("user_id = ?", userId).
		Order("create_time DESC").Find(&records).Error
	return records, err
}
