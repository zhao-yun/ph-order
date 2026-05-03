package model

// OrderStatus 定义订单状态
type OrderStatus int64

const (
	OrderInitialized  OrderStatus = 0  // 订单初始化(未支付）
	OrderPayed        OrderStatus = 1  // 订单已支付授权
	OrderAccepted     OrderStatus = 2  // Sitter 接受，此状态后订单无法修改
	OrderRejected     OrderStatus = -1 // Sitter 拒绝为
	OrderTimeout      OrderStatus = -2 // Sitter 超时未接受为
	OrderCancel       OrderStatus = -3 // 用户取消订单
	OrderRefund       OrderStatus = -4 // 订单已退款
	OrderEstablished  OrderStatus = 3  // 订单开始（Sitter输入验证码）
	OrderPreCompleted OrderStatus = 4  // 订单预结束，Sitter点击订单完成
	OrderCompleted    OrderStatus = 5  // 输入验证码后订单结束
)

// String 方法用于将 OrderStatus 转换为字符串表示='/
func (os OrderStatus) String() string {
	switch os {
	case OrderInitialized:
		return "Order Initialized"
	case OrderAccepted:
		return "Order Accepted"
	case OrderRejected:
		return "Order Rejected"
	case OrderTimeout:
		return "Order Timeout"
	case OrderEstablished:
		return "Order Established"
	case OrderPreCompleted:
		return "Order PreCompleted"
	case OrderCompleted:
		return "Order Completed"
	default:
		return "Unknown Order Status"
	}
}

// OrderModificationStatus 定义订单修改状态
type OrderModificationStatus int64

const (
	OrderModificationInitialized OrderModificationStatus = 0
	OrderModificationAccepted    OrderModificationStatus = 1
	OrderModificationRejected    OrderModificationStatus = -1
	OrderModificationTimeout     OrderModificationStatus = -2
)

func (os OrderModificationStatus) String() string {
	switch os {
	case OrderModificationInitialized:
		return "OrderModification Initialized"
	case OrderModificationAccepted:
		return "OrderModification Accepted"
	case OrderModificationRejected:
		return "OrderModification Rejected"
	case OrderModificationTimeout:
		return "OrderModification Timeout"
	}

	return "Unknown OrderModification Status"
}

// OrderModificationType 定义订单修改类型
type OrderModificationType int64

const (
	OrderModificationTypeUser   OrderModificationType = 1
	OrderModificationTypeSitter OrderModificationType = 2
)

func (ot OrderModificationType) String() string {
	switch ot {
	case OrderModificationTypeUser:
		return "OrderModification User"
	case OrderModificationTypeSitter:
		return "OrderModification Sitter"
	}
	return "Unknown OrderModification Type"

}
