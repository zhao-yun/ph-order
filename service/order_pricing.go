package service

import (
	"fmt"
	"strings"
	"time"

	"demo/model"
)

type OrderPricingResult struct {
	ServiceType      model.OrderServiceType
	Units            int
	PetTotalsCents   map[int]*int64
	PetsSubtotalCents int64
	AddonsCents      int64
	TotalCents       int64
}

func CalculateOrderPricing(req *model.CreateOrderWithPricingReq) (*OrderPricingResult, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if req.SitterID == "" {
		return nil, fmt.Errorf("SitterID is empty")
	}
	serviceType := model.OrderServiceType(req.Type)
	if serviceType != model.OrderServiceTypeBoarding &&
		serviceType != model.OrderServiceTypeDaycare &&
		serviceType != model.OrderServiceTypeWalking &&
		serviceType != model.OrderServiceTypeDropIn {
		return nil, fmt.Errorf("invalid Type: %d", req.Type)
	}
	if req.FromDate.IsZero() || req.ToDate.IsZero() {
		return nil, fmt.Errorf("FromDate/ToDate is required")
	}
	if req.ToDate.Time.Before(req.FromDate.Time) {
		return nil, fmt.Errorf("ToDate must be >= FromDate")
	}
	if len(req.PetList) == 0 {
		return nil, fmt.Errorf("PetList is empty")
	}

	pricing, err := GetServicePricingBySitterIDAndType(req.SitterID, serviceType)
	if err != nil {
		return nil, err
	}
	if !pricing.Enabled {
		return nil, fmt.Errorf("service type %s is disabled", serviceType.String())
	}

	units := dateUnitsInclusive(req.FromDate.Time, req.ToDate.Time)
	if units <= 0 {
		return nil, fmt.Errorf("invalid date range")
	}

	if serviceType == model.OrderServiceTypeWalking {
		if req.DurationMinutes == 0 {
			req.DurationMinutes = 30
		}
		if req.DurationMinutes != 30 && req.DurationMinutes != 60 {
			return nil, fmt.Errorf("invalid DurationMinutes: %d", req.DurationMinutes)
		}
	}

	perPetTotals := make(map[int]*int64, len(req.PetList))
	for i := range req.PetList {
		var v int64
		perPetTotals[i] = &v
	}

	isLongTerm := serviceType == model.OrderServiceTypeBoarding && units >= 7

	for d := 0; d < units; d++ {
		curDate := req.FromDate.Time.AddDate(0, 0, d)
		isHoliday := isHolidayCN(curDate)

		catIndex := 0
		dogIndex := 0
		for i, pet := range req.PetList {
			pt := normalizePetType(pet.PetType)
			switch pt {
			case "cat":
				if pricing.CatExclusiveCents == 0 && pricing.AdditionalCatCents == 0 {
					return nil, fmt.Errorf("cats are not supported for %s", serviceType.String())
				}
				var cents int64
				if catIndex == 0 {
					cents = pricing.CatExclusiveCents
				} else {
					if pricing.AdditionalCatCents == 0 {
						return nil, fmt.Errorf("additional cat fee is not configured for %s", serviceType.String())
					}
					cents = pricing.AdditionalCatCents
				}
				*perPetTotals[i] += cents
				catIndex++
			case "dog":
				if pricing.DogBaseCents == 0 && pricing.AdditionalDogCents == 0 {
					return nil, fmt.Errorf("dogs are not supported for %s", serviceType.String())
				}
				var cents int64
				if dogIndex == 0 {
					cents, err = dogBaseCentsForUnit(pricing, serviceType, pet, isHoliday, isLongTerm, req.DurationMinutes)
					if err != nil {
						return nil, err
					}
				} else {
					if pricing.AdditionalDogCents == 0 {
						return nil, fmt.Errorf("additional dog fee is not configured for %s", serviceType.String())
					}
					cents = pricing.AdditionalDogCents
				}
				*perPetTotals[i] += cents
				dogIndex++
			default:
				return nil, fmt.Errorf("unsupported PetType: %s", pet.PetType)
			}
		}
	}

	var petsSubtotal int64
	for _, v := range perPetTotals {
		petsSubtotal += *v
	}

	var addons int64
	if req.NeedGrooming {
		addons += pricing.GroomingCents
	}
	if req.NeedPickup {
		addons += pricing.PickupCents
	}

	return &OrderPricingResult{
		ServiceType:      serviceType,
		Units:            units,
		PetTotalsCents:   perPetTotals,
		PetsSubtotalCents: petsSubtotal,
		AddonsCents:      addons,
		TotalCents:       petsSubtotal + addons,
	}, nil
}

func dateUnitsInclusive(from time.Time, to time.Time) int {
	from = truncateDate(from)
	to = truncateDate(to)
	if to.Before(from) {
		return 0
	}
	return int(to.Sub(from).Hours()/24) + 1
}

func truncateDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func normalizePetType(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "cat", "猫", "kitty", "kitten":
		return "cat"
	case "dog", "狗", "puppy":
		return "dog"
	default:
		return s
	}
}

func dogBaseCentsForUnit(pricing model.ServicePricing, serviceType model.OrderServiceType, pet *model.OrderPet, isHoliday bool, isLongTerm bool, durationMinutes int) (int64, error) {
	if serviceType == model.OrderServiceTypeWalking && durationMinutes == 60 {
		if pricing.Walking60MinCents == 0 {
			return 0, fmt.Errorf("60 minutes price is not configured for walking")
		}
		return pricing.Walking60MinCents, nil
	}
	if isHoliday && pricing.HolidayCents != 0 {
		return pricing.HolidayCents, nil
	}
	if serviceType == model.OrderServiceTypeBoarding && isLongTerm && pricing.LongTermCents != 0 {
		return pricing.LongTermCents, nil
	}
	if pet != nil && pet.IsPuppy && pricing.PuppyCents != 0 {
		return pricing.PuppyCents, nil
	}
	if pricing.DogBaseCents == 0 {
		return 0, fmt.Errorf("dog base price is not configured for %s", serviceType.String())
	}
	return pricing.DogBaseCents, nil
}

func isHolidayCN(date time.Time) bool {
	y := date.Year()
	if y != 2026 {
		return false
	}
	d := truncateDate(date)

	ranges := [][2]time.Time{
		{time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)},
		{time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC), time.Date(2026, 2, 23, 0, 0, 0, 0, time.UTC)},
		{time.Date(2026, 4, 4, 0, 0, 0, 0, time.UTC), time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC)},
		{time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)},
		{time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 21, 0, 0, 0, 0, time.UTC)},
		{time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC), time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC)},
		{time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 10, 7, 0, 0, 0, 0, time.UTC)},
	}

	for _, r := range ranges {
		if !d.Before(r[0]) && !d.After(r[1]) {
			return true
		}
	}
	return false
}
