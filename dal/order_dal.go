package dal

import (
	"math"

	"demo/model"
	"demo/util/postgres"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetOrderPage 分页查询.
func GetOrderPage(c *gin.Context, current int, size int, params *model.OrderQueryParams) (*model.OrderPage, error) {
	// 计算偏移量.
	offset := (current - 1) * size

	// 查询总记录数.
	var total int64
	db := postgres.GetDB()

	query := db.Model(model.Order{})
	if params.ID != 0 {
		query = query.Where("id = ?", params.ID)
	}
	if params.OwnerID != "" {
		query = query.Where("owner_id = ?", params.OwnerID)
	}
	if params.SitterID != "" {
		query = query.Where("sitter_id = ?", params.SitterID)
		query = query.Where("state not in ?", []model.OrderStatus{model.OrderInitialized})
	}
	if params.TypeList != nil && len(params.TypeList) > 0 {
		query = query.Where("type IN (?)", params.TypeList)
	}
	if params.State == -1 {
		query = query.Where("state not in ?", []model.OrderStatus{model.OrderRejected, model.OrderTimeout, model.OrderCancel})
	} else if params.State != 0 {
		query = query.Where("state = ?", params.State)
	}

	if params.OrderTime != nil {
		// 修改：根据传入的日期查询当天（从 00:00:00 到 23:59:59）的订单
		// 之前是：query.Where("created_at >= ?", params.OrderTime)
		
		// 计算传入日期的第二天零点，作为范围上限
		nextDay := params.OrderTime.Time.AddDate(0, 0, 1)
		
		// 使用 from_date 进行筛选，而不是 created_at
		query = query.Where("from_date >= ? AND from_date < ?", params.OrderTime, nextDay)
	}
	if params.StartTime != nil {
		// 使用 from_date 进行筛选，而不是 created_at
		query = query.Where("from_date >= ?", params.StartTime)
	}
	if params.EndTime != nil {
		// 因为 EndTime 如果是按照 "yyyy-MM-dd" 解析的，它是当天的 00:00:00，
		// 如果要求查询到 EndTime 当天结束，需要把 EndTime 加上一天。
		// 这里为了兼容，假设前端传的是日期，则将其作为小于第二天零点的条件
		nextDay := params.EndTime.Time.AddDate(0, 0, 1)
		// 使用 from_date 进行筛选，而不是 created_at
		query = query.Where("from_date < ?", nextDay)
	}
	if params.IDList != nil && len(params.IDList) > 0 {
		query = query.Where("id IN (?)", params.IDList)
	}

	if params.KeyIDList != nil && len(params.KeyIDList) > 0 {
		query = query.Where("id IN (?)", params.KeyIDList)
	}

	if params.SitterDeleted != nil {
		query = query.Where("sitter_deleted = ?", params.SitterDeleted)
	}
	if params.UserDeleted != nil {
		query = query.Where("user_deleted = ?", params.UserDeleted)
	}

	err := query.Count(&total).Error
	if err != nil {
		logrus.Errorf("[DB] count orders failed, err：%v", err)
		return nil, err
	}

	// 计算总页数.
	totalPage := int(math.Ceil(float64(total) / float64(size)))

	// 查询当前页数据.
	var orderList []*model.Order
	err = query.Offset(offset).Limit(size).Order("created_at desc").Find(&orderList).Error
	if err != nil {
		logrus.Errorf("[DB] get orders list failed, err：%v", err)
		return nil, err
	}

	// 查询订单宠物信息
	orderIDs := make([]int64, 0, len(orderList))
	for _, order := range orderList {
		orderIDs = append(orderIDs, order.ID)
	}

	if len(orderIDs) > 0 {

		// 查询订单宠物信息
		var allOrderPets []*model.OrderPet
		err = db.Where("order_id IN (?)", orderIDs).Find(&allOrderPets).Error
		if err != nil {
			logrus.Errorf("[DB] get order pets batch failed, err: %v", err)
			return nil, err
		}

		petMap := make(map[int64][]*model.OrderPet)
		for _, pet := range allOrderPets {
			petMap[pet.OrderID] = append(petMap[pet.OrderID], pet)
		}

		// 查询正在修改订单信息
		var allModificationLogs []*model.OrderModificationLog
		err = db.Where("order_id IN (?) AND state = ?", orderIDs, model.OrderInitialized).Find(&allModificationLogs).Error
		if err != nil {
			logrus.Errorf("[DB] get order modification logs batch failed, err: %v", err)
			return nil, err
		}

		logMap := make(map[int64]*model.OrderModificationLog)
		for _, log := range allModificationLogs {
			logMap[log.OrderID] = log
		}

		for _, order := range orderList {
			order.PetList = petMap[order.ID]
			order.ModificationLog = logMap[order.ID] // 如果没有日志会赋值为nil
		}

	}
	return &model.OrderPage{
		Total:     total,
		TotalPage: totalPage,
		Size:      size,
		Current:   current,
		OrderList: orderList,
	}, nil
}

// CreateOrder 创建一个新的订单
func CreateOrder(order *model.Order) (*model.Order, error) {
	db := postgres.GetDB()

	// 开启一个事务
	tx := db.Begin()
	if tx.Error != nil {
		logrus.Errorf("begin transaction failed, err = %v", tx.Error)
		return nil, tx.Error
	}

	// 在事务中执行创建操作
	err := tx.Create(&order).Error
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Errorf("rollback failed, err = %v", rollbackErr)
		}
		logrus.Errorf("create order failed, err = %v", err)
		return nil, err
	}

	// 创建订单宠物信息
	for _, pet := range order.PetList {
		pet.OrderID = order.ID
		pet.OwnerID = order.OwnerID
		pet.SitterID = order.SitterID
		err = tx.Create(&pet).Error
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logrus.Errorf("rollback failed, err = %v", rollbackErr)
			}
			logrus.Errorf("create order pet failed, err = %v", err)
			return nil, err
		}

	}

	// 提交事务
	if commitErr := tx.Commit().Error; commitErr != nil {
		logrus.Errorf("commit failed, err = %v", commitErr)
		return nil, commitErr
	}

	return order, nil
}

