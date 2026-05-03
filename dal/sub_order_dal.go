package dal

import (
	"math"

	"demo/model"
	"demo/util/common"
	"demo/util/postgres"

	"github.com/sirupsen/logrus"
)

func GetSubOrderPage(current int, size int, params *model.SubOrderQueryParams) (*model.SubOrderPage, error) {
	offset := (current - 1) * size
	db := postgres.GetDB()

	query := db.Model(model.SubOrder{})
	if params.ID != 0 {
		query = query.Where("id = ?", params.ID)
	}
	if params.OrderID != 0 {
		query = query.Where("order_id = ?", params.OrderID)
	}
	if params.Date != nil {
		query = query.Where("date = ?", params.Date)
	}
	if params.State != 0 {
		query = query.Where("state = ?", params.State)
	}
	if params.StartTime != nil {
		query = query.Where("created_at >= ?", params.StartTime)
	}
	if params.EndTime != nil {
		query = query.Where("created_at <= ?", params.EndTime)
	}

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		logrus.Errorf("[DB] count sub orders failed, err：%v", err)
		return nil, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(size)))

	var subOrderList []*model.SubOrder
	err = query.Offset(offset).Limit(size).Order("date asc").Find(&subOrderList).Error
	if err != nil {
		logrus.Errorf("[DB] get sub orders list failed, err：%v", err)
		return nil, err
	}

	return &model.SubOrderPage{
		Total:        total,
		TotalPage:    totalPage,
		Size:         size,
		Current:      current,
		SubOrderList: subOrderList,
	}, nil
}

func CreateSubOrders(subOrders []*model.SubOrder) error {
	db := postgres.GetDB()
	return db.CreateInBatches(subOrders, len(subOrders)).Error
}

func GetSubOrderById(id int64) (*model.SubOrder, error) {
	db := postgres.GetDB()
	var subOrder *model.SubOrder
	err := db.First(&subOrder, id).Error
	if err != nil {
		logrus.Errorf("get sub order by id failed, err = %v", err)
		return nil, err
	}
	return subOrder, nil
}

func UpdateSubOrder(subOrder *model.SubOrder) error {
	db := postgres.GetDB()
	return db.Save(subOrder).Error
}

func GetSubOrdersByOrderID(orderID int64) ([]*model.SubOrder, error) {
	db := postgres.GetDB()
	var subOrderList []*model.SubOrder
	err := db.Where("order_id = ?", orderID).Order("date asc").Find(&subOrderList).Error
	if err != nil {
		logrus.Errorf("get sub orders by order id failed, err = %v", err)
		return nil, err
	}
	return subOrderList, nil
}

func GetTodayOrNearestSubOrder(orderID int64, today common.Date) (*model.SubOrder, int, error) {
	db := postgres.GetDB()
	
	var allSubOrders []*model.SubOrder
	err := db.Where("order_id = ?", orderID).Order("date asc").Find(&allSubOrders).Error
	if err != nil {
		logrus.Errorf("get sub orders failed, err = %v", err)
		return nil, 0, err
	}
	
	if len(allSubOrders) == 0 {
		return nil, 0, nil
	}
	
	var targetSubOrder *model.SubOrder
	var dayNumber int
	
	todayTime := today.Time
	
	for i, subOrder := range allSubOrders {
		subOrderTime := subOrder.Date.Time
		
		if !subOrderTime.Before(todayTime) {
			targetSubOrder = subOrder
			dayNumber = i + 1
			break
		}
	}
	
	if targetSubOrder == nil {
		targetSubOrder = allSubOrders[len(allSubOrders)-1]
		dayNumber = len(allSubOrders)
	}
	
	targetSubOrder.DayNumber = dayNumber
	
	return targetSubOrder, dayNumber, nil
}
