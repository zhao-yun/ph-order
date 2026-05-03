package model

import (
	"demo/util/common"
)

type PetRating struct {
	ID        int              `gorm:"column:id;primary_key"`
	OrderID   int64            `gorm:"column:order_id"`
	PetID     string           `gorm:"column:pet_id"`
	PetName   string           `gorm:"column:pet_name"`
	Score     int8             `gorm:"column:score"`
	CreatedAt common.Timestamp `gorm:"column:created_at"`
	UpdatedAt common.Timestamp `gorm:"column:updated_at"`
}

func (PetRating) TableName() string {
	return "pet_rating" // 请将此处替换为实际的表名
}
