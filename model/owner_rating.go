package model

import (
	"demo/util/common"
)

type OwnerRating struct {
	ID                  int              `gorm:"column:id;primary_key"`
	OrderID             int64            `gorm:"column:order_id"`
	OwnerID             string           `gorm:"column:owner_id"`
	SitterID            string           `gorm:"column:sitter_id"`
	OwnerName           string           `gorm:"column:owner_name"`
	SitterName          string           `gorm:"column:sitter_name"`
	Score               *int8            `gorm:"column:score"`
	SatisfactionLevel   *int8            `gorm:"column:satisfaction_level"`
	InstructionsClarity *int8            `gorm:"column:instructions_clarity"`
	Communication       *int8            `gorm:"column:communication"`
	SuppliesPreparation *int8            `gorm:"column:supplies_preparation"`
	RespectCourtesy     *int8            `gorm:"column:respect_courtesy"`
	Suggestions         string           `gorm:"column:suggestions"`
	CreatedAt           common.Timestamp `gorm:"column:created_at"`
	UpdatedAt           common.Timestamp `gorm:"column:updated_at"`
	PetRatingList       []*PetRating     `gorm:"-"`
}

func (OwnerRating) TableName() string {
	return "owner_rating"
}