// CreateOrderWithSubOrders 创建订单并同时创建子订单（遛狗订单）
func CreateOrderWithSubOrders(order *model.Order, subOrders []*model.SubOrder) (*model.Order, error) {
	db := postgres.GetDB()

	tx := db.Begin()
	if tx.Error != nil {
		logrus.Errorf("begin transaction failed, err = %v", tx.Error)
		return nil, tx.Error
	}

	err := tx.Create(&order).Error
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Errorf("rollback failed, err = %v", rollbackErr)
		}
		logrus.Errorf("create order failed, err = %v", err)
		return nil, err
	}

	for _, pet := range order.PetList {
		pet.OrderID = order.ID
		pet.OwnerID = order.OwnerID
		pet.SitterID = order.SitterID
		err = tx.Create(&pet).Error
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logrus.Errorf("rollback failed, err = %v", rollbackErr)
			}
			logrus.Errorf("create order pet failed, err = %v", err)
			return nil, err
		}
	}

	for _, subOrder := range subOrders {
		subOrder.OrderID = order.ID
		err = tx.Create(&subOrder).Error
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logrus.Errorf("rollback failed, err = %v", rollbackErr)
			}
			logrus.Errorf("create sub order failed, err = %v", err)
			return nil, err
		}
	}

	if commitErr := tx.Commit().Error; commitErr != nil {
		logrus.Errorf("commit failed, err = %v", commitErr)
		return nil, commitErr
	}

	return order, nil
}

// UpdateOrder 更新一个现有的订单
func UpdateOrder(order *model.Order) error {
	db := postgres.GetDB()
	result := db.Save(order)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetOrderById 通过 ID 查询订单
func GetOrderById(id int64) (*model.Order, error) {
	db := postgres.GetDB()
	var order *model.Order
	err := db.First(&order, id).Error
	if err != nil {
		logrus.Errorf("get order by id failed, err = %v", err)
		return nil, err
	}
	if order == nil {
		return nil, nil
	}

	// 查询订单宠物信息
	var orderPets []*model.OrderPet
	err = db.Where("order_id =?", id).Find(&orderPets).Error
	if err != nil {
		logrus.Errorf("get order pet list failed, err = %v", err)
		return nil, err
	}
	order.PetList = orderPets

	// 查询宠物的正在更改信息
	var orderModificationLogs *model.OrderModificationLog
	err = db.Where("order_id =? AND state =?", id, model.OrderInitialized).First(&orderModificationLogs).Error
	if err != nil {
		logrus.Errorf("get order modification logs failed, err: %v", err)
	}
	order.ModificationLog = orderModificationLogs

	return order, nil
}

// GetUsernameByKeyword 根据关键词查询Sitter名字
func GetUsernameByKeyword(keyword string, userId string) ([]string, error) {
	db := postgres.GetDB()
	var sitters []string
	err := db.Model(&model.Order{}).Select("sitter_name").
		Distinct().
		Where("sitter_name LIKE ?", keyword+"%").
		Where(db.Where("owner_id != ?", userId).Or("sitter_id != ?", userId)).
		Find(&sitters).Error
	if err != nil {
		logrus.Errorf("get sitters by keyword failed, err: %v", err)
		return nil, err
	}
	return sitters, nil
}

func GetPetNamesByKeyword(keyword string, userId string) ([]string, error) {
	db := postgres.GetDB()
	var petNames []string
	err := db.Model(&model.OrderPet{}).Select("pet_name").
		Distinct().
		Where("pet_name LIKE ?", keyword+"%").
		Where(db.Where("owner_id != ?", userId).Or("sitter_id != ?", userId)).
		Find(&petNames).Error
	if err != nil {
		logrus.Errorf("get pet names by keyword failed, err: %v", err)
		return nil, err
	}
	return petNames, nil
}
