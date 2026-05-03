package model

import (
	"demo/util/common"
)

// OrderPet 订单宠物.
type OrderPet struct {
	ID        int64            `gorm:"primaryKey;autoIncrement;comment('主键ID')"`
	OrderID   int64            `gorm:"type:int;not null;comment('订单ID')"`
	OwnerID   string           `gorm:"type:varchar(50);not null;comment('用户ID')"`
	SitterID  string           `gorm:"type:varchar(50);not null;comment('Sitter ID')"`
	PetID     string           `gorm:"type:varchar(255);not null;comment('宠物ID')"`
	PetType   string           `gorm:"type:varchar(255);not null;comment('宠物类型，如猫、狗')"`
	PetShape  int64            `gorm:"not null;comment('宠物体型 1 为小型犬 2为中型犬 3为大型犬')"`
	PetName   string           `gorm:"type:varchar(255);comment('宠物名字')"`
	Breed     string           `gorm:"type:varchar(255);comment('宠物品种')"`
	PetPrice  float64          `gorm:"type:decimal(10,2);not null;comment('宠物费用')"`
	IsPuppy   bool             `gorm:"-" json:"IsPuppy"`
	CreatedAt common.Timestamp `gorm:"type:timestamp;autoCreateTime;comment('创建时间')"`
	UpdatedAt common.Timestamp `gorm:"type:timestamp;autoUpdateTime;comment('更新时间')"`
}

// TableName 设置表名
func (OrderPet) TableName() string {
	return "order_pet"
}
