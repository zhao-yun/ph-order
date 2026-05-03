package model

import (
	"demo/util/common"

	"github.com/go-playground/validator/v10"
)

// InterviewRecord 面试记录模型
type InterviewRecord struct {
	ID              int64             `json:"ID"`
	OrderID         int64             `json:"OrderID"`         // 关联的订单ID，用于绑定业务订单
	InitiatorType   InitiatorType     `json:"InitiatorType"`   // 预约发起方类型：1-用户 2-Sitter
	InterviewType   InterviewType     `json:"InterviewType"`   // 面试类型：1-线上 2-离线会议
	AppointmentTime *common.Timestamp `json:"AppointmentTime"` // 预约面试时间
	Location        string            `json:"Location"`        // 面试地点（线上则为空）
	Message         string            `json:"Message"`         // 预约时填写的消息内容
	Status          InterviewStatus   `json:"Status"`          // 预约状态：1-待确认（刚发起） 2-接受方修改待确认 3-已接受 4-已取消 5-已拒绝
	CreatedAt       common.Timestamp  `json:"CreatedAt"`       // 预约创建时间
	UpdatedAt       common.Timestamp  `json:"UpdatedAt"`       // 记录更新时间
	UserResult      *InterviewResult  `json:"UserResult"`      // 用户面试结果
	SitterResult    *InterviewResult  `json:"SitterResult"`    // Sitter面试结果
	UserReason      string            `json:"UserReason"`      // 用户拒绝原因
	SitterReason    string            `json:"SitterReason"`    // Sitter拒绝原因
}

// TableName 设置表名
func (InterviewRecord) TableName() string {
	return "interview_record"
}

// InitiatorType 预约发起方类型枚举
type InitiatorType int

const (
	InitiatorType_User   InitiatorType = 1 // 用户发起
	InitiatorType_Sitter InitiatorType = 2 // Sitter发起
)

// 转换为中文描述（可选，用于日志/返回给前端）
func (t InitiatorType) String() string {
	return map[InitiatorType]string{
		InitiatorType_User:   "用户",
		InitiatorType_Sitter: "Sitter",
	}[t]
}

// InterviewType 面试类型枚举
type InterviewType int

const (
	InterviewType_Online  InterviewType = 1 // 线上面试
	InterviewType_Offline InterviewType = 2 // 离线会议（线下）
)

func (t InterviewType) String() string {
	return map[InterviewType]string{
		InterviewType_Online:  "线上",
		InterviewType_Offline: "离线会议",
	}[t]
}

// InterviewStatus 预约状态枚举
type InterviewStatus int

const (
	Status_NotStart        InterviewStatus = 0 // 未开始，发起面试
	Status_Pending         InterviewStatus = 1 // 待确认（刚发起），发起人可以取消，接受方可以接受或拒绝
	Status_ModifiedPending InterviewStatus = 2 // 接受方修改待确认
	Status_Accepted        InterviewStatus = 3 // 已接受，进入面试流程
	Status_Canceled        InterviewStatus = 4 // 已取消，可以重新发起面试
	Status_Rejected        InterviewStatus = 5 // 已拒绝，拒绝后可以重新发起面试
	Status_Pass            InterviewStatus = 6 // 面试通过
	Status_Failed          InterviewStatus = 7 // 面试不通过
)

type InterviewResult int

const (
	Result_PASS  InterviewResult = 1 // 面试通过
	Result_FAILD InterviewResult = 2 // 面试不通过
)

func (s InterviewStatus) String() string {
	return map[InterviewStatus]string{
		Status_Pending:         "待确认",
		Status_ModifiedPending: "接受方修改待确认",
		Status_Accepted:        "已接受",
		Status_Canceled:        "已取消",
		Status_Rejected:        "已拒绝",
	}[s]
}

