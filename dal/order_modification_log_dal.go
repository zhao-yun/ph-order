package dal

import (
	"errors"
	"fmt"

	"demo/model"
	"demo/util/json"
	"demo/util/postgres"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateOrderModificationLog 创建一个新的 OrderModificationLog 或覆盖未完成的修改
func CreateOrderModificationLog(log model.OrderModificationLog) (*model.OrderModificationLog, error) {
	db := postgres.GetDB()

	// 查询该订单是否存在由同一角色发起、且状态为初始化的修改记录
	var existingLog model.OrderModificationLog
	err := db.Where("order_id = ? AND type = ? AND state = ?", log.OrderID, log.Type, model.OrderModificationInitialized).First(&existingLog).Error

	if err == nil {
		// 记录存在，覆盖原来的记录
		log.ID = existingLog.ID
		err = db.Save(&log).Error
		if err != nil {
			logrus.Errorf("update existing order modification log failed, err = %v", err)
			return nil, err
		}
		return &log, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 记录不存在，创建新记录
		err = db.Create(&log).Error
		if err != nil {
			logrus.Errorf("create order modification log failed, err = %v", err)
			return nil, err
		}
		return &log, nil
	}

	// 发生其他数据库错误
	logrus.Errorf("query existing order modification log failed, err = %v", err)
	return nil, err
}

// GetOrderModificationLogByID 根据 ID 查询 OrderModificationLog
func GetOrderModificationLogByID(db *gorm.DB, id int64) (*model.OrderModificationLog, error) {
	var log model.OrderModificationLog
	err := db.Where("id = ?", id).First(&log).Error
	if err != nil {
		logrus.Errorf("get order modification log by ID failed, err = %v", err)
		return nil, err
	}
	return &log, nil
}

// GetIngOrderModificationLogByOrderID 根据 OrderID 查询一个正在修改中的 OrderModificationLog
func GetIngOrderModificationLogByOrderID(orderID int64) (*model.OrderModificationLog, error) {
	db := postgres.GetDB()
	var log model.OrderModificationLog
	err := db.Where("order_id =? AND state =?", orderID, model.OrderModificationInitialized).First(&log).Error
	if err != nil {
		logrus.Errorf("get order modification log by order ID failed, err = %v", err)
		return nil, err
	}
	return &log, nil
}

// GetAllOrderModificationLogs 查询所有 OrderModificationLog
func GetAllOrderModificationLogs(db *gorm.DB) ([]*model.OrderModificationLog, error) {
	var logs []*model.OrderModificationLog
	err := db.Find(&logs).Error
	if err != nil {
		logrus.Errorf("get all order modification logs failed, err = %v", err)
		return nil, err
	}
	return logs, nil
}

func UpdateOrderModificationLog(log *model.OrderModificationLog) error {
	db := postgres.GetDB()

	return db.Transaction(func(tx *gorm.DB) error {
		// ---------------------- 阶段1：更新订单修改日志 ----------------------
		logrus.Infof("开始更新订单修改日志，ID: %d", log.ID)
		if err := tx.Save(log).Error; err != nil {
			logrus.WithError(err).Errorf("更新订单修改日志失败，ID: %d", log.ID)
			return fmt.Errorf("更新订单修改日志失败: %w", err) // 包装错误，保留原始信息
		}

		// ---------------------- 阶段2：处理非接受状态的日志 ----------------------
		if log.State != model.OrderModificationAccepted {
			logrus.Infof("订单修改状态为非接受，无需更新订单，ID: %d", log.ID)
			return nil // 提前返回，事务自动提交（GORM Transaction机制）
		}

		// ---------------------- 阶段3：查询并验证订单存在性 ----------------------
		var order model.Order
		logrus.Infof("查询订单信息，OrderID: %d", log.OrderID)
		if err := tx.First(&order, log.OrderID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.WithField("OrderID", log.OrderID).Error("订单不存在")
				return fmt.Errorf("订单不存在: OrderID=%d", log.OrderID)
			}
			logrus.WithError(err).Errorf("查询订单失败，OrderID: %d", log.OrderID)
			return fmt.Errorf("查询订单失败: %w", err)
		}

		// ---------------------- 阶段4：更新订单主信息 ----------------------
		if log.Type == model.OrderModificationTypeUser {
			logrus.Infof("更新订单时间，OrderID: %d", order.ID)
			order.ToDate = log.NewDate
			if err := tx.Save(&order).Error; err != nil {
				logrus.WithError(err).Errorf("更新订单时间失败，OrderID: %d", order.ID)
				return fmt.Errorf("更新订单时间失败: %w", err)
			}
		} else if log.Type == model.OrderModificationTypeSitter {
			logrus.Infof("更新订单时间和价格，OrderID: %d", order.ID)
			if !log.NewDate.Time.IsZero() {
				order.ToDate = log.NewDate
			}
			if log.NewPrice > 0 {
				order.TotalPrice = log.NewPrice
			}
			if err := tx.Save(&order).Error; err != nil {
				logrus.WithError(err).Errorf("更新订单时间和价格失败，OrderID: %d", order.ID)
				return fmt.Errorf("更新订单时间和价格失败: %w", err)
			}
		}

		// ---------------------- 阶段5：处理宠物信息（先删除后新增） ----------------------
		if log.Type == model.OrderModificationTypeUser {
			logrus.Infof("开始处理订单宠物信息，OrderID: %d", order.ID)

			// 删除原有宠物记录
			if err := tx.Where("order_id = ?", order.ID).Delete(&model.OrderPet{}).Error; err != nil {
				logrus.WithError(err).Errorf("删除原有宠物记录失败，OrderID: %d", order.ID)
				return fmt.Errorf("删除原有宠物记录失败: %w", err)
			}
			logrus.Infof("成功删除原有宠物记录，OrderID: %d", order.ID)

			// 解析新增宠物列表
			var petList []*model.OrderPet
			if err := json.Unmarshal([]byte(log.NewPetList), &petList); err != nil {
				logrus.WithError(err).Errorf("解析宠物JSON失败，OrderID: %d", order.ID)
				return fmt.Errorf("解析宠物JSON失败: %w", err)
			}

			// 批量创建新宠物记录（使用CreateInBatches提升性能）
			if len(petList) > 0 {
				for _, pet := range petList {
					pet.OrderID = order.ID // 绑定订单ID
					pet.OwnerID = order.OwnerID
					pet.SitterID = order.SitterID
				}
				logrus.Infof("开始批量创建宠物记录，数量: %d", len(petList))
				if err := tx.CreateInBatches(petList, 50).Error; err != nil { // 每批50条，可调整
					logrus.WithError(err).Errorf("批量创建宠物记录失败，OrderID: %d", order.ID)
					return fmt.Errorf("批量创建宠物记录失败: %w", err)
				}
				logrus.Infof("成功创建宠物记录，数量: %d", len(petList))
			}
		}

		logrus.Infof("订单修改日志更新完成，ID: %d", log.ID)
		return nil // 事务自动提交
	})
}

