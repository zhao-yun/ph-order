package service

import (
	"fmt"

	"demo/model"
)

func GetSitterPricingBySitterID(sitterID string) (*model.SitterPricing, error) {
	if sitterID == "" {
		return nil, fmt.Errorf("sitterID is empty")
	}

	return &model.SitterPricing{
		SitterID: sitterID,
		Services: map[model.OrderServiceType]model.ServicePricing{
			model.OrderServiceTypeBoarding: {
				ServiceType:        model.OrderServiceTypeBoarding,
				Enabled:            true,
				DogBaseCents:       5000,
				CatExclusiveCents:  4500,
				AdditionalCatCents: 2000,
				AdditionalDogCents: 2500,
				PuppyCents:         6000,
				LongTermCents:      4000,
				HolidayCents:       7500,
				GroomingCents:      1500,
				PickupCents:        1000,
			},
			model.OrderServiceTypeDaycare: {
				ServiceType:        model.OrderServiceTypeDaycare,
				Enabled:            true,
				DogBaseCents:       3500,
				CatExclusiveCents:  3000,
				AdditionalCatCents: 1200,
				AdditionalDogCents: 1500,
				PuppyCents:         4000,
				HolidayCents:       5000,
				GroomingCents:      1500,
				PickupCents:        800,
			},
			model.OrderServiceTypeWalking: {
				ServiceType:       model.OrderServiceTypeWalking,
				Enabled:           true,
				DogBaseCents:      2500,
				AdditionalDogCents: 1000,
				PuppyCents:        3000,
				HolidayCents:      3500,
				Walking60MinCents: 4000,
			},
			model.OrderServiceTypeDropIn: {
				ServiceType:        model.OrderServiceTypeDropIn,
				Enabled:            true,
				DogBaseCents:       2000,
				CatExclusiveCents:  1800,
				AdditionalCatCents: 600,
				AdditionalDogCents: 800,
				PuppyCents:         2500,
				HolidayCents:       3000,
				GroomingCents:      1200,
			},
		},
	}, nil
}

func GetServicePricingBySitterIDAndType(sitterID string, serviceType model.OrderServiceType) (model.ServicePricing, error) {
	sitterPricing, err := GetSitterPricingBySitterID(sitterID)
	if err != nil {
		return model.ServicePricing{}, err
	}
	pricing, ok := sitterPricing.Services[serviceType]
	if !ok {
		return model.ServicePricing{}, fmt.Errorf("pricing not found for service type: %s", serviceType.String())
	}
	return pricing, nil
}

