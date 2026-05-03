package service

import (
	"fmt"
	"net/http"
	"time"

	"demo/dal"
	"demo/model"
	"demo/util/auth"
	"demo/util/common"
	"demo/util/helper"
	open_api "demo/util/open_api"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreateOrderWithPricing(c *gin.Context) {
	var req *model.CreateOrderWithPricingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("CreateOrderWithPricing ShouldBindJSON failed, err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.UserID == "" {
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "userId is required")
		return
	}
	if req.OwnerName == "" {
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "ownerName is required")
		return
	}
	if req.SitterName == "" {
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "sitterName is required")
		return
	}
	if len(req.PetList) == 0 {
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "PetList is required")
		return
	}
	for _, pet := range req.PetList {
		if pet == nil {
			open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "PetList contains nil item")
			return
		}
		if pet.PetID == "" {
			open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "PetList[].PetID is required")
			return
		}
		if pet.PetType == "" {
			open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "PetList[].PetType is required")
			return
		}
		if pet.PetShape == 0 {
			open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "PetList[].PetShape is required")
			return
		}
		if pet.PetName == "" {
			open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "PetList[].PetName is required")
			return
		}
		if pet.Breed == "" {
			open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "PetList[].Breed is required")
			return
		}
	}

	tokenUserID, err := auth.GetUserIDFromToken(c)
	if err != nil {
		logrus.Warnf("GetUserIDFromToken failed, err: %v", err)
	}
	ownerID := req.UserID
	if tokenUserID != "" {
		ownerID = tokenUserID
	}

	pricingResult, err := CalculateOrderPricing(req)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	orderNumber := fmt.Sprintf("ORD%s%s", time.Now().Format("20060102150405"), helper.RandString(8))

	order := &model.Order{
		OwnerID:            ownerID,
		SitterID:           req.SitterID,
		Type:               int64(pricingResult.ServiceType),
		FromDate:           req.FromDate,
		ToDate:             req.ToDate,
		PetList:            req.PetList,
		SubTotalPrice:      float64(pricingResult.PetsSubtotalCents) / 100.0,
		TotalPrice:         float64(pricingResult.TotalCents) / 100.0,
		State:              model.OrderInitialized,
		OwnerName:          req.OwnerName,
		SitterName:         req.SitterName,
		Contact:            req.Contact,
		AlternativeContact: req.AlternativeContact,
		Note:               req.Note,
		OrderNumber:        orderNumber,
	}

	for i, pet := range order.PetList {
		if pet == nil {
			continue
		}
		if v, ok := pricingResult.PetTotalsCents[i]; ok && v != nil {
			pet.PetPrice = float64(*v) / 100.0
		}
	}

	var created *model.Order
	if pricingResult.ServiceType == model.OrderServiceTypeWalking {
		subOrders := generateSubOrders(order)
		created, err = dal.CreateOrderWithSubOrders(order, subOrders)
	} else {
		created, err = dal.CreateOrder(order)
	}

	if err != nil {
		logrus.Errorf("CreateOrderWithPricing failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, created)
}

func generateSubOrders(order *model.Order) []*model.SubOrder {
	subOrders := []*model.SubOrder{}
	
	startDate := order.FromDate
	endDate := order.ToDate
	
	currentTime := startDate.Time
	for !currentTime.After(endDate.Time) {
		subOrder := &model.SubOrder{
			Date:  common.Date{Time: currentTime},
			State: model.OrderAccepted,
		}
		subOrders = append(subOrders, subOrder)
		
		currentTime = currentTime.AddDate(0, 0, 1)
	}
	
	return subOrders
}
