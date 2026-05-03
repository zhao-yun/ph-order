package model

import (
	"demo/util/common"
)

// BankCard 银行卡模型
type BankCard struct {
	ID                    int64            `gorm:"primaryKey;column:id" json:"ID"`
	UserID                string           `gorm:"column:user_id;type:varchar(50);not null;index:idx_user_id" json:"UserID"`
	CardNumber            string           `gorm:"column:card_number;type:varchar(20);not null;index:idx_card_number" json:"CardNumber"`
	CardType              int8             `gorm:"column:card_type;type:tinyint(1);not null" json:"CardType"`
	BankCode              string           `gorm:"column:bank_code;type:varchar(10);not null" json:"BankCode"`
	BankName              string           `gorm:"column:bank_name;type:varchar(50);not null" json:"BankName"`
	InterbankTransferCode string           `gorm:"column:interbank_transfer_code;type:varchar(20)" json:"InterbankTransferCode"`
	AccountHolder         string           `gorm:"column:account_holder;type:varchar(50);not null" json:"AccountHolder"`
	IsDefault             int8             `gorm:"column:is_default;type:tinyint(1);default:0" json:"IsDefault"`
	Status                int8             `gorm:"column:status;type:tinyint(1);default:1" json:"Status"`
	CreatedAt             common.Timestamp `gorm:"column:created_at;type:datetime;not null;autoCreateTime" json:"CreatedAt"`
	UpdatedAt             common.Timestamp `gorm:"column:updated_at;type:datetime;not null;autoUpdateTime" json:"UpdatedAt"`
}

// TableName 指定表名
func (BankCard) TableName() string {
	return "bank_card"
}

// CardType 常量定义
const (
	CardTypeSavings  = 1 // 储蓄账户
	CardTypeChecking = 2 // 支票账户
)

// Status 常量定义
const (
	StatusDisabled = 0 // 禁用
	StatusEnabled  = 1 // 启用
)
