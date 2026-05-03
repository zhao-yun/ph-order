package service

import (
	"context"
	"errors"

	"demo/dal"
	"demo/model"
)

// SubmitReport 提交举报
func SubmitReport(ctx context.Context, req *model.ReportRequest, reporterID string) (*model.ReportResponse, error) {
	// 检查是否举报自己

	if reporterID == req.ReportedID {
		return nil, errors.New("不能举报自己")
	}

	report := &model.Report{
		ReporterID:   reporterID,
		ReportedID:   req.ReportedID,
		ReportedName: req.ReportedName,
		ReportReason: req.ReportReason,
		ReportDesc:   req.ReportDesc,
		ReportImg:    req.ReportImg,
		Status:       model.ReportStatusProcessing,
	}

	if err := dal.CreateReport(ctx, report); err != nil {
		return nil, err
	}

	return &model.ReportResponse{
		ID:           report.ID,
		ReportedID:   report.ReportedID,
		ReportedName: report.ReportedName,
		ReportReason: report.ReportReason,
		ReportDesc:   report.ReportDesc,
		ReportImg:    report.ReportImg,
		Status:       report.Status,
		CreatedAt:    report.CreatedAt,
	}, nil
}

// GetReportDetail 获取举报详情
func GetReportDetail(ctx context.Context, id int64, userID string) (*model.ReportResponse, error) {
	report, err := dal.GetReportByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if report == nil || (report.ReporterID != userID && report.ReportedID != userID) {
		return nil, errors.New("举报记录不存在或无权查看")
	}

	return &model.ReportResponse{
		ID:           report.ID,
		ReportedID:   report.ReportedID,
		ReportedName: report.ReportedName,
		ReportReason: report.ReportReason,
		ReportDesc:   report.ReportDesc,
		ReportImg:    report.ReportImg,
		Status:       report.Status,
		CreatedAt:    report.CreatedAt,
	}, nil
}

// GetMyReports 获取我的举报记录
func GetMyReports(ctx context.Context, userID string, page, pageSize int) ([]*model.ReportResponse, int64, error) {
	reports, total, err := dal.GetReportsByReporter(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var result []*model.ReportResponse
	for _, report := range reports {
		result = append(result, &model.ReportResponse{
			ID:           report.ID,
			ReportedID:   report.ReportedID,
			ReportedName: report.ReportedName,
			ReportReason: report.ReportReason,
			ReportDesc:   report.ReportDesc,
			ReportImg:    report.ReportImg,
			Status:       report.Status,
			CreatedAt:    report.CreatedAt,
		})
	}

	return result, total, nil
}

// GetReportsAgainstMe 获取针对我的举报
func GetReportsAgainstMe(ctx context.Context, userID string, page, pageSize int) ([]*model.ReportResponse, int64, error) {
	reports, total, err := dal.GetReportsByReported(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var result []*model.ReportResponse
	for _, report := range reports {
		result = append(result, &model.ReportResponse{
			ID:           report.ID,
			ReportedID:   report.ReportedID,
			ReportedName: report.ReportedName,
			ReportReason: report.ReportReason,
			ReportDesc:   report.ReportDesc,
			ReportImg:    report.ReportImg,
			Status:       report.Status,
			CreatedAt:    report.CreatedAt,
		})
	}

	return result, total, nil
}

// ProcessReport 处理举报
func ProcessReport(ctx context.Context, id int64, handlerID string) error {
	return dal.UpdateReportStatus(ctx, id, model.ReportStatusProcessing, handlerID)
}

// ResolveReport 解决举报
func ResolveReport(ctx context.Context, id int64, result string) error {
	return dal.UpdateReportResult(ctx, id, result)
}

// GetPendingReports 获取待处理举报列表(管理员用)
func GetPendingReports(ctx context.Context, page, pageSize int) ([]*model.ReportResponse, int64, error) {
	reports, total, err := dal.GetPendingReports(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var result []*model.ReportResponse
	for _, report := range reports {
		result = append(result, &model.ReportResponse{
			ID:           report.ID,
			ReportedID:   report.ReportedID,
			ReportedName: report.ReportedName,
			ReportReason: report.ReportReason,
			ReportDesc:   report.ReportDesc,
			ReportImg:    report.ReportImg,
			Status:       report.Status,
			CreatedAt:    report.CreatedAt,
		})
	}

	return result, total, nil
}
