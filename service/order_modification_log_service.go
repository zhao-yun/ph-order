package service

import (
	"errors"
	"net/http"

	"demo/dal"
	"demo/model"
	"demo/util/json"
	"demo/util/open_api"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72"
)

// UserUpdateOrder 创建订单修改记录
func UserUpdateOrder(c *gin.Context) {
	var param *model.UserUpdateOrderReq

	// 绑定 JSON 数据到模型
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("CreateOrderModificationLog ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询订单
	order, err := dal.GetOrderById(param.OrderID)
	if err != nil {
		logrus.Errorf("GetOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//ownerId, err := auth.GetUserID(c)
	//if err != nil {
	//	logrus.Errorf("auth failed, err: %v", err)
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}

	// 订单是否是该用户的
	//if order.OwnerID != ownerId {
	//	logrus.Errorf(" No permission for this order %v", err)
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}

	log := &model.OrderModificationLog{
		OrderID:         param.OrderID,
		OwnerID:         order.OwnerID,
		SitterID:        order.SitterID,
		PreviousDate:    order.ToDate,
		NewDate:         param.ToDate,
		PreviousPetList: json.ToJSON(order.PetList),
		NewPetList:      json.ToJSON(param.PetList),
		State:           model.OrderModificationInitialized,
		Type:            model.OrderModificationTypeUser,
	}

	// 创建订单修改记录
	createdLog, err := dal.CreateOrderModificationLog(*log)
	if err != nil {
		logrus.Errorf("CreateOrderModificationLog failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 如果订单价格变多
	newPrice := 100.00
	if newPrice > 0.0 {
		param := &stripe.PaymentIntentParams{
			Amount:        stripe.Int64(int64(newPrice * 100)), // $10.00
			Currency:      stripe.String("usd"),
			CaptureMethod: stripe.String("manual"), // 关键：延迟扣款
			// 可以添加其他参数如 customer, metadata 等
		}
		res, err := HandleCreateCheckoutSession(c, param)
		if err != nil {
			open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		open_api.OpenApiSuccessResponse(c, res)
		return
	}

	open_api.OpenApiSuccessResponse(c, createdLog)
}

func SitterUpdateOrder(c *gin.Context) {
	var param *model.SitterUpdateOrderReq
	// 绑定 JSON 数据到模型
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("CreateOrderModificationLog ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 查询订单
	order, err := dal.GetOrderById(param.OrderID)
	if err != nil {
		logrus.Errorf("GetOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//sitterId, err := auth.GetSitterID(c)
	//if err != nil {
	//	logrus.Errorf("auth failed, err: %v", err)
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	// 订单是否是该用户的
	//if order.SitterID != sitterId {
	//	logrus.Errorf(" No permission for this order %v", err)
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}

	newSubTotal := 0.0
	if len(param.PetList) > 0 {
		for _, p := range param.PetList {
			if p != nil {
				newSubTotal += p.PetPrice
			}
		}
	} else {
		newSubTotal = order.SubTotalPrice
	}
	
	newTotalPrice := newSubTotal + order.ServiceFee + order.Taxes

	log := &model.OrderModificationLog{
		OrderID:         param.OrderID,
		OwnerID:         order.OwnerID,
		SitterID:        order.SitterID,
		PreviousDate:    order.ToDate,
		NewDate:         param.ToDate,
		PreviousPetList: json.ToJSON(order.PetList),
		NewPetList:      json.ToJSON(param.PetList),
		PreviousPrice:   order.TotalPrice,
		NewPrice:        newTotalPrice,
		State:           model.OrderModificationInitialized,
		Type:            model.OrderModificationTypeSitter,
	}
	// 创建订单修改记录
	createdLog, err := dal.CreateOrderModificationLog(*log)
	if err != nil {
		logrus.Errorf("CreateOrderModificationLog failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	open_api.OpenApiSuccessResponse(c, map[string]interface{}{
		"log": createdLog,
		"priceChanges": map[string]interface{}{
			"previousPrice": order.TotalPrice,
			"newPrice":      newTotalPrice,
			"difference":    newTotalPrice - order.TotalPrice,
		},
	})
}

// SitterConfirmModification Sitter确认订单修改
func SitterConfirmModification(c *gin.Context) {
	var param *model.SitterConfirmModificationReq

	// 绑定 JSON 数据到模型
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("CreateOrderModificationLog ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询订单正在修改的记录
	log, err := dal.GetIngOrderModificationLogByOrderID(param.OrderID)
	if err != nil {
		logrus.Errorf("GetIngOrderModificationLogByOrderID failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 判断是否存在
	if log == nil {
		logrus.Errorf("No order modification log found for order ID %d", param.OrderID)
		err := errors.New("no order modification log found")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//// 判断是否是用户发起的修改记录
	//if log.Type != model.OrderModificationTypeUser {
	//	logrus.Errorf("No permission for this order modification log")
	//	err := errors.New("no permission for this order modification log")
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	// 修改订单
	log.State = param.State
	err = dal.UpdateOrderModificationLog(log)
	if err != nil {
		logrus.Errorf("UpdateOrderModificationLog failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	open_api.OpenApiSuccessResponse(c, log)

}

// UserConfirmModification 用户确认订单修改
func UserConfirmModification(c *gin.Context) {
	var param *model.SitterConfirmModificationReq

	// 绑定 JSON 数据到模型
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("CreateOrderModificationLog ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询订单正在修改的记录
	log, err := dal.GetIngOrderModificationLogByOrderID(param.OrderID)
	if err != nil {
		logrus.Errorf("GetIngOrderModificationLogByOrderID failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 判断是否存在
	if log == nil {
		logrus.Errorf("No order modification log found for order ID %d", param.OrderID)
		err := errors.New("no order modification log found")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 判断是否是用户发起的修改记录
	//if log.Type != model.OrderModificationTypeUser {
	//	logrus.Errorf("No permission for this order modification log")
	//	err := errors.New("no permission for this order modification log")
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	// 修改订单
	log.State = param.State
	err = dal.UpdateOrderModificationLog(log)
	if err != nil {
		logrus.Errorf("UpdateOrderModificationLog failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	open_api.OpenApiSuccessResponse(c, log)

}
