package model

import (
	"time"
)

// UserWallet 用户钱包主表模型
type UserWallet struct {
	ID        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement" json:"ID"`                                                        // 主键ID
	UserID    string    `gorm:"column:user_id;type:bigint;not null;uniqueIndex:uk_user_id" json:"UserID"`                                        // 用户ID（关联用户表）
	Balance   float64   `gorm:"column:balance;type:decimal(12,2);not null;default:0.00" json:"Balance"`                                          // 账户余额（元）
	Status    int8      `gorm:"column:status;type:tinyint;not null;default:1" json:"Status"`                                                     // 钱包状态：1-正常，2-冻结
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"CreatedAt"`                             // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"UpdatedAt"` // 更新时间 	// 软删除字段（可选）
}

// TableName 指定表名
func (u UserWallet) TableName() string {
	return "user_wallet"
}

// WalletStatus 钱包状态枚举
type WalletStatus int8

const (
	WalletStatusNormal WalletStatus = 1 // 正常
	WalletStatusFrozen WalletStatus = 2 // 冻结
)

// 转换为字符串描述
func (s WalletStatus) String() string {
	switch s {
	case WalletStatusNormal:
		return "正常"
	case WalletStatusFrozen:
		return "冻结"
	default:
		return "未知"
	}
}

type WithdrawReq struct {
	Amount float64 `json:"Amount"` // 提现金额
}
