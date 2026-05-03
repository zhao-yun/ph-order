package service

import (
	"net/http"
	"strconv"

	"demo/dal"
	"demo/model"
	"demo/util/helper"
	open_api "demo/util/open_api"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UpdateSitterLocation 已废弃，合并到了 WalkRecord/AppendWalkPath 中
// func UpdateSitterLocation(c *gin.Context) {
// 	...
// }

func GetLatestSitterLocation(c *gin.Context) {
	params := &model.SitterLocationQueryParams{
		SitterID:   c.Query("SitterID"),
		OrderID:    helper.S2I64(c.Query("OrderID")),
		SubOrderID: helper.S2I64(c.Query("SubOrderID")),
	}

	if params.SitterID == "" && params.OrderID == 0 && params.SubOrderID == 0 {
		logrus.Errorf("at least one of SitterID, OrderID, SubOrderID is required")
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "at least one of SitterID, OrderID, SubOrderID is required")
		return
	}

	location, err := dal.GetLatestSitterLocation(params)
	if err != nil {
		logrus.Errorf("GetLatestSitterLocation failed, err: %v", err)
		open_api.OpenApiSuccessResponse(c, nil)
		return
	}

	response := &model.SitterLocationResponse{
		ID:         location.ID,
		SitterID:   location.SitterID,
		OrderID:    location.OrderID,
		SubOrderID: location.SubOrderID,
		Lat:        location.Lat,
		Lng:        location.Lng,
		Timestamp:  location.Timestamp,
		CreatedAt:  location.CreatedAt.String(),
	}

	open_api.OpenApiSuccessResponse(c, response)
}

func GetSitterLocationsBySubOrderId(c *gin.Context) {
	subOrderId := c.Query("SubOrderID")
	if subOrderId == "" {
		logrus.Errorf("SubOrderID is required")
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "SubOrderID is required")
		return
	}

	limitStr := c.Query("Limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	locations, err := dal.GetSitterLocationsBySubOrderID(helper.S2I64(subOrderId), limit)
	if err != nil {
		logrus.Errorf("GetSitterLocationsBySubOrderID failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]*model.SitterLocationResponse, 0, len(locations))
	for _, loc := range locations {
		response = append(response, &model.SitterLocationResponse{
			ID:         loc.ID,
			SitterID:   loc.SitterID,
			OrderID:    loc.OrderID,
			SubOrderID: loc.SubOrderID,
			Lat:        loc.Lat,
			Lng:        loc.Lng,
			Timestamp:  loc.Timestamp,
			CreatedAt:  loc.CreatedAt.String(),
		})
	}

	open_api.OpenApiSuccessResponse(c, response)
}
