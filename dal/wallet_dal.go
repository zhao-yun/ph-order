package dal

import (
	"context"
	"database/sql"
	"time"

	"demo/model"
	"demo/util/common"

	"gorm.io/gorm"
)

// WalletDAL 钱包数据访问层
type WalletDAL struct {
	db *gorm.DB
}

// NewWalletDAL 创建钱包DAL实例
func NewWalletDAL(db *gorm.DB) *WalletDAL {
	return &WalletDAL{db: db}
}

// GetUserWalletByUserID 获取用户钱包信息
func (d *WalletDAL) GetUserWalletByUserID(ctx context.Context, userID string) (*model.UserWallet, error) {
	var wallet model.UserWallet
	result := d.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

// UpdateWalletBalance 更新钱包余额(使用事务保证原子性)
func (d *WalletDAL) UpdateWalletBalance(ctx context.Context, tx *gorm.DB, userID string, newBalance float64) error {
	db := d.db.WithContext(ctx)
	if tx != nil {
		db = tx
	}

	return db.Model(&model.UserWallet{}).
		Where("user_id = ?", userID).
		Update("balance", newBalance).Error
}

// CreateTransaction 创建交易记录
func (d *WalletDAL) CreateTransaction(ctx context.Context, tx *gorm.DB, transaction *model.WalletTransaction) error {
	db := d.db.WithContext(ctx)
	if tx != nil {
		db = tx
	}

	return db.Create(transaction).Error
}

// ListTransactions 分页查询交易记录
func (d *WalletDAL) ListTransactions(
	ctx context.Context,
	userID string,
	transactionType int8,
	orderTypeList []int64,
	startTime, endTime *common.Timestamp,
	createTime *common.Timestamp,
	timeType string, // "order_created" 或 "transaction"
	page, pageSize int,
) ([]*model.WalletTransaction, int64, error) {

	query := d.db.WithContext(ctx).
		Where("user_id = ?", userID)

	// 交易类型筛选
	if transactionType != 0 {
		query = query.Where("transaction_type = ?", transactionType)
	}

	// 订单类型筛选
	if orderTypeList != nil && len(orderTypeList) > 0 {
		query = query.Where("order_type IN ?", orderTypeList)
	}

	// 时间筛选
	if startTime != nil {
		query = query.Where("transaction_time >= ?", startTime)
	}
	if endTime != nil {
		// 如果传入的是当天的 00:00:00，表示查询到当天结束，所以加一天
		nextDay := endTime.Time.AddDate(0, 0, 1)
		query = query.Where("transaction_time < ?", nextDay)
	}
	if createTime != nil {
		// 最近时间，表示查询 createTime 及其之后的数据
		// 由于传进来的是某一天的 00:00:00，所以直接用 >= 即可
		query = query.Where("transaction_time >= ?", createTime)
	}

	// 计算总数
	var total int64
	if err := query.Model(&model.WalletTransaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var transactions []*model.WalletTransaction
	offset := (page - 1) * pageSize
	err := query.Order("transaction_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&transactions).Error

	return transactions, total, err
}

// GetTodayIncome 获取今日收益
func (d *WalletDAL) GetTodayIncome(ctx context.Context, userID string) (float64, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var sum struct {
		Total float64 `gorm:"column:total"`
	}

	err := d.db.WithContext(ctx).
		Table("wallet_transaction").
		Select("SUM(amount) as total").
		Where("user_id = ?", userID).
		Where("transaction_type = ?", model.TransactionTypeIncome).
		Where("transaction_time BETWEEN ? AND ?", startOfDay, endOfDay).
		Scan(&sum).Error

	if err != nil && err != sql.ErrNoRows {
		return 0.0, err
	}

	return sum.Total, nil
}

// GetPendingSettlement 获取未结算收入
func (d *WalletDAL) GetPendingSettlement(ctx context.Context, userID string) (float64, error) {
	// 假设未结算收入存储在单独的表中，这里简化处理
	// 实际应用中应根据业务逻辑查询
	var wallet model.UserWallet
	result := d.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&wallet)

	if result.Error != nil {
		return 0.0, result.Error
	}

	// 这里假设未结算收入是一个独立字段
	// 实际应用中可能需要关联其他表计算
	return 0.0, nil
}

// BeginTx 开启事务
func (d *WalletDAL) BeginTx(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx).Begin()
}