// CreateInterviewReq 发起面试预约请求参数
type CreateInterviewReq struct {
	OrderID         int64             `json:"OrderID" binding:"required" validate:"max=64"`          // 订单ID（必填，最大64字符）
	InitiatorType   InitiatorType     `json:"InitiatorType" binding:"required" validate:"oneof=1 2"` // 发起方类型（1-用户，2-Sitter）
	InterviewType   InterviewType     `json:"InterviewType" binding:"required" validate:"oneof=1 2"` // 面试类型（1-线上，2-离线）
	AppointmentTime *common.Timestamp `json:"AppointmentTime" binding:"required"`                    // 预约时间（必填）
	Location        string            `json:"Location" validate:"max=255"`                           // 面试地点
	Message         string            `json:"Message" validate:"max=500"`                            // 预约消息
}

// GetInterviewByOrderIDReq 根据订单ID查询面试预约请求参数
type GetInterviewByOrderIDReq struct {
	OrderID string `json:"OrderID" binding:"required" validate:"max=64"` // 订单ID（必填）
}

// UpdateInterviewStatusReq 更新面试预约状态请求参数
type UpdateInterviewStatusReq struct {
	ID     int64           `json:"ID" binding:"required"`                                // 预约ID（必填）
	Status InterviewStatus `json:"Status" binding:"required" validate:"oneof=1 2 3 4 5"` // 目标状态（1-待确认，2-修改待确认，3-已接受，4-已取消，5-已拒绝）
}

// ModifyInterviewReq 修改面试预约请求参数
type ModifyInterviewReq struct {
	ID              int64             `json:"ID" binding:"required"`               // 预约ID（必填）
	InitiatorType   InitiatorType     `json:"InitiatorType"  validate:"oneof=1 2"` // 发起方类型（1-用户，2-Sitter）
	AppointmentTime *common.Timestamp `json:"AppointmentTime"`                     // 新预约时间（必填）
	Location        string            `json:"Location" validate:"max=255"`         // 新面试地点（最大255字符）
	Message         string            `json:"Message" validate:"max=500"`          // 新预约消息（最大500字符）
}

// ConfirmInterviewModificationReq 确认预约修改请求参数
type ConfirmInterviewModificationReq struct {
	ID            int64         `json:"ID" binding:"required"`                                 // 预约ID（必填）
	Confirm       bool          `json:"Confirm" binding:"required"`                            // 是否确认修改（必填）
	InitiatorType InitiatorType `json:"InitiatorType" binding:"required" validate:"oneof=1 2"` // 发起方类型（1-用户，2-Sitter）
}

// CancelInterviewReq 取消面试预约请求参数
type CancelInterviewReq struct {
	ID            int64         `json:"ID" binding:"required"`                                 // 预约ID（必填）
	Reason        string        `json:"Reason" validate:"max=500"`                             // 取消原因（最大500字符）
	InitiatorType InitiatorType `json:"InitiatorType" binding:"required" validate:"oneof=1 2"` // 发起方类型（1-用户，2-Sitter）
}

// GetInterviewResultReq 获取面试结果请求参数
type GetInterviewResultReq struct {
	OrderID string `json:"OrderID" binding:"required" validate:"max=64"` // 订单ID（必填）
}

// SubmitInterviewResultReq 提交面试结果请求参数
type SubmitInterviewResultReq struct {
	ID     int64           `json:"ID" binding:"required" validate:"max=64"`          // 面试记录ID（必填）
	Result InterviewResult `json:"Result" binding:"required" validate:"oneof=0 1 2"` // 面试结果（0-未操作，1-通过，2-不通过）
	Reason string          `json:"Reason" validate:"max=500"`                        // 不通过原因（选填，最大500字符）
}

// ReplyInterviewReq 接受/拒绝面试预约请求参数
type ReplyInterviewReq struct {
	ID            int64           `json:"ID" binding:"required"`                                // 预约ID（必填）
	Status        InterviewStatus `json:"Status" binding:"required" validate:"oneof=1 2 3 4 5"` // 3-接受，5-拒绝）
	InitiatorType InitiatorType   `json:"InitiatorType" binding:"required" validate:"oneof=1 2"`
}

// 自定义验证错误处理（示例）
func (r CreateInterviewReq) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
