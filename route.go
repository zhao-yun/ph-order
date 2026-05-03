package main

import (
	"time"

	"demo/handler"
	"demo/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// routeInit 路由初始化
func routeInit() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		// 允许的前端域名，多个域名用逗号分隔（生产环境需指定具体域名，不要用*）
		AllowOrigins:     []string{"*"},                                                // 替换为你的前端实际域名
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, // 允许的HTTP方法
		AllowHeaders:     []string{"*"},                                                // 允许的请求头
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,           // 允许前端携带Cookie
		MaxAge:           12 * time.Hour, // 预检请求的缓存时间
	}))

	r.POST("/pay", handler.HandleCreateCheckoutSession)
	r.POST("/pay/capture", handler.HandleCapturePayment)
	// 分页查询订单
	r.GET("/Order/GetOrdersPage", service.GetOrdersPage)
	// 用户分页查询我的订单
	r.GET("/Order/UserGetMyOrdersPage", service.UserGetMyOrdersPage)
	// Sitter分页查询我的订单
	r.GET("/Order/SitterGetMyOrdersPage", service.SitterGetMyOrdersPage)
	// 通过Id查询订单
	r.GET("/Order/GetOrderById", service.GetOrderById)
	// 创建订单
	r.POST("/Order/CreateOrder", service.CreateOrder)
	r.POST("/Order/CreateOrderWithPricing", service.CreateOrderWithPricing)
	// 发起支付
	r.POST("/Order/Pay", service.PayOrder)
	// App端支付确认
	r.POST("/Order/ConfirmPayment", service.ConfirmPayment)
	// Sitter接受/拒绝订单
	r.POST("/Order/SitterHandleInvite", service.SitterHandleOrderInvite)
	// Sitter设置订单建立验证码
	r.POST("/Order/SitterSetCreateCode", service.SitterSetCreateCode)
	// Sitter完成订单
	r.POST("/Order/SitterFinishOrder", service.SitterOrderFinish)
	// Sitter设置订单完成验证码
	r.POST("/Order/SitterSetFinishCode", service.SitterSetFinishCode)
	// 添加小费
	r.POST("/Order/AddTip", service.AddTips)
	// 用户删除订单
	r.POST("/Order/UserDeleteOrder", service.UserDeleteOrder)
	// Sitter删除订单
	r.POST("/Order/SitterDeleteOrder", service.SitterDeleteOrder)
	// 查询联想词
	r.GET("/Order/GetAssociatedWords", service.GetAssociatedWords)

	// 用户更新订单
	r.POST("/Order/UserUpdateOrder", service.UserUpdateOrder)
	// Sitter更新订单
	r.POST("/Order/SitterUpdateOrder", service.SitterUpdateOrder)
	// Sitter确认订单修改
	r.POST("/Order/SitterConfirmModification", service.SitterConfirmModification)
	// 用户确认订单修改
	r.POST("/Order/UserConfirmModification", service.UserConfirmModification)

	// 用户取消订单
	r.POST("/Order/UserCancelOrder", service.UserCancelOrder)

	// 子订单管理
	r.GET("/SubOrder/GetSubOrdersPage", service.GetSubOrdersPage)
	r.GET("/SubOrder/GetSubOrderById", service.GetSubOrderById)
	r.GET("/SubOrder/GetSubOrdersByOrderId", service.GetSubOrdersByOrderId)
	r.GET("/SubOrder/GetTodayOrNearestSubOrder", service.GetTodayOrNearestSubOrder)
	r.POST("/SubOrder/SitterHandleSubOrder", service.SitterHandleSubOrder)
	r.POST("/SubOrder/SitterSetSubOrderStartCode", service.SitterSetSubOrderStartCode)
	r.POST("/SubOrder/SitterSubOrderFinish", service.SitterSubOrderFinish)
	r.POST("/SubOrder/SitterSetSubOrderEndCode", service.SitterSetSubOrderEndCode)

	// 遛狗轨迹管理
	r.POST("/WalkRecord/AppendWalkPath", service.AppendWalkPath)
	r.GET("/WalkRecord/GetWalkRecordBySubOrderId", service.GetWalkRecordBySubOrderId)
	r.POST("/SubOrder/UpdateWalkThumbnail", service.UpdateSubOrderWalkThumbnail)

	// Sitter实时位置管理
	r.GET("/SitterLocation/GetLatest", service.GetLatestSitterLocation)
	r.GET("/SitterLocation/GetBySubOrderId", service.GetSitterLocationsBySubOrderId)

	// 评价管理
	sitterRatingGroup := r.Group("/sitterRating")
	{
		// 获取评价列表（分页）
		sitterRatingGroup.GET("/list", service.GetSitterRatingsPage)

		// 创建评价
		sitterRatingGroup.POST("/create", service.CreateSitterRating)

		// 更新评价
		sitterRatingGroup.POST("/update", service.UpdateSitterRating)

		// 删除评价
		sitterRatingGroup.POST("/delete", service.DeleteSitterRating)

		// 通过ID获取评价
		sitterRatingGroup.GET("/detail", service.GetSitterRatingById)

		// 通过订单ID获取评价
		sitterRatingGroup.GET("/byOrder", service.GetSitterRatingByOrderId)
	}

	ownerRatingGroup := r.Group("/ownerRating")
	{
		// 创建评价
		ownerRatingGroup.POST("/create", service.CreateOwnerRating)
	}

	// 评分查询接口
	r.GET("/rating/average", service.GetUserAverageRatings)

	err := r.Run(":8000")
	if err != nil {
		logrus.Errorf("listening port failed, err: %v", err)
		panic(err)
	}
}
