package service

import (
	"testing"
	"time"

	"demo/model"
	"demo/util/common"
)

func d(y int, m time.Month, day int) common.Date {
	return common.Date{Time: time.Date(y, m, day, 0, 0, 0, 0, time.UTC)}
}

func TestCalculateOrderPricing_CatsBoarding(t *testing.T) {
	req := &model.CreateOrderWithPricingReq{
		SitterID: "s1",
		Type:     1,
		FromDate: d(2026, 3, 1),
		ToDate:   d(2026, 3, 1),
		PetList: []*model.OrderPet{
			{PetType: "Cat"},
			{PetType: "Cat"},
			{PetType: "Cat"},
		},
	}

	res, err := CalculateOrderPricing(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res.PetsSubtotalCents != 8500 {
		t.Fatalf("want 8500, got %d", res.PetsSubtotalCents)
	}
	if res.TotalCents != 8500 {
		t.Fatalf("want 8500, got %d", res.TotalCents)
	}
}

func TestCalculateOrderPricing_BoardingDogMultiDay(t *testing.T) {
	req := &model.CreateOrderWithPricingReq{
		SitterID: "s1",
		Type:     1,
		FromDate: d(2026, 3, 1),
		ToDate:   d(2026, 3, 3),
		PetList: []*model.OrderPet{
			{PetType: "Dog"},
		},
	}

	res, err := CalculateOrderPricing(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res.Units != 3 {
		t.Fatalf("want 3 units, got %d", res.Units)
	}
	if res.PetsSubtotalCents != 15000 {
		t.Fatalf("want 15000, got %d", res.PetsSubtotalCents)
	}
}

func TestCalculateOrderPricing_BoardingDogLongTerm(t *testing.T) {
	req := &model.CreateOrderWithPricingReq{
		SitterID: "s1",
		Type:     1,
		FromDate: d(2026, 3, 1),
		ToDate:   d(2026, 3, 8),
		PetList: []*model.OrderPet{
			{PetType: "Dog"},
		},
	}

	res, err := CalculateOrderPricing(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res.Units != 8 {
		t.Fatalf("want 8 units, got %d", res.Units)
	}
	if res.PetsSubtotalCents != 32000 {
		t.Fatalf("want 32000, got %d", res.PetsSubtotalCents)
	}
}

func TestCalculateOrderPricing_Walking60Minutes(t *testing.T) {
	req := &model.CreateOrderWithPricingReq{
		SitterID:        "s1",
		Type:            3,
		FromDate:        d(2026, 3, 1),
		ToDate:          d(2026, 3, 1),
		DurationMinutes: 60,
		PetList: []*model.OrderPet{
			{PetType: "Dog"},
		},
	}

	res, err := CalculateOrderPricing(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res.PetsSubtotalCents != 4000 {
		t.Fatalf("want 4000, got %d", res.PetsSubtotalCents)
	}
}

func TestCalculateOrderPricing_BoardingHoliday(t *testing.T) {
	req := &model.CreateOrderWithPricingReq{
		SitterID: "s1",
		Type:     1,
		FromDate: d(2026, 10, 1),
		ToDate:   d(2026, 10, 1),
		PetList: []*model.OrderPet{
			{PetType: "Dog"},
		},
	}

	res, err := CalculateOrderPricing(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res.PetsSubtotalCents != 7500 {
		t.Fatalf("want 7500, got %d", res.PetsSubtotalCents)
	}
}
