package dal

import (
	"demo/model"
	"demo/util/postgres"

	"github.com/sirupsen/logrus"
)

func CreateSitterLocation(location *model.SitterLocation) error {
	db := postgres.GetDB()
	return db.Create(location).Error
}

func GetLatestSitterLocation(params *model.SitterLocationQueryParams) (*model.SitterLocation, error) {
	db := postgres.GetDB()
	query := db.Model(model.SitterLocation{})

	if params.SitterID != "" {
		query = query.Where("sitter_id = ?", params.SitterID)
	}
	if params.OrderID != 0 {
		query = query.Where("order_id = ?", params.OrderID)
	}
	if params.SubOrderID != 0 {
		query = query.Where("sub_order_id = ?", params.SubOrderID)
	}

	var location *model.SitterLocation
	err := query.Order("timestamp DESC").First(&location).Error
	if err != nil {
		logrus.Errorf("get latest sitter location failed, err = %v", err)
		return nil, err
	}
	return location, nil
}

func GetSitterLocationsBySubOrderID(subOrderID int64, limit int) ([]*model.SitterLocation, error) {
	db := postgres.GetDB()
	var locations []*model.SitterLocation
	err := db.Where("sub_order_id = ?", subOrderID).Order("timestamp DESC").Limit(limit).Find(&locations).Error
	if err != nil {
		logrus.Errorf("get sitter locations by sub order id failed, err = %v", err)
		return nil, err
	}
	return locations, nil
}
