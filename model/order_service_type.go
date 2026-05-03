package model

import (
	"fmt"
	"strings"
)

type OrderServiceType int64

const (
	OrderServiceTypeUnknown OrderServiceType = 0
	OrderServiceTypeBoarding OrderServiceType = 1
	OrderServiceTypeDaycare  OrderServiceType = 2
	OrderServiceTypeWalking  OrderServiceType = 3
	OrderServiceTypeDropIn   OrderServiceType = 4
)

func (t OrderServiceType) String() string {
	switch t {
	case OrderServiceTypeBoarding:
		return "Boarding"
	case OrderServiceTypeDaycare:
		return "Daycare"
	case OrderServiceTypeWalking:
		return "Walking"
	case OrderServiceTypeDropIn:
		return "DropIn"
	default:
		return "Unknown"
	}
}

func ParseOrderServiceType(s string) (OrderServiceType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "boarding":
		return OrderServiceTypeBoarding, nil
	case "daycare", "day care":
		return OrderServiceTypeDaycare, nil
	case "walking", "walk":
		return OrderServiceTypeWalking, nil
	case "dropin", "drop-in", "drop in":
		return OrderServiceTypeDropIn, nil
	default:
		return OrderServiceTypeUnknown, fmt.Errorf("unknown service type: %s", s)
	}
}

