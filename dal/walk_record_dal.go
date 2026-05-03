package dal

import (
	"demo/model"
	"demo/util/postgres"

	"github.com/sirupsen/logrus"
)

func CreateWalkRecord(walkRecord *model.WalkRecord) error {
	db := postgres.GetDB()
	return db.Create(walkRecord).Error
}

func GetWalkRecordById(id int64) (*model.WalkRecord, error) {
	db := postgres.GetDB()
	var walkRecord *model.WalkRecord
	err := db.First(&walkRecord, id).Error
	if err != nil {
		logrus.Errorf("get walk record by id failed, err = %v", err)
		return nil, err
	}
	return walkRecord, nil
}

func UpdateWalkRecord(walkRecord *model.WalkRecord) error {
	db := postgres.GetDB()
	return db.Save(walkRecord).Error
}

func GetWalkRecordsBySubOrderID(subOrderID int64) ([]*model.WalkRecord, error) {
	db := postgres.GetDB()
	var walkRecordList []*model.WalkRecord
	err := db.Where("sub_order_id = ?", subOrderID).Order("created_at desc").Find(&walkRecordList).Error
	if err != nil {
		logrus.Errorf("get walk records by sub order id failed, err = %v", err)
		return nil, err
	}
	return walkRecordList, nil
}
