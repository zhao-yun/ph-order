package model

import "demo/util/common"

type SubOrder struct {
	ID                 int64           `gorm:"primaryKey;autoIncrement;comment('子订单ID')"`
	OrderID            int64           `gorm:"not null;comment('主订单ID')"`
	Date               common.Date     `gorm:"type:date;not null;comment('日期')"`
	State              OrderStatus     `gorm:"not null;comment('子订单状态')"`
	StartCode          string          `gorm:"type:varchar(10);comment('开始验证码')"`
	EndCode            string          `gorm:"type:varchar(10);comment('结束验证码')"`
	SitterHandleAt     common.Timestamp `gorm:"type:timestamp;comment('Sitter处理时间')"`
	SitterFinishAt     common.Timestamp `gorm:"type:timestamp;comment('Sitter完成时间')"`
	WalkThumbnailUrl   string          `gorm:"type:varchar(500);comment('轨迹缩略图URL')" json:"WalkThumbnailUrl"`
	CreatedAt          common.Timestamp `gorm:"type:timestamp;autoCreateTime;comment('创建时间')"`
	UpdatedAt          common.Timestamp `gorm:"type:timestamp;autoUpdateTime;comment('更新时间')"`
	DayNumber          int             `gorm:"-" json:"DayNumber"`
}

func (SubOrder) TableName() string {
	return "sub_orders"
}

type SubOrderPage struct {
	Total        int64        `json:"Total"`
	TotalPage    int          `json:"TotalPage"`
	Current      int          `json:"Current"`
	Size         int          `json:"Size"`
	SubOrderList []*SubOrder  `json:"SubOrderList"`
}

type SubOrderQueryParams struct {
	ID        int64             `json:"ID" form:"id"`
	OrderID   int64             `json:"OrderID" form:"order_id"`
	Date      *common.Date      `json:"Date" form:"date"`
	State     OrderStatus       `json:"State" form:"state"`
	StartTime *common.Timestamp `json:"StartTime"`
	EndTime   *common.Timestamp `json:"EndTime"`
}

type SitterSetSubOrderStartCodeReq struct {
	SubOrderID int64  `json:"SubOrderID" form:"SubOrderID"`
	Code       string `json:"Code" form:"Code"`
}

type SitterSetSubOrderEndCodeReq struct {
	SubOrderID int64  `json:"SubOrderID" form:"SubOrderID"`
	Code       string `json:"Code" form:"Code"`
}

type UpdateSubOrderWalkThumbnailReq struct {
	SubOrderID int64  `json:"SubOrderID" form:"SubOrderID" binding:"required"`
	Url        string `json:"Url" form:"Url" binding:"required"`
}
