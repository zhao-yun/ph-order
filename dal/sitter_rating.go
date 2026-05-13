package dal

import (
	"demo/model"
	"demo/util/postgres"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetSitterRatingPage 分页查询 sitter_rating 表
func GetSitterRatingPage(c *gin.Context, current int, size int, params model.SitterRatingQueryParams) (*model.SitterRatingPage, error) {
	offset := (current - 1) * size
	var total int64
	db := postgres.GetDB()

	query := db.Model(&model.SitterRating{})
	if params.ID != 0 {
		query = query.Where("id = ?", params.ID)
	}
	if params.OrderID != 0 {
		query = query.Where("order_id = ?", params.OrderID)
	}
	if params.UserID != "" && params.SitterID != "" {
		query = query.Where(db.Where("user_id = ?", params.UserID).Or("sitter_id = ?", params.SitterID))
	} else if params.UserID != "" {
		query = query.Where("user_id = ?", params.UserID)
	} else if params.SitterID != "" {
		query = query.Where("sitter_id = ?", params.SitterID)
	}

	err := query.Count(&total).Error
	if err != nil {
		logrus.Errorf("[DB] count sitter ratings failed, err: %v", err)
		return nil, err
	}

	totalPage := int((total + int64(size) - 1) / int64(size)) // 等价于 math.Ceil

	var ratingList []*model.SitterRating
	err = query.Offset(offset).Limit(size).Order("created_at desc").Find(&ratingList).Error
	if err != nil {
		logrus.Errorf("[DB] get sitter ratings list failed, err: %v", err)
		return nil, err
	}

	return &model.SitterRatingPage{
		Total:     total,
		TotalPage: totalPage,
		Size:      size,
		Current:   current,
		List:      ratingList,
	}, nil
}

// CreateSitterRating 创建 sitter_rating 记录
func CreateSitterRating(rating *model.SitterRating) (*model.SitterRating, error) {
	db := postgres.GetDB()

	tx := db.Begin()
	if tx.Error != nil {
		logrus.Errorf("begin transaction failed, err: %v", tx.Error)
		return nil, tx.Error
	}

	err := tx.Create(rating).Error
	if err != nil {
		tx.Rollback()
		logrus.Errorf("create sitter rating failed, err: %v", err)
		return nil, err
	}

	// 3. 更新订单状态为已评价
	err = tx.Model(&model.Order{}).Where("id = ?", rating.OrderID).Update("sitter_rating_state", 1).Error
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

// UpdateSitterRating 更新 sitter_rating 记录
func UpdateSitterRating(rating *model.SitterRating) error {
	db := postgres.GetDB()

	// 使用 map 避免更新零值字段
	updateData := make(map[string]interface{})
	if rating.Score != 0 {
		updateData["score"] = rating.Score
	}
	if rating.Punctuality != nil {
		updateData["punctuality"] = rating.Punctuality
	}
	if rating.Responsibility != nil {
		updateData["responsibility"] = rating.Responsibility
	}
	if rating.Communication != nil {
		updateData["communication"] = rating.Communication
	}
	if rating.PetCareSkills != nil {
		updateData["pet_care_skills"] = rating.PetCareSkills
	}
	if rating.Cleanliness != nil {
		updateData["cleanliness"] = rating.Cleanliness
	}
	if rating.Suggestions != "" {
		updateData["suggestions"] = rating.Suggestions
	}

	result := db.Model(&model.SitterRating{}).Where("id = ?", rating.ID).Updates(updateData)
	if result.Error != nil {
		logrus.Errorf("update sitter rating failed, err: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return nil // 可以选择返回 ErrRecordNotFound
	}

	return nil
}

// DeleteSitterRating 删除 sitter_rating 记录
func DeleteSitterRating(id int) error {
	db := postgres.GetDB()

	tx := db.Begin()
	if tx.Error != nil {
		logrus.Errorf("begin transaction failed, err: %v", tx.Error)
		return tx.Error
	}

	var rating model.SitterRating
	result := tx.Where("id = ?", id).First(&rating)
	if result.Error != nil {
		tx.Rollback()
		logrus.Errorf("get sitter rating failed before delete, err: %v", result.Error)
		return result.Error
	}

	result = tx.Where("id = ?", id).Delete(&model.SitterRating{})
	if result.Error != nil {
		tx.Rollback()
		logrus.Errorf("delete sitter rating failed, err: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil // 可以选择返回 ErrRecordNotFound
	}

	if err := tx.Model(&model.Order{}).Where("id = ?", rating.OrderID).Update("sitter_rating_state", 0).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("reset sitter rating state failed, err: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		logrus.Errorf("commit transaction failed, err: %v", err)
		return err
	}

	return nil
}

// GetSitterRatingById 通过 ID 查询 sitter_rating
func GetSitterRatingById(id int) (*model.SitterRating, error) {
	db := postgres.GetDB()
	var rating model.SitterRating

	err := db.First(&rating, id).Error
	if err != nil {
		logrus.Errorf("get sitter rating by id failed, err: %v", err)
		return nil, err
	}

	return &rating, nil
}

// GetSitterRatingByOrderId 通过订单 ID 查询 sitter_rating
func GetSitterRatingByOrderId(orderId int64) (*model.SitterRating, error) {
	db := postgres.GetDB()
	var rating model.SitterRating

	err := db.Where("order_id = ?", orderId).First(&rating).Error
	if err != nil {
		logrus.Errorf("get sitter rating by order id failed, err: %v", err)
		return nil, err
	}

	return &rating, nil
}

// GetAverageSitterRating 获取用户作为Sitter的平均分
func GetAverageSitterRating(sitterID string) (float64, error) {
	db := postgres.GetDB()
	var avg float64
	err := db.Model(&model.SitterRating{}).Where("sitter_id = ?", sitterID).Select("COALESCE(AVG(score), 0)").Scan(&avg).Error
	if err != nil {
		logrus.Errorf("get average sitter rating failed, err: %v", err)
		return 0, err
	}
	return avg, nil
}
