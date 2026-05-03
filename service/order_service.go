package service

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"demo/dal"
	"demo/model"
	"demo/util/auth"
	"demo/util/common"
	"demo/util/helper"
	"demo/util/open_api"
	"demo/util/redis"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72"
	_ "github.com/stripe/stripe-go/v72"
)

// GetOrdersPage 分页查询订单数据.
func GetOrdersPage(c *gin.Context) {

	// 参数校验.
	params, current, size, err := paramCheck(c)
	if err != nil {
		logrus.Errorf("param check failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 查询数据
	orderPage, err := dal.GetOrderPage(c, current, size, params)
	if err != nil {
		logrus.Errorf("GetTestDataPage failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 处理订单状态.
	handelOrderListState(orderPage.OrderList)

	// 处理订单面试状态
	for _, order := range orderPage.OrderList {
		handelOrderInterviewState(order)
	}

	open_api.OpenApiSuccessResponse(c, orderPage)

}

// UserGetMyOrdersPage 分页查询我的订单数据.
func UserGetMyOrdersPage(c *gin.Context) {

	// 参数校验.
	params, current, size, err := paramCheck(c)
	if err != nil {
		logrus.Errorf("param check failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//ownerId, err := auth.GetUserID(c)
	//if err != nil {
	//	logrus.Errorf("GetUserID failed, err: %v", err)
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//params.OwnerID = ownerId

	userId := c.Query("UserID")
	if userId == "" {
		logrus.Errorf("invalid UserID parameter, UserID = %v", userId)
		err = errors.New("invalid UserID parameter")
		return
	}
	params.OwnerID = userId

	if (params.PetIDList != nil && len(params.PetIDList) > 0) || (params.PetTypeList != nil && len(params.PetTypeList) > 0) {
		orderIdList, err := dal.GetOrderIdByOrderPetList(params)
		if err != nil {
			logrus.Errorf("GetOrderIdByOrderPetList failed, err: %v", err)
			open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		params.IDList = orderIdList
	}

	if params.Keyword != "" {
		orderIdList, err := dal.GetOrderIdListByKeyword(params.Keyword)
		if err != nil {
			logrus.Errorf("GetOrderIdListByKeyword failed, err: %v", err)
			open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		params.KeyIDList = orderIdList
	}

	userDeleted := 0
	params.UserDeleted = &userDeleted

	// 查询数据
	orderPage, err := dal.GetOrderPage(c, current, size, params)
	if err != nil {
		logrus.Errorf("GetOrderPage failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	handelOrderListState(orderPage.OrderList)
	for _, order := range orderPage.OrderList {
		handelOrderInterviewState(order)
	}

	open_api.OpenApiSuccessResponse(c, orderPage)

}

// SitterGetMyOrdersPage Sitter分页查询我的订单数据.
func SitterGetMyOrdersPage(c *gin.Context) {

	// 参数校验.
	params, current, size, err := paramCheck(c)
	if err != nil {
		logrus.Errorf("param check failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//sitterId, err := auth.GetSitterID(c)
	//if err != nil {
	//	logrus.Errorf("GetUserID failed, err: %v", err)
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//params.SitterID = sitterId

	userId := c.Query("UserID")
	if userId == "" {
		logrus.Errorf("invalid UserID parameter, UserID = %v", userId)
		err = errors.New("invalid UserID parameter")
		return
	}
	params.SitterID = userId

	if (params.PetIDList != nil && len(params.PetIDList) > 0) || (params.PetTypeList != nil && len(params.PetTypeList) > 0) {
		orderIdList, err := dal.GetOrderIdByOrderPetList(params)
		if err != nil {
			logrus.Errorf("GetOrderIdByOrderPetList failed, err: %v", err)
			open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		params.IDList = orderIdList
	}

	if params.Keyword != "" {
		orderIdList, err := dal.GetOrderIdListByKeyword(params.Keyword)
		if err != nil {
			logrus.Errorf("GetOrderIdListByKeyword failed, err: %v", err)
			open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		params.KeyIDList = orderIdList
	}

	sitterDeleted := 0
	params.SitterDeleted = &sitterDeleted
	// 查询数据
	testDataPage, err := dal.GetOrderPage(c, current, size, params)
	if err != nil {
		logrus.Errorf("GetTestDataPage failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	handelOrderListState(testDataPage.OrderList)
	for _, order := range testDataPage.OrderList {
		handelOrderInterviewState(order)
	}
	open_api.OpenApiSuccessResponse(c, testDataPage)

}

// handelOrderListState 处理订单列表状态.
func handelOrderListState(orderList []*model.Order) {
	// 处理订单状态.
	for _, order := range orderList {
		handleOrderState(order)
	}
}

// handelOrderInterviewState 处理订单列表面试状态.
func handelOrderInterviewState(order *model.Order) {
	if order.State != model.OrderAccepted {
		return
	}
	// 查询订单面试状态
	interviewRecord, _ := dal.GetInterviewRecordByOrderID(order.ID)
	if interviewRecord == nil || interviewRecord.Status == model.Status_Canceled || interviewRecord.Status == model.Status_Rejected {
		status := model.Status_NotStart
		order.InterviewStatus = &status
	} else {
		order.InterviewStatus = &interviewRecord.Status
	}
}

// handleOrderState 处理订单状态.
func handleOrderState(order *model.Order) {
	if order.InterviewStatus == nil {
		interviewStatus := model.Status_NotStart
		order.InterviewStatus = &interviewStatus
	}
	// 判断用户和Sitter有没有输入验证码.
	if order.State == model.OrderAccepted {
		//order.Code = "3247"
		code, err := redis.Get(redis.RedisPrefixCreateCode + helper.I642S(order.ID))
		if err != nil {
			logrus.Errorf("redis get failed, err: %v", err)
			return
		}
		order.Code = code
	}

	if order.State == model.OrderPreCompleted {
		code, err := redis.Get(redis.RedisPrefixFinishCode + helper.I642S(order.ID))
		if err != nil {
			logrus.Errorf("redis get failed, err: %v", err)
			return
		}
		order.Code = code
	}

	// 判定有没有评价
	if order.State == model.OrderCompleted {
		// 查询Sitter评价
		rating, err := dal.GetSitterRatingByOrderId(order.ID)
		if err != nil {
			logrus.Errorf("GetSitterRatingByOrderId failed, err: %v", err)
			return
		}
		order.SitterRating = rating

		// 查询用户评价
		ownerRating, err := dal.GetOwnerRatingByOrderId(order.ID)
		if err != nil {
			logrus.Errorf("GetOwnerRatingByOrderId failed, err: %v", err)
			return
		}
		order.OwnerRating = ownerRating
	}
}

// paramCheck 参数校验.
func paramCheck(c *gin.Context) (params *model.OrderQueryParams, current int, size int, err error) {
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

	// 处理时间参数（将时间字符串转为时间）
	var orderTime *common.Timestamp
	orderTimeStr := c.Query("OrderTime")
	if orderTimeStr != "" {
		// 如果只是传入日期，需要解析为当天的零点
		t, err := time.Parse("2006-01-02", orderTimeStr)
		if err == nil {
			orderTime = &common.Timestamp{Time: t}
		} else {
			// 兼容带时间的格式
			orderTime, err = common.ParseTimestamp(orderTimeStr)
			if err != nil {
				logrus.Errorf("parseTimestamp failed, err: %v", err)
				return nil, 0, 0, err
			}
		}
	}

	var startTime *common.Timestamp
	var endTime *common.Timestamp
	startTimeStr := c.Query("StartTime")
	endTimeStr := c.Query("EndTime")
	if startTimeStr != "" {
		// 兼容日期格式和时间戳格式
		t, err := time.Parse("2006-01-02", startTimeStr)
		if err == nil {
			startTime = &common.Timestamp{Time: t}
		} else {
			startTime, err = common.ParseTimestamp(startTimeStr)
			if err != nil {
				logrus.Errorf("parseTimestamp failed for StartTime, err: %v", err)
				return nil, 0, 0, err
			}
		}
	}
	if endTimeStr != "" {
		// 兼容日期格式和时间戳格式
		t, err := time.Parse("2006-01-02", endTimeStr)
		if err == nil {
			// 如果是结束日期，通常是指当天的23:59:59或者第二天的00:00:00，这里解析为当天的零点，
			// 在dal层处理时如果判断是EndTime，应该包含当天的整个时间段
			endTime = &common.Timestamp{Time: t}
		} else {
			endTime, err = common.ParseTimestamp(endTimeStr)
			if err != nil {
				logrus.Errorf("parseTimestamp failed for EndTime, err: %v", err)
				return nil, 0, 0, err
			}
		}
	}

	typeListStr := c.QueryArray("TypeList") // 获取字符串数组
	var typeList []int64
	for _, v := range typeListStr {
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			// 处理错误
		}
		typeList = append(typeList, num)
	}

	petTypeListStr := c.QueryArray("PetTypeList") // 获取字符串数组
	var petTypeList []string
	for _, v := range petTypeListStr {
		petTypeList = append(petTypeList, v)
	}

	petIDListStr := c.QueryArray("PetIDList") // 获取字符串数组
	var petIDList []string
	for _, v := range petIDListStr {
		petIDList = append(petIDList, v)
	}

	params = &model.OrderQueryParams{
		ID:          helper.S2I64(c.Query("ID")),
		TypeList:    typeList,
		State:       model.OrderStatus(helper.S2I64(c.Query("State"))),
		SitterID:    c.Query("SitterID"),
		OrderTime:   orderTime,
		StartTime:   startTime,
		EndTime:     endTime,
		PetTypeList: petTypeList,
		PetIDList:   petIDList,
		Keyword:     c.Query("Keyword"),
	}

	return
}

// GetOrderById 通过订单ID查询订单.
func GetOrderById(c *gin.Context) {

	orderId := c.Query("OrderID")

	// 查询数据
	order, err := dal.GetOrderById(helper.S2I64(orderId))
	if err != nil {
		logrus.Errorf("GetOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	handleOrderState(order)

	open_api.OpenApiSuccessResponse(c, order)

}

// CreateOrder 创建订单
func CreateOrder(c *gin.Context) {
	var order *model.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		logrus.Errorf("CreateOrder ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户ID
	ownerID, err := auth.GetUserID(c)
	if err != nil {
		logrus.Errorf("GetUserID failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	order.OwnerID = ownerID

	order.State = model.OrderInitialized

	// 如果是遛狗订单，生成子订单
	if order.Type == int64(model.OrderServiceTypeWalking) {
		// 计算日期范围内的天数
		days := int(order.ToDate.Sub(order.FromDate.Time).Hours()/24) + 1
		subOrders := make([]*model.SubOrder, 0, days)
		
		// 为每一天生成一个子订单
		for i := 0; i < days; i++ {
			date := order.FromDate.Time.AddDate(0, 0, i)
			subOrder := &model.SubOrder{
				OrderID:          order.ID,
				Date:             common.Date{Time: date},
				State:            model.OrderAccepted, // 子订单初始状态为已接受
				SitterHandleAt:   common.Timestamp{Time: time.Now()},
			}
			subOrders = append(subOrders, subOrder)
		}
		
		// 创建订单和子订单
		order, err = dal.CreateOrderWithSubOrders(order, subOrders)
	} else {
		// 普通订单，直接创建
		order, err = dal.CreateOrder(order)
	}
	
	if err != nil {
		logrus.Errorf("CreateOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, order)

}

// 支付授权
func PayOrder(c *gin.Context) {
	var param *model.PayOrderReq

	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("PayOrder ShouldBindJSON failed, err: %v", err)
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

	// 校验订单状态
	if order.State != model.OrderInitialized {
		logrus.Errorf("invalid order state, state = %v", order.State)
		err := errors.New("invalid order state")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 查询订单是否为自己的
	//ownerID, err := auth.GetUserID(c)
	//if err != nil {
	//	logrus.Errorf("GetUserID failed, err: %v", err)
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//if order.OwnerID != ownerID {
	//	logrus.Errorf(" No permission for this order")
	//	err := errors.New("no permission for this order")
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}

	// 支付授权
	// 发起支付
	payParam := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(int64(math.Round(order.TotalPrice * 100))), // 修复浮点数精度
		Currency:      stripe.String("usd"),
		CaptureMethod: stripe.String("manual"), // 关键：延迟扣款
		// 可以添加其他参数如 customer, metadata 等
	}
	res, err := HandleCreateCheckoutSession(c, payParam)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	// 保存 PaymentIntentId 供后续扣款使用
	if piID, ok := res["paymentIntentId"].(string); ok {
		redis.Set("payment_intent:"+helper.I642S(order.ID), piID, 30*24*time.Hour)
	}

	open_api.OpenApiSuccessResponse(c, res)

}

// ConfirmPayment App端支付完成后调用的确认接口
func ConfirmPayment(c *gin.Context) {
	var param struct {
		OrderID int64 `json:"OrderID"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("ConfirmPayment ShouldBindJSON failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 查询订单
	order, err := dal.GetOrderById(param.OrderID)
	if err != nil || order == nil {
		logrus.Errorf("GetOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "order not found")
		return
	}

	// 仅当订单为初始化状态时才变更为已支付
	if order.State != model.OrderInitialized {
		logrus.Errorf("invalid order state for ConfirmPayment, current state = %v", order.State)
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "invalid order state")
		return
	}

	order.State = model.OrderPayed
	err = dal.UpdateOrder(order)
	if err != nil {
		logrus.Errorf("UpdateOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, nil)
}

// SitterHandleOrderInvite Sitter接受/拒绝订单
func SitterHandleOrderInvite(c *gin.Context) {
	var param *model.SitterHandleOrderInviteReq

	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("SitterHandleOrderInviteReq ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if param.State != model.OrderAccepted && param.State != model.OrderRejected {
		logrus.Errorf("invalid state parameter, state = %v", param.State)
		err := errors.New("invalid state parameter")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 查询订单
	order, err := dal.GetOrderById(param.OrderID)
	if err != nil {
		logrus.Errorf("GetOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if param.State == model.OrderAccepted {
		// 接单前必须确保订单已支付
		if order.State != model.OrderPayed {
			logrus.Errorf("cannot accept unpaid order, order.State = %v", order.State)
			err := errors.New("cannot accept unpaid order")
			open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	//sitterId := "ec5d3518-b0d1-7084-9840-4304844857f9"
	//// 订单是否是该用户的
	//if order.SitterID != sitterId {
	//	logrus.Errorf(" No permission for this order")
	//	err := errors.New("no permission for this order")
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}

	order.State = param.State
	order.SitterHandleAt = common.Timestamp{
		Time: time.Now(),
	}
	// 更新订单状态
	err = dal.UpdateOrder(order)
	if err != nil {
		logrus.Errorf("UpdateOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 如果是遛狗订单，为每个子订单生成开始验证码
	if order.Type == int64(model.OrderServiceTypeWalking) {
		subOrders, err := dal.GetSubOrdersByOrderID(order.ID)
		if err == nil {
			for _, subOrder := range subOrders {
				subOrder.SitterHandleAt = order.SitterHandleAt
				if param.State == model.OrderAccepted {
					code := common.GenerateRandomCode(4)
					subOrder.StartCode = code
				}
				dal.UpdateSubOrder(subOrder)
			}
		}
	} else {
		// 接受订单要发送验证码
		if param.State == model.OrderAccepted {
			code := common.GenerateRandomCode(4)
			err = redis.Set(redis.RedisPrefixCreateCode+helper.I642S(order.ID), code, 24*time.Hour)
			if err != nil {
				logrus.Errorf("redis set failed, err: %v", err)
				open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	open_api.OpenApiSuccessResponse(c, order)

}

// SitterSetCreateCode Sitter输入建立验证码
func SitterSetCreateCode(c *gin.Context) {
	var param *model.SitterSetCreateCodeReq
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("SitterSetCreateCodeReq ShouldBindJSON failed, err: %v", err)
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
	if order == nil {
		logrus.Errorf("order not found")
		err := errors.New("order not found")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 如果是遛狗订单，提示使用子订单接口
	if order.Type == int64(model.OrderServiceTypeWalking) {
		logrus.Errorf("walking order should use sub order api")
		err := errors.New("walking order should use sub order api")
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	code, err := redis.Get(redis.RedisPrefixCreateCode + helper.I642S(order.ID))
	if err != nil {
		logrus.Errorf("redis get failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 验证code是否正确
	if code != param.Code {
		logrus.Errorf("invalid code")
		err := errors.New("invalid code")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 判断用户是否验证成功
	//userState, err := redis.Get(redis.RedisPrefixCreateCodeUserState + helper.I642S(order.ID))
	//if err != nil {
	//	logrus.Errorf("redis get failed, err: %v", err)
	//	open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//if userState != "1" {
	//	// 如果用户没有验证成功，将Sitter状态更新为已验证
	//	err = redis.Set(redis.RedisPrefixCreateCodeSitterState+helper.I642S(order.ID), "1", 24*time.Hour)
	//	if err != nil {
	//		logrus.Errorf("redis set failed, err: %v", err)
	//		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
	//		return
	//	}
	//	open_api.OpenApiSuccessResponse(c, nil)
	//	return
	//}
	// 如果用户已经验证成功，将订单状态更新为已建立
	order.State = model.OrderEstablished
	err = dal.UpdateOrder(order)
	if err != nil {
		logrus.Errorf("UpdateOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 生成结束验证码
	code = common.GenerateRandomCode(4)
	err = redis.Set(redis.RedisPrefixFinishCode+helper.I642S(order.ID), code, 24*time.Hour)
	if err != nil {
		logrus.Errorf("redis set failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	open_api.OpenApiSuccessResponse(c, nil)
}

// SitterOrderFinish Sitter完成订单
func SitterOrderFinish(c *gin.Context) {
	var param *model.SitterFinishOrderReq
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("SitterOrderFinishReq ShouldBindJSON failed, err: %v", err)
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
	if order == nil {
		logrus.Errorf("order not found")
		err := errors.New("order not found")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 如果是遛狗订单，提示使用子订单接口
	if order.Type == int64(model.OrderServiceTypeWalking) {
		logrus.Errorf("walking order should use sub order api")
		err := errors.New("walking order should use sub order api")
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 订单状态是否是已建立
	if order.State != model.OrderEstablished {
		logrus.Errorf("invalid order state")
		err := errors.New("invalid order state")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	order.State = model.OrderPreCompleted
	err = dal.UpdateOrder(order)

	// 生成结束验证码
	code := common.GenerateRandomCode(4)
	err = redis.Set(redis.RedisPrefixFinishCode+helper.I642S(order.ID), code, 24*time.Hour)
	if err != nil {
		logrus.Errorf("redis set failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	open_api.OpenApiSuccessResponse(c, nil)

}

// SitterSetFinishCode Sitter输入结束验证码
func SitterSetFinishCode(c *gin.Context) {
	var param *model.SitterSetFinishCodeReq
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("SitterSetFinishCodeReq ShouldBindJSON failed, err: %v", err)
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
	if order == nil {
		logrus.Errorf("order not found")
		err := errors.New("order not found")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 如果是遛狗订单，提示使用子订单接口
	if order.Type == int64(model.OrderServiceTypeWalking) {
		logrus.Errorf("walking order should use sub order api")
		err := errors.New("walking order should use sub order api")
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	code, err := redis.Get(redis.RedisPrefixFinishCode + helper.I642S(order.ID))
	if err != nil {
		logrus.Errorf("redis get failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 验证code是否正确
	if code != param.Code {
		logrus.Errorf("invalid code")
		err := errors.New("invalid code")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 修改状态
	order.State = model.OrderCompleted
	err = dal.UpdateOrder(order)
	if err != nil {
		logrus.Errorf("UpdateOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 订单完成，自动进行 Stripe 扣款
	piID, _ := redis.Get("payment_intent:" + helper.I642S(order.ID))
	if piID != "" && piID != "0" {
		stripe.Key = "sk_test_51RnEzzQNvpWOmO7YCBM6O9GULyc8XiHbwper8eluVrEnqKFrug36nMMsiNMHWScyT7Qzoizr94JofSt3uHnXrotT008REuqzrr"
		captureErr := stripe.GetBackend(stripe.APIBackend).Call(
			http.MethodPost,
			"/v1/payment_intents/"+piID+"/capture",
			stripe.Key,
			nil,
			&stripe.PaymentIntent{},
		)
		if captureErr != nil {
			logrus.Errorf("stripe capture failed, piID: %s, err: %v", piID, captureErr)
		}
	}

	open_api.OpenApiSuccessResponse(c, nil)

}

// UserCancelOrder 用户取消订单
func UserCancelOrder(c *gin.Context) {
	var param *model.UserCancelOrderReq
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("UserCancelOrderReq ShouldBindJSON failed, err: %v", err)
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
	if order == nil {
		logrus.Errorf("order not found")
		err := errors.New("order not found")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 订单状态是否是已建立
	if order.State != model.OrderAccepted && order.State != model.OrderInitialized {
		logrus.Errorf("invalid order state")
		err := errors.New("invalid order state")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	order.State = model.OrderCancel
	order.CancelAt = common.Timestamp{
		Time: time.Now(),
	}
	order.CancelReason = param.Reason

	err = dal.UpdateOrder(order)
	if err != nil {
		logrus.Errorf("UpdateOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	open_api.OpenApiSuccessResponse(c, nil)
}

// AddTips 添加小费
func AddTips(c *gin.Context) {
	var tipsParam *model.AddTipsReq
	if err := c.ShouldBindJSON(&tipsParam); err != nil {
		logrus.Errorf("AddTipsReq ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 查询订单
	order, err := dal.GetOrderById(tipsParam.OrderID)
	if err != nil {
		logrus.Errorf("GetOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if order == nil {
		logrus.Errorf("order not found")
		err := errors.New("order not found")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 添加小费
	order.TipsPrice = tipsParam.Tips
	err = dal.UpdateOrder(order)
	if err != nil {
		logrus.Errorf("UpdateOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 发起支付
	param := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(int64(math.Round(order.TipsPrice * 100))), // 修复浮点数精度
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
}

// UserDeleteOrder 用户删除订单
func UserDeleteOrder(c *gin.Context) {
	var param *model.DeleteOrderReq
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("UserDeleteOrderReq ShouldBindJSON failed, err: %v", err)
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
	if order == nil {
		logrus.Errorf("order not found")
		err := errors.New("order not found")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 用户删除订单
	order.UserDeleted = 1
	err = dal.UpdateOrder(order)
	if err != nil {
		logrus.Errorf("UpdateOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	open_api.OpenApiSuccessResponse(c, nil)
}

// SitterDeleteOrder Sitter删除订单
func SitterDeleteOrder(c *gin.Context) {
	var param *model.DeleteOrderReq
	if err := c.ShouldBindJSON(&param); err != nil {
		logrus.Errorf("SitterDeleteOrderReq ShouldBindJSON failed, err: %v", err)
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
	if order == nil {
		logrus.Errorf("order not found")
		err := errors.New("order not found")
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Sitter删除订单
	order.SitterDeleted = 1
	err = dal.UpdateOrder(order)
	if err != nil {
		logrus.Errorf("UpdateOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	open_api.OpenApiSuccessResponse(c, nil)
}

// GetAssociatedWords 查询联想词
func GetAssociatedWords(c *gin.Context) {
	userId, err := auth.GetUserID(c)
	if err != nil {
		logrus.Errorf("auth failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	keyword := c.Query("Keyword")
	if keyword == "" {
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "keyword is empty")
		return
	}

	var res []string

	// 查询订单相关用户名字
	userNameList, err := dal.GetUsernameByKeyword(keyword, userId)
	if err == nil {
		res = append(res, userNameList...)
	}

	// 查询宠物名字
	petNameList, err := dal.GetPetNamesByKeyword(keyword, userId)
	if err == nil {
		res = append(res, petNameList...)
	}

	open_api.OpenApiSuccessResponse(c, res)
}
