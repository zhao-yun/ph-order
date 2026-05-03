package dal

import (
	"context"
	"errors"

	"demo/model"
	"demo/util/postgres"

	"gorm.io/gorm"
)

// Create 创建银行卡记录
func Create(ctx context.Context, card *model.BankCard) error {
	db := postgres.GetDB()
	return db.WithContext(ctx).Create(card).Error
}

// GetByID 根据ID获取银行卡
func GetByID(ctx context.Context, id int64) (*model.BankCard, error) {
	db := postgres.GetDB()
	var card model.BankCard
	err := db.WithContext(ctx).Where("id = ?", id).First(&card).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &card, err
}

// GetByUserID 获取用户的所有银行卡
func GetByUserID(ctx context.Context, userID string) ([]*model.BankCard, error) {
	db := postgres.GetDB()
	var cards []*model.BankCard
	err := db.WithContext(ctx).Where("user_id = ?", userID).Find(&cards).Error
	return cards, err
}

// GetDefaultCard 获取用户的默认银行卡
func GetDefaultCard(ctx context.Context, userID string) (*model.BankCard, error) {
	db := postgres.GetDB()
	var card model.BankCard
	err := db.WithContext(ctx).Where("user_id = ? AND is_default = 1", userID).First(&card).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &card, err
}

// Update 更新银行卡信息
func Update(ctx context.Context, card *model.BankCard) error {
	db := postgres.GetDB()
	return db.WithContext(ctx).Save(card).Error
}

// Delete 删除银行卡
func Delete(ctx context.Context, id int64) error {
	db := postgres.GetDB()
	return db.WithContext(ctx).Delete(&model.BankCard{}, id).Error
}

// SetDefault 设置用户默认银行卡
func SetDefault(ctx context.Context, userID string, cardID int64) error {
	db := postgres.GetDB()

	// 开启事务
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先取消所有卡的默认状态
		if err := tx.Model(&model.BankCard{}).
			Where("user_id = ?", userID).
			Update("is_default", 0).Error; err != nil {
			return err
		}

		// 设置指定卡为默认
		if err := tx.Model(&model.BankCard{}).
			Where("id = ? AND user_id = ?", cardID, userID).
			Update("is_default", 1).Error; err != nil {
			return err
		}

		return nil
	})
}

// CheckCardExists 检查银行卡是否已存在
func CheckCardExists(ctx context.Context, userID, cardNumber string) (bool, error) {
	db := postgres.GetDB()
	var count int64
	err := db.WithContext(ctx).Model(&model.BankCard{}).
		Where("user_id = ? AND card_number = ?", userID, cardNumber).
		Count(&count).Error
	return count > 0, err
}

// ListActiveCards 获取用户所有启用的银行卡
func ListActiveCards(ctx context.Context, userID string) ([]*model.BankCard, error) {
	db := postgres.GetDB()
	var cards []*model.BankCard
	err := db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, model.StatusEnabled).
		Find(&cards).Error
	return cards, err
}

// DisableCard 禁用银行卡
func DisableCard(ctx context.Context, cardID int64, userID string) error {
	db := postgres.GetDB()
	return db.WithContext(ctx).
		Model(&model.BankCard{}).
		Where("id = ? AND user_id = ?", cardID, userID).
		Update("status", model.StatusDisabled).Error
}
