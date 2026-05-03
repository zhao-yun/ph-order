package service

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"demo/dal"
	"demo/model"
	"demo/util/common"
	"demo/util/helper"
	open_api "demo/util/open_api"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetSubOrdersPage(c *gin.Context) {
	params, current, size, err := subOrderParamCheck(c)
	if err != nil {
		logrus.Errorf("param check failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	subOrderPage, err := dal.GetSubOrderPage(current, size, params)
	if err != nil {
		logrus.Errorf("GetSubOrderPage failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, subOrderPage)
}

func GetSubOrderById(c *gin.Context) {
	subOrderId := c.Query("SubOrderID")
	subOrder, err := dal.GetSubOrderById(helper.S2I64(subOrderId))
	if err != nil {
		logrus.Errorf("GetSubOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	open_api.OpenApiSuccessResponse(c, subOrder)
}

func GetSubOrdersByOrderId(c *gin.Context) {
	orderId := c.Query("OrderID")
	subOrderList, err := dal.GetSubOrdersByOrderID(helper.S2I64(orderId))
	if err != nil {
		logrus.Errorf("GetSubOrdersByOrderID failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	open_api.OpenApiSuccessResponse(c, subOrderList)
}

func SitterHandleSubOrder(c *gin.Context) {
	open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "sub orders are initialized as accepted, use sub order code APIs directly")
}

func SitterSetSubOrderStartCode(c *gin.Context) {
	var param *model.SitterSetSubOrderStartCodeReq
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("SitterSetSubOrderStartCodeReq ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subOrder, err := dal.GetSubOrderById(param.SubOrderID)
	if err != nil {
		logrus.Errorf("GetSubOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if subOrder == nil {
		logrus.Errorf("sub order not found")
		err := errors.New("sub order not found")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if subOrder.State != model.OrderAccepted {
		logrus.Errorf("invalid sub order state")
		err := errors.New("invalid sub order state")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if subOrder.StartCode != param.Code {
		logrus.Errorf("invalid code")
		err := errors.New("invalid code")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	subOrder.State = model.OrderEstablished
	endCode := common.GenerateRandomCode(4)
	subOrder.EndCode = endCode
	err = dal.UpdateSubOrder(subOrder)
	if err != nil {
		logrus.Errorf("UpdateSubOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, nil)
}

func SitterSubOrderFinish(c *gin.Context) {
	var param *model.SitterFinishOrderReq
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("SitterSubOrderFinishReq ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subOrder, err := dal.GetSubOrderById(param.OrderID)
	if err != nil {
		logrus.Errorf("GetSubOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if subOrder == nil {
		logrus.Errorf("sub order not found")
		err := errors.New("sub order not found")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if subOrder.State != model.OrderEstablished {
		logrus.Errorf("invalid sub order state")
		err := errors.New("invalid sub order state")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	subOrder.State = model.OrderPreCompleted
	subOrder.SitterFinishAt = common.Timestamp{
		Time: time.Now(),
	}
	err = dal.UpdateSubOrder(subOrder)
	if err != nil {
		logrus.Errorf("UpdateSubOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, nil)

}

func SitterSetSubOrderEndCode(c *gin.Context) {
	var param *model.SitterSetSubOrderEndCodeReq
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("SitterSetSubOrderEndCodeReq ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subOrder, err := dal.GetSubOrderById(param.SubOrderID)
	if err != nil {
		logrus.Errorf("GetSubOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if subOrder == nil {
		logrus.Errorf("sub order not found")
		err := errors.New("sub order not found")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if subOrder.State != model.OrderPreCompleted {
		logrus.Errorf("invalid sub order state")
		err := errors.New("invalid sub order state")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if subOrder.EndCode != param.Code {
		logrus.Errorf("invalid code")
		err := errors.New("invalid code")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	subOrder.State = model.OrderCompleted
	err = dal.UpdateSubOrder(subOrder)
	if err != nil {
		logrus.Errorf("UpdateSubOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, nil)
}

func subOrderParamCheck(c *gin.Context) (params *model.SubOrderQueryParams, current int, size int, err error) {
	currentStr := c.Query("Current")
	sizeStr := c.Query("Size")

	current = 1
	size = 10

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
		if err != nil || size < 1 || size > 100 {
			logrus.Errorf("invalid size parameter, size = %v", size)
			err = errors.New("invalid size parameter")
			return
		}
	}

	var startTime *common.Timestamp
	var endTime *common.Timestamp
	startTimeStr := c.Query("StartTime")
	endTimeStr := c.Query("EndTime")
	if startTimeStr != "" {
		startTime, err = common.ParseTimestamp(startTimeStr)
		if err != nil {
			logrus.Errorf("parseTimestamp failed, err: %v", err)
			return
		}
	}
	if endTimeStr != "" {
		endTime, err = common.ParseTimestamp(endTimeStr)
		if err != nil {
			logrus.Errorf("parseTimestamp failed, err: %v", err)
			return
		}
	}

	params = &model.SubOrderQueryParams{
		ID:        helper.S2I64(c.Query("ID")),
		OrderID:   helper.S2I64(c.Query("OrderID")),
		State:     model.OrderStatus(helper.S2I64(c.Query("State"))),
		StartTime: startTime,
		EndTime:   endTime,
	}

	return
}

func GetTodayOrNearestSubOrder(c *gin.Context) {
	orderID := c.Query("OrderID")
	if orderID == "" {
		logrus.Errorf("OrderID is required")
		err := errors.New("OrderID is required")
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	todayStr := c.Query("Date")
	var today common.Date
	if todayStr == "" {
		today = common.Date{Time: time.Now()}
	} else {
		err := today.UnmarshalJSON([]byte(`"` + todayStr + `"`))
		if err != nil {
			logrus.Errorf("parse date failed, err: %v", err)
			open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "invalid date format")
			return
		}
	}

	subOrder, dayNumber, err := dal.GetTodayOrNearestSubOrder(helper.S2I64(orderID), today)
	if err != nil {
		logrus.Errorf("GetTodayOrNearestSubOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if subOrder == nil {
		open_api.OpenApiSuccessResponse(c, nil)
		return
	}

	type Result struct {
		SubOrder  *model.SubOrder `json:"SubOrder"`
		DayNumber int             `json:"DayNumber"`
	}

	result := &Result{
		SubOrder:  subOrder,
		DayNumber: dayNumber,
	}

	open_api.OpenApiSuccessResponse(c, result)
}
