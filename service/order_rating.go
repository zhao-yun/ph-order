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

	userID, err := auth.GetUserID(c)
	if err != nil {
		logrus.Errorf("auth failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	sitterID, err := auth.GetSitterID(c)
	if err != nil {
		logrus.Errorf("auth failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	params := model.SitterRatingQueryParams{
		ID:       helper.S2I64(c.Query("ID")),
		OrderID:  helper.S2I64(c.Query("OrderID")),
		UserID:   userID,
		SitterID: sitterID,
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
	if order.OwnerID != userId {
		logrus.Errorf("no permission for this order")
		open_api.OpenApiErrorResponse(c, http.StatusForbidden, "no permission for this order")
		return
	}
	// 查询评价是否存在
	sitterRating, _ := dal.GetSitterRatingByOrderId(rating.OrderID)
	if sitterRating != nil {
		rating.ID = sitterRating.ID
		rating.OrderID = sitterRating.OrderID
		rating.UserID = userId
		rating.SitterID = order.SitterID
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

	if rating.ID <= 0 {
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "ID is required")
		return
	}

	existingRating, err := dal.GetSitterRatingById(int(rating.ID))
	if err != nil {
		logrus.Errorf("GetSitterRatingById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	order, err := dal.GetOrderById(existingRating.OrderID)
	if err != nil {
		logrus.Errorf("GetOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := ensureOrderOwner(c, order); err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	rating.OrderID = existingRating.OrderID
	rating.UserID = order.OwnerID
	rating.SitterID = order.SitterID

	// 更新评价
	err = dal.UpdateSitterRating(rating)
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

	rating, err := dal.GetSitterRatingById(id)
	if err != nil {
		logrus.Errorf("GetSitterRatingById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	order, err := dal.GetOrderById(rating.OrderID)
	if err != nil {
		logrus.Errorf("GetOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := ensureOrderOwner(c, order); err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusForbidden, err.Error())
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

	order, err := dal.GetOrderById(rating.OrderID)
	if err != nil {
		logrus.Errorf("GetOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := ensureOrderParticipant(c, order); err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, rating)
}

// GetSitterRatingByOrderId 通过订单ID查询评价
func GetSitterRatingByOrderId(c *gin.Context) {
	orderIdStr := c.Query("OrderID")
	orderId := helper.S2I64(orderIdStr)

	order, err := dal.GetOrderById(orderId)
	if err != nil {
		logrus.Errorf("GetOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := ensureOrderParticipant(c, order); err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

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
	if err := ensureOrderSitter(c, order); err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}
	// 查询评价是否存在
	ownerRating, _ := dal.GetOwnerRatingByOrderId(rating.OrderID)
	if ownerRating != nil {
		rating.ID = ownerRating.ID
		rating.OrderID = ownerRating.OrderID
		rating.OwnerID = order.OwnerID
		rating.SitterID = order.SitterID
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

	currentUserID, err := auth.GetUserID(c)
	if err != nil {
		logrus.Errorf("auth failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	currentSitterID, err := auth.GetSitterID(c)
	if err != nil {
		logrus.Errorf("auth failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if userId != currentUserID && userId != currentSitterID {
		logrus.Errorf("no permission to view average rating")
		open_api.OpenApiErrorResponse(c, http.StatusForbidden, "no permission to view average rating")
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

func ensureOrderOwner(c *gin.Context, order *model.Order) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return err
	}
	if order == nil || order.OwnerID != userID {
		return errors.New("no permission for this order")
	}
	return nil
}

func ensureOrderSitter(c *gin.Context, order *model.Order) error {
	sitterID, err := auth.GetSitterID(c)
	if err != nil {
		return err
	}
	if order == nil || order.SitterID != sitterID {
		return errors.New("no permission for this order")
	}
	return nil
}

func ensureOrderParticipant(c *gin.Context, order *model.Order) error {
	if order == nil {
		return errors.New("order not found")
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		return err
	}
	if order.OwnerID == userID {
		return nil
	}

	sitterID, err := auth.GetSitterID(c)
	if err != nil {
		return err
	}
	if order.SitterID == sitterID {
		return nil
	}

	return errors.New("no permission for this order")
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