// DeleteOrderModificationLog 删除一个 OrderModificationLog
func DeleteOrderModificationLog(db *gorm.DB, id int64) error {
	result := db.Where("id = ?", id).Delete(&model.OrderModificationLog{})
	if result.Error != nil {
		logrus.Errorf("delete order modification log failed, err = %v", result.Error)
		return result.Error
	}
	return nil
}

// GetOrderModificationLogsByOrderID 根据 OrderID 查询所有 OrderModificationLog
func GetOrderModificationLogsByOrderID(db *gorm.DB, orderID int64) ([]*model.OrderModificationLog, error) {
	var logs []*model.OrderModificationLog
	err := db.Where("order_id = ?", orderID).Find(&logs).Error
	if err != nil {
		logrus.Errorf("get order modification logs by order ID failed, err = %v", err)
		return nil, err
	}
	return logs, nil
}

// GetOrderModificationLogsWithPagination 分页查询 OrderModificationLog
func GetOrderModificationLogsWithPagination(db *gorm.DB, page, size int) ([]*model.OrderModificationLog, int64, error) {
	var logs []*model.OrderModificationLog
	var total int64

	// 计算偏移量
	offset := (page - 1) * size

	// 查询总数
	err := db.Model(&model.OrderModificationLog{}).Count(&total).Error
	if err != nil {
		logrus.Errorf("count order modification logs failed, err = %v", err)
		return nil, 0, err
	}

	// 查询分页数据
	err = db.Offset(offset).Limit(size).Find(&logs).Error
	if err != nil {
		logrus.Errorf("get order modification logs with pagination failed, err = %v", err)
		return nil, 0, err
	}

	return logs, total, nil
}
