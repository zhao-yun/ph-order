package model

import (
	"time"
)

type Report struct {
	ID           int64     `gorm:"primaryKey;column:id" json:"ID"`
	ReporterID   string    `gorm:"column:reporter_id;type:varchar(50);not null" json:"ReporterID"`
	ReportedID   string    `gorm:"column:reported_id;type:varchar(50);not null" json:"ReportedID"`
	ReportedName string    `gorm:"column:reported_name;type:varchar(50)" json:"ReportedName"`
	ReportReason string    `gorm:"column:report_reason;type:varchar(50)" json:"ReportReason"`
	ReportDesc   string    `gorm:"column:report_desc;type:text;not null" json:"ReportDesc"`
	ReportImg    string    `gorm:"column:report_img;type:text" json:"ReportImg"`
	Status       int8      `gorm:"column:status;type:tinyint(1);default:0" json:"Status"`
	HandlerID    string    `gorm:"column:handler_id;type:varchar(50)" json:"HandlerID"`
	HandleResult string    `gorm:"column:handle_result;type:text" json:"HandleResult"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null;autoCreateTime" json:"CreatedAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null;autoUpdateTime" json:"UpdatedAt"`
}

func (Report) TableName() string {
	return "report"
}

// JSON 自定义JSON类型
type JSON []byte

// 举报类型常量
const (
	ReportTypeFraud         = "fraud"         // 欺诈
	ReportTypeHarassment    = "harassment"    // 骚扰
	ReportTypeInappropriate = "inappropriate" // 不当内容
	ReportTypeOther         = "other"         // 其他
)

// 处理状态常量
const (
	ReportStatusPending    = 0 // 待处理
	ReportStatusProcessing = 1 // 处理中
	ReportStatusResolved   = 2 // 已处理
)

// ReportRequest 举报请求结构体
type ReportRequest struct {
	ReportedID   string `json:"ReportedID" binding:"required"`   // 被举报人ID
	ReportedName string `json:"ReportedName"`                    // 被举报人姓名
	ReportReason string `json:"ReportReason" binding:"required"` // 举报类型
	ReportDesc   string `json:"ReportDesc" binding:"required"`   // 举报描述
	ReportImg    string `json:"ReportImg"`                       // 举报图片
}

// ReportResponse 举报响应结构体
type ReportResponse struct {
	ID           int64     `json:"ID"`
	ReportedID   string    `json:"ReportedID"`
	ReportedName string    `json:"ReportedName"`
	ReportReason string    `json:"ReportReason"`
	ReportDesc   string    `json:"ReportDesc"`
	ReportImg    string    `json:"ReportImg"`
	Status       int8      `json:"Status"`
	CreatedAt    time.Time `json:"CreatedAt"`
}
