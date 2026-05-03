package service

import (
	"context"
	"errors"
	"sync"

	"demo/dal"
	"demo/model"
	"demo/util/common"
	"demo/util/postgres"
)

// WalletService 钱包服务
type WalletService struct {
	walletDAL *dal.WalletDAL
	mu        sync.RWMutex // 读写锁保证线程安全
}

// NewWalletService 创建钱包服务实例
func NewWalletService() *WalletService {
	walletDAL := dal.NewWalletDAL(postgres.GetDB())
	return &WalletService{
		walletDAL: walletDAL,
	}
}

// WalletBasicInfo 钱包基本信息
type WalletBasicInfo struct {
	Balance           float64 `json:"Balance"`           // 账户余额
	TodayIncome       float64 `json:"TodayIncome"`       // 今日收益
	PendingSettlement float64 `json:"PendingSettlement"` // 未结算收入
}

// GetWalletBasicInfo 获取钱包基本信息
func (s *WalletService) GetWalletBasicInfo(ctx context.Context, userID string) (*WalletBasicInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 获取账户余额
	wallet, err := s.walletDAL.GetUserWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 获取今日收益
	todayIncome, err := s.walletDAL.GetTodayIncome(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 获取未结算收入
	pendingSettlement, err := s.walletDAL.GetPendingSettlement(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &WalletBasicInfo{
		Balance:           wallet.Balance,
		TodayIncome:       todayIncome,
		PendingSettlement: pendingSettlement,
	}, nil
}

// TransactionFilter 交易记录筛选条件
type TransactionFilter struct {
	UserID          string            `json:"UserID"`
	TransactionType int8              `json:"TransactionType"` // 1-收入，2-支出，0-全部
	OrderTypeList   []int64           `json:"OrderType"`       // 订单类型
	StartTime       *common.Timestamp `json:"StartTime"`
	EndTime         *common.Timestamp `json:"EndTime"`
	CreateTime      *common.Timestamp `json:"CreateTime"`
	TimeType        string            `json:"TimeType"` // "order_created" 或 "transaction"
	Page            int               `json:"Page"`     // 页码，从1开始
	PageSize        int               `json:"PageSize"` // 每页条数

}

// TransactionListResult 交易记录列表结果
type TransactionListResult struct {
	Total        int64                      `json:"Total"`        // 总条数
	Transactions []*model.WalletTransaction `json:"Transactions"` // 交易记录列表
}

// ListTransactions 获取交易记录列表
func (s *WalletService) ListTransactions(ctx context.Context, filter *TransactionFilter) (*TransactionListResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	transactions, total, err := s.walletDAL.ListTransactions(
		ctx,
		filter.UserID,
		filter.TransactionType,
		filter.OrderTypeList,
		filter.StartTime,
		filter.EndTime,
		filter.CreateTime,
		filter.TimeType,
		filter.Page,
		filter.PageSize,
	)

	if err != nil {
		return nil, err
	}

	return &TransactionListResult{
		Total:        total,
		Transactions: transactions,
	}, nil
}

// CreateTransaction 创建交易记录(线程安全)
func (s *WalletService) CreateTransaction(ctx context.Context, transaction *model.WalletTransaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 开启事务
	tx := s.walletDAL.BeginTx(ctx)
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建交易记录
	if err := s.walletDAL.CreateTransaction(ctx, tx, transaction); err != nil {
		tx.Rollback()
		return err
	}

	// 更新钱包余额
	if err := s.walletDAL.UpdateWalletBalance(ctx, tx, transaction.UserID, transaction.BalanceAfter); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// Withdraw 提现
func (s *WalletService) Withdraw(ctx context.Context, userID string, param *model.WithdrawReq) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 校验金额
	if param.Amount <= 0 {
		return errors.New("提现金额必须大于0")
	}

	// 校验用户钱包状态
	wallet, err := s.walletDAL.GetUserWalletByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// 校验余额是否充足
	if wallet.Balance < param.Amount {
		return errors.New("余额不足，无法提现")
	}

	// 扣减
	wallet.Balance -= param.Amount
	// 更新钱包余额
	if err := s.walletDAL.UpdateWalletBalance(ctx, postgres.GetDB(), wallet.UserID, wallet.Balance); err != nil {
		return err
	}

	return nil

}
