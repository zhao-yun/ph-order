package model

import (
	"time"
)

// WalletTransaction 钱包收支记录明细表模型
type WalletTransaction struct {
	ID              int64      `gorm:"column:id;type:bigint;primaryKey;autoIncrement" json:"ID"`                            // 主键ID
	UserID          string     `gorm:"column:user_id;type:bigint;not null" json:"UserID"`                                   // 用户ID
	TransactionType int8       `gorm:"column:transaction_type;type:tinyint;not null" json:"TransactionType"`                // 交易类型：1-收入，2-支出
	OrderType       string     `gorm:"column:order_type;type:varchar(50);not null" json:"OrderType"`                        // 订单类型
	OrderID         *int64     `gorm:"column:order_id;type:bigint" json:"OrderID,omitempty"`                                // 关联订单ID（可选）
	OrderCreatedAt  *time.Time `gorm:"column:order_created_at;type:datetime" json:"OrderCreatedAt,omitempty"`               // 订单开始时间
	Amount          float64    `gorm:"column:amount;type:decimal(12,2);not null" json:"Amount"`                             // 交易金额（元，正数）
	BalanceAfter    float64    `gorm:"column:balance_after;type:decimal(12,2);not null" json:"BalanceAfter"`                // 交易后余额（元）
	TransactionTime time.Time  `gorm:"column:transaction_time;type:datetime;not null" json:"TransactionTime"`               // 交易时间
	Remark          *string    `gorm:"column:remark;type:varchar(255)" json:"Remark,omitempty"`                             // 交易备注
	CreatedAt       time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"CreatedAt"` // 记录创建时间
}

// TableName 指定表名
func (w WalletTransaction) TableName() string {
	return "wallet_transaction"
}

// TransactionType 交易类型枚举
type TransactionType int8

const (
	TransactionTypeIncome  TransactionType = 1 // 收入
	TransactionTypeExpense TransactionType = 2 // 支出
)

// 转换为字符串描述
func (t TransactionType) String() string {
	switch t {
	case TransactionTypeIncome:
		return "收入"
	case TransactionTypeExpense:
		return "支出"
	default:
		return "未知"
	}
}

// 订单类型常量（宠物服务相关场景）
const (
	OrderTypeBoarding   = "Boarding"    // 宠物寄宿
	OrderTypeDayCare    = "Day care"    // 日间照料
	OrderTypeDogWalking = "Dog Walking" // 遛狗服务
	OrderTypeDropIn     = "Drop-In"     // 上门照看
	// 可扩展其他类型
	OrderTypeRecharge = "Recharge" // 账户充值
	OrderTypeRefund   = "Refund"   // 订单退款
)

// 获取订单类型的中文描述（可选，用于展示）
func GetOrderTypeCN(orderType string) string {
	switch orderType {
	case OrderTypeBoarding:
		return "宠物寄宿"
	case OrderTypeDayCare:
		return "日间照料"
	case OrderTypeDogWalking:
		return "遛狗服务"
	case OrderTypeDropIn:
		return "上门照看"
	case OrderTypeRecharge:
		return "账户充值"
	case OrderTypeRefund:
		return "订单退款"
	default:
		return orderType
	}
}
