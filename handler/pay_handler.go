package handler

import (
	"net/http"

	"demo/util/json"
	"demo/util/open_api"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/payout"
)

// CreateCheckoutSession 创建 Checkout Session
//func CreateCheckoutSession(c *gin.Context) {
//	stripe.Key = "sk_test_51RnEzzQNvpWOmO7YCBM6O9GULyc8XiHbwper8eluVrEnqKFrug36nMMsiNMHWScyT7Qzoizr94JofSt3uHnXrotT008REuqzrr"
//
//	var request struct {
//		Amount     int64  `json:"amount"`
//		Currency   string `json:"currency"`
//		CustomerID string `json:"customerId"`
//		SuccessURL string `json:"successUrl"` // 支付成功后的跳转页面
//		CancelURL  string `json:"cancelUrl"`  // 支付取消后的跳转页面
//	}
//	if err := c.ShouldBindJSON(&request); err != nil {
//		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "Missing required fields")
//		return
//	}
//
//	request.SuccessURL = "https://baidu.com"
//	request.CancelURL = "https://baidu.com"
//	params := &stripe.CheckoutSessionParams{
//		//Customer: stripe.String(request.CustomerID),
//		PaymentMethodTypes: []*string{
//			stripe.String("card"),
//			//stripe.String("google_pay"),
//			stripe.String("link"),
//			// stripe.String("apple_pay"), // 待 Apple Pay 启用
//		},
//		LineItems: []*stripe.CheckoutSessionLineItemParams{
//			{
//				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
//					Currency: stripe.String(request.Currency),
//					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
//						Name: stripe.String("Test Product"),
//					},
//					UnitAmount: stripe.Int64(request.Amount), // 金额（以分为单位）
//				},
//				Quantity: stripe.Int64(1),
//			},
//		},
//		Mode:       stripe.String("payment"),
//		SuccessURL: stripe.String(request.SuccessURL),
//		CancelURL:  stripe.String(request.CancelURL),
//		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
//			CaptureMethod: stripe.String("manual"), // 延后扣款
//		},
//	}
//
//	session, err := session.New(params)
//	if err != nil {
//		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	open_api.OpenApiSuccessResponse(c, map[string]interface{}{
//		"sessionId":       session.ID,
//		"paymentIntentId": session.PaymentIntent.ID, // 用于后续捕获
//	})
//}

// 创建托管支付页面的 Checkout Session
func HandleCreateCheckoutSession(c *gin.Context) {
	stripe.Key = "sk_test_51RnEzzQNvpWOmO7YCBM6O9GULyc8XiHbwper8eluVrEnqKFrug36nMMsiNMHWScyT7Qzoizr94JofSt3uHnXrotT008REuqzrr"
	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(1000), // $10.00
		Currency:      stripe.String("usd"),
		CaptureMethod: stripe.String("manual"), // 关键：延迟扣款
		// 可以添加其他参数如 customer, metadata 等
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, map[string]interface{}{
		"clientSecret":    pi.ClientSecret,
		"paymentIntentId": pi.ID,
	})
}

// 手动扣款
func HandleCapturePayment(c *gin.Context) {
	stripe.Key = "sk_test_51RnEzzQNvpWOmO7YCBM6O9GULyc8XiHbwper8eluVrEnqKFrug36nMMsiNMHWScyT7Qzoizr94JofSt3uHnXrotT008REuqzrr"
	paymentIntentID := c.Query("PaymentIntentId")
	if paymentIntentID == "" {
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "payment_intent_id is required")
		return
	}

	pi, err := paymentintent.Capture(paymentIntentID, nil)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, map[string]interface{}{
		"status": pi.Status,
		//"amountCaptured": pi.AmountCaptured,
		//"captured":       pi.Captured,
	})
}

func handleCreatePayout(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Amount      int64  `json:"amount"`      // 提现金额(分)
		Currency    string `json:"currency"`    // 货币代码(usd,eur等)
		Description string `json:"description"` // 可选描述
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 创建提现参数
	params := &stripe.PayoutParams{
		Amount:      stripe.Int64(request.Amount),
		Currency:    stripe.String(request.Currency),
		Description: stripe.String(request.Description),
		Method:      stripe.String("instant"), // 即时到账(额外费用)
		// 默认是标准到账(2-5个工作日)
	}

	// 创建提现
	p, err := payout.New(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":           p.ID,
		"amount":       p.Amount,
		"currency":     p.Currency,
		"status":       p.Status,
		"arrival_date": p.ArrivalDate,
		"created":      p.Created,
	})
}
