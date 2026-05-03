package service

import (
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

func HandleCreateCheckoutSession(c *gin.Context, params *stripe.PaymentIntentParams) (map[string]interface{}, error) {
	stripe.Key = "sk_test_51RnEzzQNvpWOmO7YCBM6O9GULyc8XiHbwper8eluVrEnqKFrug36nMMsiNMHWScyT7Qzoizr94JofSt3uHnXrotT008REuqzrr"

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"clientSecret":    pi.ClientSecret,
		"paymentIntentId": pi.ID,
	}, nil
}
