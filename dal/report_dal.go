package dal

import (
	"context"
	"encoding/json"
	"errors"

	"demo/model"
	"demo/util/postgres"

	"gorm.io/gorm"
)

// CreateReport 创建举报记录
func CreateReport(ctx context.Context, report *model.Report) error {
	db := postgres.GetDB()
	return db.WithContext(ctx).Create(report).Error
}

// GetReportByID 根据ID获取举报记录
func GetReportByID(ctx context.Context, id int64) (*model.Report, error) {
	db := postgres.GetDB()
	var report model.Report
	err := db.WithContext(ctx).Where("id = ?", id).First(&report).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &report, err
}

// GetReportsByReporter 获取举报人的举报记录
func GetReportsByReporter(ctx context.Context, reporterID string, page, pageSize int) ([]*model.Report, int64, error) {
	db := postgres.GetDB()
	var reports []*model.Report
	var total int64

	query := db.WithContext(ctx).Where("reporter_id = ?", reporterID)

	// 获取总数
	if err := query.Model(&model.Report{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&reports).Error

	return reports, total, err
}

// GetReportsByReported 获取被举报人的举报记录
func GetReportsByReported(ctx context.Context, reportedID string, page, pageSize int) ([]*model.Report, int64, error) {
	db := postgres.GetDB()
	var reports []*model.Report
	var total int64

	query := db.WithContext(ctx).Where("reported_id = ?", reportedID)

	if err := query.Model(&model.Report{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&reports).Error

	return reports, total, err
}

// UpdateReportStatus 更新举报状态
func UpdateReportStatus(ctx context.Context, id int64, status int8, handlerID string) error {
	db := postgres.GetDB()
	return db.WithContext(ctx).Model(&model.Report{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"handler_id": handlerID,
		}).Error
}

// UpdateReportResult 更新处理结果
func UpdateReportResult(ctx context.Context, id int64, result string) error {
	db := postgres.GetDB()
	return db.WithContext(ctx).Model(&model.Report{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        model.ReportStatusResolved,
			"handle_result": result,
		}).Error
}

// GetPendingReports 获取待处理的举报
func GetPendingReports(ctx context.Context, page, pageSize int) ([]*model.Report, int64, error) {
	db := postgres.GetDB()
	var reports []*model.Report
	var total int64

	query := db.WithContext(ctx).Where("status = ?", model.ReportStatusPending)

	if err := query.Model(&model.Report{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&reports).Error

	return reports, total, err
}

// ConvertEvidenceToJSON 将证据字符串数组转换为JSON格式
func ConvertEvidenceToJSON(evidence []string) (model.JSON, error) {
	if evidence == nil {
		return nil, nil
	}
	return json.Marshal(evidence)
}

// ParseEvidenceFromJSON 从JSON解析证据
func ParseEvidenceFromJSON(evidence model.JSON) ([]string, error) {
	if evidence == nil {
		return nil, nil
	}
	var result []string
	err := json.Unmarshal(evidence, &result)
	return result, err
}
