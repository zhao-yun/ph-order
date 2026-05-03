package model

import (
	"time"

	"demo/util/common"
)

type OrderModificationLog struct {
	ID              int64  `gorm:"primaryKey;autoIncrement"`
	OrderID         int64  `gorm:"not null"`
	OwnerID         string `gorm:"type:varchar(50);not null"`
	SitterID        string `gorm:"type:varchar(50);not null"`
	PreviousDate    common.Date
	NewDate         common.Date
	PreviousPetList string // 假设宠物信息以 JSON 格式存储
	NewPetList      string
	PreviousPrice   float64
	NewPrice        float64
	State           OrderModificationStatus `gorm:"not null;default:0"`
	Type            OrderModificationType   `gorm:"not null;"`
	CreatedAt       time.Time               `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time               `gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName 设置表名
func (OrderModificationLog) TableName() string {
	return "order_modification_log"
}

// UserUpdateOrderReq 用户更新订单请求
type UserUpdateOrderReq struct {
	OrderID int64       `json:"OrderID" form:"orderId"` // 订单ID
	ToDate  common.Date `json:"ToDate" form:"toDate"`
	PetList []*OrderPet `json:"PetList" form:"petList"`
}

type SitterUpdateOrderReq struct {
	OrderID int64       `json:"OrderID" form:"orderId"` // 订单ID
	PetList []*OrderPet `json:"PetList" form:"petList"` // 包含修改后的宠物价格
	ToDate  common.Date `json:"ToDate" form:"toDate"`
}

type SitterConfirmModificationReq struct {
	OrderID int64                   `json:"OrderID" form:"orderId"` // 订单ID
	State   OrderModificationStatus `json:"State" form:"state"`
}

type UserConfirmModificationReq struct {
	OrderID int64                   `json:"OrderID" form:"orderId"` // 订单ID
	State   OrderModificationStatus `json:"State" form:"state"`
}
