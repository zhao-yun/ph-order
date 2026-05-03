package model

import (
	"time"
)

// SitterRating 对应数据库中的sitter_rating表
type SitterRating struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement" json:"ID"`
	OrderID        int64     `gorm:"column:order_id;unique;not null" json:"OrderID"`
	UserID         string    `gorm:"column:user_id;type:varchar(255);not null" json:"UserID"`
	SitterID       string    `gorm:"column:sitter_id;type:varchar(255);not null" json:"SitterID"`
	Score          int8      `gorm:"column:score;type:tinyint" json:"Score,omitempty"`
	Punctuality    *int8     `gorm:"column:punctuality;type:tinyint" json:"Punctuality,omitempty"`
	Responsibility *int8     `gorm:"column:responsibility;type:tinyint" json:"Responsibility,omitempty"`
	Communication  *int8     `gorm:"column:communication;type:tinyint" json:"Communication,omitempty"`
	PetCareSkills  *int8     `gorm:"column:pet_care_skills;type:tinyint" json:"PetCareSkills,omitempty"`
	Cleanliness    *int8     `gorm:"column:cleanliness;type:tinyint" json:"Cleanliness,omitempty"`
	Suggestions    string    `gorm:"column:suggestions;type:text" json:"Suggestions,omitempty"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"CreatedAt,omitempty"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"UpdatedAt,omitempty"`
}

// TableName 指定表名
func (SitterRating) TableName() string {
	return "sitter_rating"
}

// SitterRatingPage 封装 sitter_rating 表的分页查询结果
type SitterRatingPage struct {
	Total     int64           `json:"Total`
	TotalPage int             `json:"TotalPage"`
	Current   int             `json:"Current"`
	Size      int             `json:"Size"`
	List      []*SitterRating `json:"list"` // 评价列表
}

type SitterRatingQueryParams struct {
	ID       int64  `form:"id" json:"id"`               // 评分ID
	OrderID  int64  `form:"order_id" json:"order_id"`   // 订单ID
	UserID   string `form:"user_id" json:"user_id"`     // 用户ID
	SitterID string `form:"sitter_id" json:"sitter_id"` // 宠物保姆ID
}
