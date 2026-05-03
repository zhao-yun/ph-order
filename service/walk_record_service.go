package service

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"demo/dal"
	"demo/model"
	"demo/util/helper"
	open_api "demo/util/open_api"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const earthRadius = 6371000.0

func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	lat1Rad := lat1 * math.Pi / 180.0
	lng1Rad := lng1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0
	lng2Rad := lng2 * math.Pi / 180.0

	dLat := lat2Rad - lat1Rad
	dLng := lng2Rad - lng1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func calculatePathDistance(path []model.LatLng) float64 {
	if len(path) < 2 {
		return 0
	}

	totalDistance := 0.0
	for i := 1; i < len(path); i++ {
		distance := haversine(
			path[i-1].Lat, path[i-1].Lng,
			path[i].Lat, path[i].Lng,
		)
		totalDistance += distance
	}
	return totalDistance
}

func calculatePathDuration(path []model.LatLng) int64 {
	if len(path) < 2 {
		return 0
	}

	minTimestamp := path[0].Timestamp
	maxTimestamp := path[0].Timestamp

	for _, point := range path {
		if point.Timestamp < minTimestamp {
			minTimestamp = point.Timestamp
		}
		if point.Timestamp > maxTimestamp {
			maxTimestamp = point.Timestamp
		}
	}

	return maxTimestamp - minTimestamp
}

func AppendWalkPath(c *gin.Context) {
	var req model.CreateWalkRecordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("CreateWalkRecordReq ShouldBindJSON failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	subOrder, err := dal.GetSubOrderById(req.SubOrderID)
	if err != nil {
		logrus.Errorf("GetSubOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "sub order not found")
		return
	}

	if subOrder == nil {
		logrus.Errorf("sub order not found")
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "sub order not found")
		return
	}

	// 获取主订单以获取 SitterID
	order, err := dal.GetOrderById(subOrder.OrderID)
	if err != nil || order == nil {
		logrus.Errorf("GetOrderById failed or order not found, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "order not found")
		return
	}

	existingRecords, err := dal.GetWalkRecordsBySubOrderID(req.SubOrderID)
	if err != nil {
		logrus.Errorf("GetWalkRecordsBySubOrderID failed, err: %v", err)
	}

	var walkRecord *model.WalkRecord
	var existingPath []model.LatLng

	if len(existingRecords) > 0 {
		walkRecord = existingRecords[0]
		if len(walkRecord.Path) > 0 {
			if err := json.Unmarshal(walkRecord.Path, &existingPath); err != nil {
				logrus.Errorf("unmarshal existing path failed, err: %v", err)
			}
		}
	} else {
		walkRecord = &model.WalkRecord{
			OrderID:    subOrder.OrderID,
			SubOrderID: req.SubOrderID,
		}
	}

	combinedPath := append(existingPath, req.Path...)
	pathJson, err := json.Marshal(combinedPath)
	if err != nil {
		logrus.Errorf("marshal path failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "invalid path data")
		return
	}

	walkRecord.Path = pathJson

	if len(existingRecords) > 0 {
		err = dal.UpdateWalkRecord(walkRecord)
	} else {
		err = dal.CreateWalkRecord(walkRecord)
	}

	if err != nil {
		logrus.Errorf("save walk record failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "save walk record failed")
		return
	}

	// 更新 Sitter 实时位置（取传入路径点中最新的一个）
	if len(req.Path) > 0 {
		var latestPoint model.LatLng
		latestPoint = req.Path[0]
		for _, point := range req.Path {
			if point.Timestamp > latestPoint.Timestamp {
				latestPoint = point
			}
		}

		location := &model.SitterLocation{
			SitterID:   order.SitterID,
			OrderID:    subOrder.OrderID,
			SubOrderID: req.SubOrderID,
			Lat:        latestPoint.Lat,
			Lng:        latestPoint.Lng,
			Timestamp:  latestPoint.Timestamp,
		}

		err = dal.CreateSitterLocation(location)
		if err != nil {
			logrus.Errorf("CreateSitterLocation failed, err: %v", err)
			// 这里不中断返回，因为轨迹已经保存成功，只打印错误
		}
	}

	walkRecord.PathData = combinedPath
	open_api.OpenApiSuccessResponse(c, walkRecord)
}

func GetWalkRecordBySubOrderId(c *gin.Context) {
	subOrderId := c.Query("SubOrderID")
	if subOrderId == "" {
		logrus.Errorf("SubOrderID is required")
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "SubOrderID is required")
		return
	}

	walkRecordList, err := dal.GetWalkRecordsBySubOrderID(helper.S2I64(subOrderId))
	if err != nil {
		logrus.Errorf("GetWalkRecordsBySubOrderID failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	subOrder, err := dal.GetSubOrderById(helper.S2I64(subOrderId))
	if err != nil {
		logrus.Errorf("GetSubOrderById failed, err: %v", err)
	}

	if len(walkRecordList) == 0 {
		var walkThumbnailUrl string
		if subOrder != nil {
			walkThumbnailUrl = subOrder.WalkThumbnailUrl
		}
		result := &model.WalkDetail{
			WalkThumbnailUrl: walkThumbnailUrl,
		}
		open_api.OpenApiSuccessResponse(c, result)
		return
	}

	walkRecord := walkRecordList[0]

	var pathData []model.LatLng
	if len(walkRecord.Path) > 0 {
		if err := json.Unmarshal(walkRecord.Path, &pathData); err == nil {
			walkRecord.PathData = pathData
		}
	}

	durationSeconds := calculatePathDuration(pathData)
	distanceMeters := calculatePathDistance(pathData)

	var walkThumbnailUrl string
	if subOrder != nil {
		walkThumbnailUrl = subOrder.WalkThumbnailUrl
	}

	detail := &model.WalkDetail{
		ID:               strconv.FormatInt(walkRecord.ID, 10),
		OrderID:          walkRecord.OrderID,
		SubOrderID:       walkRecord.SubOrderID,
		Path:             walkRecord.PathData,
		DurationSeconds:  durationSeconds,
		DistanceMeters:   distanceMeters,
		WalkThumbnailUrl: walkThumbnailUrl,
		CreatedAt:        walkRecord.CreatedAt.String(),
	}

	open_api.OpenApiSuccessResponse(c, detail)
}

func UpdateSubOrderWalkThumbnail(c *gin.Context) {
	var req model.UpdateSubOrderWalkThumbnailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("UpdateSubOrderWalkThumbnailReq ShouldBindJSON failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	subOrder, err := dal.GetSubOrderById(req.SubOrderID)
	if err != nil {
		logrus.Errorf("GetSubOrderById failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "sub order not found")
		return
	}

	if subOrder == nil {
		logrus.Errorf("sub order not found")
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "sub order not found")
		return
	}

	subOrder.WalkThumbnailUrl = req.Url
	err = dal.UpdateSubOrder(subOrder)
	if err != nil {
		logrus.Errorf("UpdateSubOrder failed, err: %v", err)
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "update sub order failed")
		return
	}

	open_api.OpenApiSuccessResponse(c, subOrder)
}
