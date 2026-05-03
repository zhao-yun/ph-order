package dal

import (
	"demo/model"
	"demo/util/postgres"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateOrderPet 创建一个新的 OrderPet
func CreateOrderPet(orderPet model.OrderPet) (*model.OrderPet, error) {
	db := postgres.GetDB() // 假设 postgres.GetDB() 返回一个 *gorm.DB 实例
	err := db.Create(&orderPet).Error
	if err != nil {
		logrus.Errorf("create order pet failed, err = %v", err)
		return nil, err
	}
	return &orderPet, nil
}

// GetOrderPetByID 根据 ID 查询 OrderPet
func GetOrderPetByID(db *gorm.DB, id string) (*model.OrderPet, error) {
	var orderPet model.OrderPet
	err := db.Where("id = ?", id).First(&orderPet).Error
	if err != nil {
		logrus.Errorf("get order pet by ID failed, err = %v", err)
		return nil, err
	}
	return &orderPet, nil
}

// GetAllOrderPets 查询所有 OrderPet
func GetAllOrderPets(db *gorm.DB) ([]*model.OrderPet, error) {
	var orderPets []*model.OrderPet
	err := db.Find(&orderPets).Error
	if err != nil {
		logrus.Errorf("get all order pets failed, err = %v", err)
		return nil, err
	}
	return orderPets, nil
}

// UpdateOrderPet 更新一个现有的 OrderPet
func UpdateOrderPet(db *gorm.DB, orderPet *model.OrderPet) error {
	result := db.Save(orderPet)
	if result.Error != nil {
		logrus.Errorf("update order pet failed, err = %v", result.Error)
		return result.Error
	}
	return nil
}

// DeleteOrderPet 删除一个 OrderPet
func DeleteOrderPet(db *gorm.DB, id string) error {
	result := db.Where("id = ?", id).Delete(&model.OrderPet{})
	if result.Error != nil {
		logrus.Errorf("delete order pet failed, err = %v", result.Error)
		return result.Error
	}
	return nil
}

// GetOrderPetsWithPagination 分页查询 OrderPet
func GetOrderPetsWithPagination(db *gorm.DB, page, size int) ([]*model.OrderPet, int64, error) {
	var orderPets []*model.OrderPet
	var total int64

	// 计算偏移量
	offset := (page - 1) * size

	// 查询总数
	err := db.Model(&model.OrderPet{}).Count(&total).Error
	if err != nil {
		logrus.Errorf("count order pets failed, err = %v", err)
		return nil, 0, err
	}

	// 查询分页数据
	err = db.Offset(offset).Limit(size).Find(&orderPets).Error
	if err != nil {
		logrus.Errorf("get order pets with pagination failed, err = %v", err)
		return nil, 0, err
	}

	return orderPets, total, nil
}

// GetOrderIdByOrderPetList 根据订单宠物列表查询订单ID
func GetOrderIdByOrderPetList(params *model.OrderQueryParams) ([]int64, error) {
	db := postgres.GetDB()
	// 查询宠物类型存在
	orderPetQuery := db.Model(model.OrderPet{})
	if params.PetTypeList != nil && len(params.PetTypeList) > 0 {
		orderPetQuery = orderPetQuery.Where("pet_type IN (?)", params.PetTypeList)
	}
	if params.PetIDList != nil && len(params.PetIDList) > 0 {
		orderPetQuery = orderPetQuery.Where("pet_id IN (?)", params.PetIDList)
	}
	if params.OwnerID != "" {
		orderPetQuery = orderPetQuery.Where("owner_id = ?", params.OwnerID)
	}
	if params.SitterID != "" {
		orderPetQuery = orderPetQuery.Where("sitter_id = ?", params.SitterID)
	}
	var orderPet []model.OrderPet
	err := orderPetQuery.Model(model.OrderPet{}).Scan(&orderPet).Error
	if err != nil {
		logrus.Errorf("[DB] get order id list failed, err: %v", err)
		return nil, err
	}
	orderIdList := make([]int64, 0, len(orderPet))
	for _, pet := range orderPet {
		orderIdList = append(orderIdList, pet.OrderID)
	}
	orderIdList = append(orderIdList, -1)
	// 关联查询订单
	return orderIdList, nil
}

func GetOrderIdListByKeyword(keyword string) ([]int64, error) {
	db := postgres.GetDB()
	// 查询宠物类型存在
	orderPetQuery := db.Model(model.OrderPet{})
	if keyword != "" {
		orderPetQuery = orderPetQuery.Where("pet_name LIKE ?", "%"+keyword+"%").Or("pet_type LIKE ?", "%"+keyword+"%").Or("breed LIKE ?", "%"+keyword+"%")
	}
	var orderPet []model.OrderPet
	err := orderPetQuery.Model(model.OrderPet{}).Scan(&orderPet).Error
	if err != nil {
		logrus.Errorf("[DB] get order id list failed, err: %v", err)
		return nil, err
	}
	orderIdList := make([]int64, 0, len(orderPet))
	for _, pet := range orderPet {
		orderIdList = append(orderIdList, pet.OrderID)
	}
	orderIdList = append(orderIdList, -1)
	// 关联查询订单
	return orderIdList, nil
}
