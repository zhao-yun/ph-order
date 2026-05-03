package model

import (
	"demo/util/common"
	"encoding/json"
)

type WalkRecord struct {
	ID         int64            `gorm:"primaryKey;autoIncrement;comment('轨迹ID')" json:"id"`
	OrderID    int64            `gorm:"not null;comment('主订单ID')" json:"orderId"`
	SubOrderID int64            `gorm:"not null;comment('子订单ID')" json:"subOrderId"`
	Path       json.RawMessage  `gorm:"type:json;comment('路径数据 LatLng数组')" json:"-"`
	PathData   []LatLng         `gorm:"-" json:"path,omitempty"`
	CreatedAt  common.Timestamp `gorm:"type:timestamp;autoCreateTime;comment('创建时间')" json:"createdAt,omitempty"`
	UpdatedAt  common.Timestamp `gorm:"type:timestamp;autoUpdateTime;comment('更新时间')" json:"-"`
}

type LatLng struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type WalkDetail struct {
	ID                string   `json:"id"`
	OrderID           int64    `json:"orderId"`
	SubOrderID        int64    `json:"subOrderId"`
	Path              []LatLng `json:"path"`
	DurationSeconds   int64    `json:"durationSeconds"`
	DistanceMeters    float64  `json:"distanceMeters"`
	WalkThumbnailUrl  string   `json:"walkThumbnailUrl"`
	CreatedAt         string   `json:"createdAt,omitempty"`
}

type CreateWalkRecordReq struct {
	SubOrderID int64    `json:"SubOrderID" form:"SubOrderID" binding:"required"`
	Path       []LatLng `json:"Path" form:"Path" binding:"required"`
}

func (WalkRecord) TableName() string {
	return "walk_records"
}
