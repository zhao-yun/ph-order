package service

import (
	"errors"
	"net/http"
	"strconv"

	"demo/dal"
	"demo/model"
	"demo/util/auth"
	"demo/util/helper"
	"demo/util/open_api"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetSitterRatingsPage 分页查询评价数据
func GetSitterRatingsPage(c *gin.Context) {
	// 参数校验
	current, size, err := paramRatingCheck(c)
	if err != nil {
		logrus.Errorf("param check failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	params := model.SitterRatingQueryParams{
		ID:       helper.S2I64(c.Query("ID")),
		OrderID:  helper.S2I64(c.Query("OrderID")),
		UserID:   c.Query("UserID"),
		SitterID: c.Query("SitterID"),
	}

	// 查询数据
	ratingPage, err := dal.GetSitterRatingPage(c, current, size, params)
	if err != nil {
		logrus.Errorf("GetSitterRatingsPage failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, ratingPage)
}

// CreateSitterRating 创建评价
func CreateSitterRating(c *gin.Context) {
	var rating *model.SitterRating

	if err := c.ShouldBindJSON(&rating); err != nil {
		logrus.Errorf("CreateSitterRating ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询订单
	order, err := dal.GetOrderById(rating.OrderID)
	if err != nil {
		logrus.Errorf("GetOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 订单是否是该用户的
	userId, err := auth.GetUserID(c)
	if err != nil {
		logrus.Errorf("auth failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//if order.OwnerID != userId {
	//	logrus.Errorf(" No permission for this order %v", err)
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	// 查询评价是否存在
	sitterRating, _ := dal.GetSitterRatingByOrderId(rating.OrderID)
	if sitterRating != nil {
		rating.ID = sitterRating.ID
		err := dal.UpdateSitterRating(rating)
		if err != nil {
			logrus.Errorf("UpdateSitterRating failed, err: %v", err)
			open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		open_api.OpenApiSuccessResponse(c, rating)
		return
	}

	rating.UserID = userId
	rating.SitterID = order.SitterID

	// 创建评价
	rating, err = dal.CreateSitterRating(rating)
	if err != nil {
		logrus.Errorf("CreateSitterRating failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, rating)
}

// UpdateSitterRating 更新评价
func UpdateSitterRating(c *gin.Context) {
	var rating *model.SitterRating

	if err := c.ShouldBindJSON(&rating); err != nil {
		logrus.Errorf("UpdateSitterRating ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新评价
	err := dal.UpdateSitterRating(rating)
	if err != nil {
		logrus.Errorf("UpdateSitterRating failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, nil)
}

// DeleteSitterRating 删除评价
func DeleteSitterRating(c *gin.Context) {
	idStr := c.Query("ID")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		logrus.Errorf("invalid id parameter, id = %v", idStr)
		err = errors.New("invalid id parameter")
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 删除评价
	err = dal.DeleteSitterRating(id)
	if err != nil {
		logrus.Errorf("DeleteSitterRating failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, nil)
}

// GetSitterRatingById 通过ID查询评价
func GetSitterRatingById(c *gin.Context) {
	idStr := c.Query("ID")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		logrus.Errorf("invalid id parameter, id = %v", idStr)
		err = errors.New("invalid id parameter")
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 查询评价
	rating, err := dal.GetSitterRatingById(id)
	if err != nil {
		logrus.Errorf("GetSitterRatingById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, rating)
}

// GetSitterRatingByOrderId 通过订单ID查询评价
func GetSitterRatingByOrderId(c *gin.Context) {
	orderIdStr := c.Query("OrderID")
	orderId := helper.S2I64(orderIdStr)

	// 查询评价
	rating, err := dal.GetSitterRatingByOrderId(orderId)
	if err != nil {
		logrus.Errorf("GetSitterRatingByOrderId failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, rating)
}

// CreateOwnerRating 创建评价
func CreateOwnerRating(c *gin.Context) {
	var rating *model.OwnerRating
	if err := c.ShouldBindJSON(&rating); err != nil {
		logrus.Errorf("CreateOwnerRating ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 查询订单
	order, err := dal.GetOrderById(rating.OrderID)
	if err != nil {
		logrus.Errorf("GetOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 订单是否是该用户的
	//sitterID, err := auth.GetSitterID(c)
	//if err != nil {
	//	logrus.Errorf("auth failed, err: %v", err)
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//if order.SitterID != sitterID {
	//	logrus.Errorf("No permission for this order")
	//	err := errors.New("no permission for this order")
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	// 查询评价是否存在
	ownerRating, _ := dal.GetOwnerRatingByOrderId(rating.OrderID)
	if ownerRating != nil {
		rating.ID = ownerRating.ID
		err := dal.UpdateOwnerRating(rating)
		if err != nil {
			logrus.Errorf("UpdateOwnerRating failed, err: %v", err)
			open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		open_api.OpenApiSuccessResponse(c, rating)
		return
	}
	rating.SitterID = order.SitterID
	rating.OwnerID = order.OwnerID
	// 创建评价
	rating, err = dal.CreateOwnerRating(rating)
	if err != nil {
		logrus.Errorf("CreateOwnerRating failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	open_api.OpenApiSuccessResponse(c, rating)
}

// GetUserAverageRatings 获取用户作为 Sitter 和 Owner 的平均评分
func GetUserAverageRatings(c *gin.Context) {
	userId := c.Query("UserID")
	if userId == "" {
		logrus.Errorf("UserID is required")
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "UserID is required")
		return
	}

	// 1. 获取作为 Sitter 的平均分
	sitterAvg, err := dal.GetAverageSitterRating(userId)
	if err != nil {
		logrus.Errorf("GetAverageSitterRating failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 2. 获取作为 Owner 的平均分
	ownerAvg, err := dal.GetAverageOwnerRating(userId)
	if err != nil {
		logrus.Errorf("GetAverageOwnerRating failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, map[string]interface{}{
		"UserID":             userId,
		"SitterAverageScore": sitterAvg,
		"OwnerAverageScore":  ownerAvg,
	})
}

// paramCheck 参数校验.
func paramRatingCheck(c *gin.Context) (current int, size int, err error) {
	currentStr := c.Query("Current")
	sizeStr := c.Query("Size")

	// 分页默认值
	current = 1
	size = 10

	// 参数校验
	if currentStr != "" {
		current, err = strconv.Atoi(currentStr)
		if err != nil || current < 1 {
			logrus.Errorf("invalid current parameter, current = %v", current)
			err = errors.New("invalid current parameter")
			return
		}
	}

	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		// 分页最大数量限制.
		if err != nil || size < 1 || size > 100 {
			logrus.Errorf("invalid size parameter, size = %v", size)
			err = errors.New("invalid size parameter")
			return
		}
	}
	return
}
