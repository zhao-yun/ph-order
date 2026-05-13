package dal

import (
	"demo/model"
	"demo/util/postgres"

	"github.com/sirupsen/logrus"
)

// CreateOwnerRating 创建主人评分（含宠物评分）
func CreateOwnerRating(rating *model.OwnerRating) (*model.OwnerRating, error) {
	db := postgres.GetDB()

	tx := db.Begin()
	if tx.Error != nil {
		logrus.Errorf("begin transaction failed, err: %v", tx.Error)
		return nil, tx.Error
	}

	// 1. 保存主人评分
	err := tx.Create(rating).Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("create owner rating failed, err: %v", err)
		return nil, err
	}

	// 2. 保存关联的宠物评分（如果有）
	if len(rating.PetRatingList) > 0 {
		for i := range rating.PetRatingList {
			rating.PetRatingList[i].OrderID = rating.OrderID
			err := tx.Create(&rating.PetRatingList[i]).Error
			if err != nil {
				tx.Rollback()
				logrus.Errorf("create pet rating failed, err: %v", err)
				return nil, err
			}
		}
	}

	// 3. 更新订单状态为已评价
	err = tx.Model(&model.Order{}).Where("id = ?", rating.OrderID).Update("user_rating_state", 1).Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("update order state failed, err: %v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		logrus.Errorf("commit transaction failed, err: %v", err)
		return nil, err
	}

	return rating, nil
}

// UpdateOwnerRating 更新主人评分
func UpdateOwnerRating(rating *model.OwnerRating) error {
	db := postgres.GetDB()

	tx := db.Begin()
	if tx.Error != nil {
		logrus.Errorf("begin transaction failed, err: %v", tx.Error)
		return tx.Error
	}

	// 1. 更新主人评分
	updateData := make(map[string]interface{})
	if rating.Score != nil {
		updateData["score"] = rating.Score
	}
	if rating.SatisfactionLevel != nil {
		updateData["satisfaction_level"] = rating.SatisfactionLevel
	}
	if rating.InstructionsClarity != nil {
		updateData["instructions_clarity"] = rating.InstructionsClarity
	}
	if rating.Communication != nil {
		updateData["communication"] = rating.Communication
	}
	if rating.SuppliesPreparation != nil {
		updateData["supplies_preparation"] = rating.SuppliesPreparation
	}
	if rating.RespectCourtesy != nil {
		updateData["respect_courtesy"] = rating.RespectCourtesy
	}
	if rating.Suggestions != "" {
		updateData["suggestions"] = rating.Suggestions
	}

	result := tx.Model(&model.OwnerRating{}).Where("id = ?", rating.ID).Updates(updateData)
	if result.Error != nil {
		tx.Rollback()
		logrus.Errorf("update owner rating failed, err: %v", result.Error)
		return result.Error
	}

	// 删除旧宠物评分
	result = tx.Where("order_id = ?", rating.OrderID).Delete(&model.PetRating{})
	if result.Error != nil {
		tx.Rollback()
		logrus.Errorf("delete associated pet ratings failed, err: %v", result.Error)
		return result.Error
	}

	// 2. 保存关联的宠物评分（如果有）
	if len(rating.PetRatingList) > 0 {
		for i := range rating.PetRatingList {
			rating.PetRatingList[i].OrderID = rating.OrderID
			err := tx.Create(&rating.PetRatingList[i]).Error
			if err != nil {
				tx.Rollback()
				logrus.Errorf("create pet rating failed, err: %v", err)
				return err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		logrus.Errorf("commit transaction failed, err: %v", err)
		return err
	}

	return nil
}

// DeleteOwnerRating 删除主人评分（级联删除关联的宠物评分）
func DeleteOwnerRating(id int) error {
	db := postgres.GetDB()

	tx := db.Begin()
	if tx.Error != nil {
		logrus.Errorf("begin transaction failed, err: %v", tx.Error)
		return tx.Error
	}

	var rating model.OwnerRating
	result := tx.Where("id = ?", id).First(&rating)
	if result.Error != nil {
		tx.Rollback()
		logrus.Errorf("get owner rating failed before delete, err: %v", result.Error)
		return result.Error
	}

	// 1. 先删除关联的宠物评分
	result = tx.Where("owner_rating_id = ?", id).Delete(&model.PetRating{})
	if result.Error != nil {
		tx.Rollback()
		logrus.Errorf("delete associated pet ratings failed, err: %v", result.Error)
		return result.Error
	}

	// 2. 删除主人评分
	result = tx.Where("id = ?", id).Delete(&model.OwnerRating{})
	if result.Error != nil {
		tx.Rollback()
		logrus.Errorf("delete owner rating failed, err: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil // 可选择返回 ErrRecordNotFound
	}

	if err := tx.Model(&model.Order{}).Where("id = ?", rating.OrderID).Update("user_rating_state", 0).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("reset owner rating state failed, err: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		logrus.Errorf("commit transaction failed, err: %v", err)
		return err
	}

	return nil
}

// GetOwnerRatingById 通过 ID 查询主人评分（含宠物评分）
func GetOwnerRatingById(id int) (*model.OwnerRating, error) {
	db := postgres.GetDB()
	var rating model.OwnerRating

	// 查询主人评分
	err := db.First(&rating, id).Error
	if err != nil {
		logrus.Errorf("get owner rating by id failed, err: %v", err)
		return nil, err
	}

	// 查询关联的宠物评分
	var petRatings []*model.PetRating
	err = db.Where("order_id = ?", rating.OwnerID).Find(&petRatings).Error
	if err != nil {
		logrus.Errorf("get associated pet ratings failed, err: %v", err)
		return nil, err
	}

	rating.PetRatingList = petRatings
	return &rating, nil
}

// GetAverageOwnerRating 获取用户作为Owner的平均分
func GetAverageOwnerRating(ownerID string) (float64, error) {
	db := postgres.GetDB()
	var avg float64
	err := db.Model(&model.OwnerRating{}).Where("owner_id = ?", ownerID).Select("COALESCE(AVG(score), 0)").Scan(&avg).Error
	if err != nil {
		logrus.Errorf("get average owner rating failed, err: %v", err)
		return 0, err
	}
	return avg, nil
}

// GetOwnerRatingByOrderId 通过订单 ID 查询主人评分（含宠物评分）
func GetOwnerRatingByOrderId(orderId int64) (*model.OwnerRating, error) {
	db := postgres.GetDB()
	var rating model.OwnerRating

	// 查询主人评分
	err := db.Where("order_id = ?", orderId).First(&rating).Error
	if err != nil {
		logrus.Errorf("get owner rating by order id failed, err: %v", err)
		return nil, err
	}

	// 查询关联的宠物评分
	var petRatings []*model.PetRating
	err = db.Where("order_id = ?", orderId).Find(&petRatings).Error
	if err != nil {
		logrus.Errorf("get associated pet ratings failed, err: %v", err)
		return nil, err
	}

	rating.PetRatingList = petRatings
	return &rating, nil
}
