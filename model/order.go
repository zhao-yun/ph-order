package model

import (
	"demo/util/common"
)

type Order struct {
	ID                int64                 `gorm:"primaryKey;autoIncrement;comment('订单ID')"`
	OwnerID           string                `gorm:"type:varchar(50);not null;comment('用户ID')"`
	SitterID          string                `gorm:"type:varchar(50);not null;comment('Sitter ID')"`
	Type              int64                 `gorm:"not null;comment('订单类型')"`
	FromDate          common.Date           `gorm:"type:date;not null;comment('开始时间')"`
	ToDate            common.Date           `gorm:"type:date;not null;comment('结束时间')"`
	TipsPrice         float64               `gorm:"type:decimal;comment('小费')"`
	TotalPrice        float64               `gorm:"type:decimal;comment('总费用')"`
	SubTotalPrice     float64               `gorm:"type:decimal;comment('宠物总费用')"`
	ServiceFee        float64               `gorm:"type:decimal;comment('服务费')"`
	Taxes             float64               `gorm:"type:decimal;comment('税')"`
	State             OrderStatus           `gorm:"not null;comment('订单状态')"`
	SitterHandleAt    common.Timestamp      `gorm:"type:timestamp;comment('Sitter 处理时间')"`
	SitterFinishAt    common.Timestamp      `gorm:"type:timestamp;comment('Sitter 完成时间')"`
	CreatedAt         common.Timestamp      `gorm:"type:timestamp;autoCreateTime;comment('创建时间')"`
	UpdatedAt         common.Timestamp      `gorm:"type:timestamp;autoUpdateTime;comment('更新时间')"`
	PetList           []*OrderPet           `gorm:"-"`
	ModificationLog   *OrderModificationLog `gorm:"-"`
	OwnerName         string                `gorm:"type:varchar(255);comment('用户名字')"`
	SitterName        string                `gorm:"type:varchar(255);comment('Sitter名字')"`
	Contact           string                `gorm:"type:varchar(255);comment('联系人')"`
	AlternativeContact string               `gorm:"type:varchar(255);comment('备用联系人')"`
	Note              string                `gorm:"type:text;comment('备注')"`
	CancelAt          common.Timestamp      `gorm:"type:timestamp;comment('取消时间')"`
	CancelReason      string                `gorm:"type:varchar(255);comment('取消原因')"`
	RefundPrice       float64               `gorm:"type:decimal;comment('退款费用')"`
	OrderNumber       string                `gorm:"type:varchar(255);comment('订单编号')"`
	UserDeleted       int                   `gorm:"default:0;comment('用户是否删除')"`
	SitterDeleted     int                   `gorm:"default:0;comment('Sitter是否删除')"`
	UserCode          string                `gorm:"-"`
	UserCodeState     int64                 `gorm:"-"`
	SitterCodeState   int64                 `gorm:"-"`
	Code              string                `gorm:"-"`
	UserRatingState   int64                 `gorm:"default:0;comment('用户评分状态')"`
	SitterRatingState int                   `gorm:"default:0;comment('Sitter评分状态')"`
	OwnerRating       *OwnerRating          `gorm:"-"`
	SitterRating      *SitterRating         `gorm:"-"`
	InterviewStatus   *InterviewStatus      `gorm:"-"`
}

// TableName 设置表名
func (Order) TableName() string {
	return "orders"
}

type OrderPage struct {
	Total     int64    `json:"Total`
	TotalPage int      `json:"TotalPage"`
	Current   int      `json:"Current"`
	Size      int      `json:"Size"`
	OrderList []*Order `json:"OrderList"`
}

type OrderQueryParams struct {
	ID            int64             `json:"ID" form:"id"`                        // 订单ID
	OwnerID       string            `json:"OwnerID" form:"owner_id"`             // 用户ID
	SitterID      string            `json:"SitterID" form:"sitter_id"`           // Sitter ID
	TypeList      []int64           `json:"TypeList" form:"type"`                // 订单类型
	State         OrderStatus       `json:"State" form:"state"`                  // 订单状态
	OrderTime     *common.Timestamp `json:"OrderTime"`                           // 订单时间（最近时间）
	StartTime     *common.Timestamp `json:"StartTime"`                           // 订单时间（最早时间）
	EndTime       *common.Timestamp `json:"EndTime"`                             // 订单时间（最晚时间）
	PetTypeList   []string          `json:"PetTypeList"`                         // 宠物类型
	IDList        []int64           `json:"IDList"`                              // 订单ID列表
	PetIDList     []string          `json:"PetIDList"`                           // 宠物ID列表
	Keyword       string            `json:"Keyword" form:"keyword"`              // 搜索关键词
	KeyIDList     []int64           `json:"KeyIDList"`                           // 关键词订单ID列表
	UserDeleted   *int              `json:"UserDeleted" form:"user_deleted"`     // 用户是否删除
	SitterDeleted *int              `json:"SitterDeleted" form:"sitter_deleted"` // Sitter是否删除
}

type SitterHandleOrderInviteReq struct {
	OrderID int64       `json:"OrderID" form:"OrderID"` // 订单ID
	State   OrderStatus `json:"State" form:"State"`     // 状态
}

type UserSetCreateCodeReq struct {
	OrderID int64  `json:"OrderID" form:"OrderID"`
	Code    string `json:"Code" form:"Code"` // 邀请码
}

type SitterSetCreateCodeReq struct {
	OrderID int64  `json:"OrderID" form:"OrderID"`
	Code    string `json:"Code" form:"Code"` // 邀请码
}

type SitterFinishOrderReq struct {
	OrderID int64 `json:"OrderID" form:"OrderID"`
}

type SitterSetFinishCodeReq struct {
	OrderID int64  `json:"OrderID" form:"OrderID"`
	Code    string `json:"Code" form:"Code"` // 邀请码
}

type UserCancelOrderReq struct {
	OrderID int64  `json:"OrderID" form:"OrderID"`
	Reason  string `json:"Reason" form:"Reason"`
}

type AddTipsReq struct {
	OrderID int64   `json:"OrderID" form:"OrderID"`
	Tips    float64 `json:"Tips" form:"Tips"`
}

type DeleteOrderReq struct {
	OrderID int64 `json:"OrderID" form:"OrderID"`
}

type PayOrderReq struct {
	OrderID int64 `json:"OrderID" form:"OrderID"`
}
