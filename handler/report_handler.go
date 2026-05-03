package handler

import (
	"net/http"
	"strconv"

	"demo/model"
	"demo/service"
	"demo/util/auth"
	"demo/util/open_api"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct{}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{}
}

// SubmitReport 提交举报
// @Summary 提交举报
// @Description 用户提交举报信息
// @Tags 举报管理
// @Accept json
// @Produce json
// @Param report body model.ReportRequest true "举报信息"
// @Success 200 {object} model.Response "成功返回"
// @Failure 400 {object} model.Response "参数错误"
// @Failure 500 {object} model.Response "内部错误"
// @Router /reports [post]
func SubmitReport(c *gin.Context) {
	var req model.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	report, err := service.SubmitReport(c.Request.Context(), &req, userID)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "提交举报失败: "+err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, report)
}

// GetReportDetail 获取举报详情
// @Summary 获取举报详情
// @Description 获取举报的详细信息
// @Tags 举报管理
// @Produce json
// @Param id path int true "举报ID"
// @Success 200 {object} model.Response{data=model.ReportResponse} "成功返回"
// @Failure 400 {object} model.Response "参数错误"
// @Failure 403 {object} model.Response "无权访问"
// @Failure 404 {object} model.Response "举报不存在"
// @Router /reports/{id} [get]
func GetReportDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "举报ID格式错误")
		return
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	report, err := service.GetReportDetail(c.Request.Context(), id, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "举报记录不存在或无权查看" {
			status = http.StatusNotFound
		}
		open_api.OpenApiErrorResponse(c, status, err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, report)
}

// GetMyReports 获取我的举报记录
// @Summary 获取我的举报记录
// @Description 获取当前用户提交的所有举报
// @Tags 举报管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} model.Response{data=[]model.ReportResponse,total=int64} "成功返回"
// @Failure 500 {object} model.Response "内部错误"
// @Router /reports/my [get]
func GetMyReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	userID, err := auth.GetUserID(c)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	reports, total, err := service.GetMyReports(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "获取举报记录失败: "+err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, gin.H{
		"list":  reports,
		"total": total,
	})
}

// GetReportsAgainstMe 获取针对我的举报
// @Summary 获取针对我的举报
// @Description 获取其他用户对当前用户的举报
// @Tags 举报管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} model.Response{data=[]model.ReportResponse,total=int64} "成功返回"
// @Failure 500 {object} model.Response "内部错误"
// @Router /reports/against-me [get]
func GetReportsAgainstMe(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	userID, err := auth.GetUserID(c)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	reports, total, err := service.GetReportsAgainstMe(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "获取举报记录失败: "+err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, gin.H{
		"list":  reports,
		"total": total,
	})
}

// ProcessReport 处理举报(管理员)
// @Summary 处理举报
// @Description 管理员开始处理举报
// @Tags 举报管理(管理员)
// @Produce json
// @Param id path int true "举报ID"
// @Success 200 {object} model.Response "成功返回"
// @Failure 400 {object} model.Response "参数错误"
// @Failure 500 {object} model.Response "内部错误"
// @Router /admin/reports/{id}/process [put]
func ProcessReport(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "举报ID格式错误")
		return
	}

	handlerID, err := auth.GetUserID(c)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	if err := service.ProcessReport(c.Request.Context(), id, handlerID); err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "处理举报失败: "+err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, nil)
}

// ResolveReport 解决举报(管理员)
// @Summary 解决举报
// @Description 管理员提交举报处理结果
// @Tags 举报管理(管理员)
// @Accept json
// @Produce json
// @Param id path int true "举报ID"
// @Param result body string true "处理结果"
// @Success 200 {object} model.Response "成功返回"
// @Failure 400 {object} model.Response "参数错误"
// @Failure 500 {object} model.Response "内部错误"
// @Router /admin/reports/{id}/resolve [put]
func ResolveReport(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "举报ID格式错误")
		return
	}

	var result struct {
		Result string `json:"result" binding:"required"`
	}
	if err := c.ShouldBindJSON(&result); err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := service.ResolveReport(c.Request.Context(), id, result.Result); err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "解决举报失败: "+err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, nil)
}

// GetPendingReports 获取待处理举报(管理员)
// @Summary 获取待处理举报
// @Description 管理员获取所有待处理的举报
// @Tags 举报管理(管理员)
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} model.Response{data=[]model.ReportResponse,total=int64} "成功返回"
// @Failure 500 {object} model.Response "内部错误"
// @Router /admin/reports/pending [get]
func GetPendingReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	reports, total, err := service.GetPendingReports(c.Request.Context(), page, pageSize)
	if err != nil {
		open_api.OpenApiErrorResponse(c, http.StatusInternalServerError, "获取待处理举报失败: "+err.Error())
		return
	}

	open_api.OpenApiSuccessResponse(c, gin.H{
		"list":  reports,
		"total": total,
	})
}
