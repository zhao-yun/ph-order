package model

import "demo/util/common"

type SitterLocation struct {
	ID         int64            `gorm:"primaryKey;autoIncrement;comment('位置记录ID')" json:"id"`
	SitterID   string           `gorm:"type:varchar(50);not null;comment('Sitter ID')" json:"sitterId"`
	OrderID    int64            `gorm:"not null;comment('主订单ID')" json:"orderId"`
	SubOrderID int64            `gorm:"not null;comment('子订单ID')" json:"subOrderId"`
	Lat        float64          `gorm:"type:decimal(10,6);not null;comment('纬度')" json:"lat"`
	Lng        float64          `gorm:"type:decimal(10,6);not null;comment('经度')" json:"lng"`
	Timestamp  int64            `gorm:"not null;comment('时间戳')" json:"timestamp"`
	CreatedAt  common.Timestamp `gorm:"type:timestamp;autoCreateTime;comment('创建时间')" json:"createdAt"`
}

type SitterLocationReq struct {
	SitterID   string  `json:"SitterID" form:"SitterID" binding:"required"`
	OrderID    int64   `json:"OrderID" form:"OrderID" binding:"required"`
	SubOrderID int64   `json:"SubOrderID" form:"SubOrderID" binding:"required"`
	Lat        float64 `json:"Lat" form:"Lat" binding:"required"`
	Lng        float64 `json:"Lng" form:"Lng" binding:"required"`
	Timestamp  int64   `json:"Timestamp" form:"Timestamp" binding:"required"`
}

type SitterLocationResponse struct {
	ID         int64   `json:"id"`
	SitterID   string  `json:"sitterId"`
	OrderID    int64   `json:"orderId"`
	SubOrderID int64   `json:"subOrderId"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	Timestamp  int64   `json:"timestamp"`
	CreatedAt  string  `json:"createdAt"`
}

type SitterLocationQueryParams struct {
	SitterID   string `json:"SitterID" form:"sitter_id"`
	OrderID    int64  `json:"OrderID" form:"order_id"`
	SubOrderID int64  `json:"SubOrderID" form:"sub_order_id"`
}

func (SitterLocation) TableName() string {
	return "sitter_locations"
}
