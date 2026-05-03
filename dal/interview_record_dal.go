package dal

import (
	"demo/model"
	"demo/util/postgres"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateInterviewRecord 创建面试预约记录
func CreateInterviewRecord(interviewRecord model.InterviewRecord) (*model.InterviewRecord, error) {
	db := postgres.GetDB()
	err := db.Create(&interviewRecord).Error
	if err != nil {
		logrus.Errorf("create interview record failed, err = %v", err)
		return nil, err
	}
	return &interviewRecord, nil
}

// GetInterviewRecordByID 根据ID查询面试预约记录
func GetInterviewRecordByID(id int64) (*model.InterviewRecord, error) {
	db := postgres.GetDB()
	var interviewRecord model.InterviewRecord
	err := db.Where("id = ?", id).First(&interviewRecord).Error
	if err != nil {
		logrus.Errorf("get interview record by ID failed, err = %v", err)
		return nil, err
	}
	return &interviewRecord, nil
}

// GetInterviewRecordByOrderID 根据订单ID查询面试预约记录
func GetInterviewRecordByOrderID(orderID int64) (*model.InterviewRecord, error) {
	db := postgres.GetDB()
	var interviewRecord model.InterviewRecord
	err := db.Where("order_id = ?", orderID).Order("created_at desc").First(&interviewRecord).Error
	if err != nil {
		logrus.Errorf("get interview record by orderID failed, err = %v", err)
		return nil, err
	}
	return &interviewRecord, nil
}

// UpdateInterviewRecord 更新面试预约记录
func UpdateInterviewRecord(interviewRecord *model.InterviewRecord) error {
	db := postgres.GetDB()
	result := db.Save(interviewRecord)
	if result.Error != nil {
		logrus.Errorf("update interview record failed, err = %v", result.Error)
		return result.Error
	}
	return nil
}

// DeleteInterviewRecord 删除面试预约记录
func DeleteInterviewRecord(id uint64) error {
	db := postgres.GetDB()
	result := db.Where("id = ?", id).Delete(&model.InterviewRecord{})
	if result.Error != nil {
		logrus.Errorf("delete interview record failed, err = %v", result.Error)
		return result.Error
	}
	return nil
}

// QueryInterviewRecords 条件查询面试预约记录
func QueryInterviewRecords(orderID string, status model.InterviewStatus) ([]*model.InterviewRecord, error) {
	db := postgres.GetDB()
	var interviewRecords []*model.InterviewRecord
	query := db.Model(&model.InterviewRecord{})

	// 条件过滤
	if orderID != "" {
		query = query.Where("order_id = ?", orderID)
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	err := query.Find(&interviewRecords).Error
	if err != nil {
		logrus.Errorf("query interview records failed, err = %v", err)
		return nil, err
	}
	return interviewRecords, nil
}

// GetInterviewRecordsWithPagination 分页查询面试预约记录
func GetInterviewRecordsWithPagination(orderID string, status model.InterviewStatus, page, size int) ([]*model.InterviewRecord, int64, error) {
	db := postgres.GetDB()
	var interviewRecords []*model.InterviewRecord
	var total int64

	// 计算偏移量
	offset := (page - 1) * size

	// 构建查询
	query := db.Model(&model.InterviewRecord{})
	if orderID != "" {
		query = query.Where("order_id = ?", orderID)
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	// 查询总数
	err := query.Count(&total).Error
	if err != nil {
		logrus.Errorf("count interview records failed, err = %v", err)
		return nil, 0, err
	}

	// 查询分页数据
	err = query.Offset(offset).Limit(size).Find(&interviewRecords).Error
	if err != nil {
		logrus.Errorf("get interview records with pagination failed, err = %v", err)
		return nil, 0, err
	}

	return interviewRecords, total, nil
}

// UpdateInterviewStatus 更新面试预约状态
func UpdateInterviewStatus(id uint64, status model.InterviewStatus) error {
	db := postgres.GetDB()
	result := db.Model(&model.InterviewRecord{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		logrus.Errorf("update interview status failed, err = %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		logrus.Warnf("update interview status: record not found, id = %d", id)
		return gorm.ErrRecordNotFound
	}
	return nil
}
