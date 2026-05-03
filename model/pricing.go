package model

type ServicePricing struct {
	ServiceType        OrderServiceType `json:"ServiceType"`
	Enabled            bool             `json:"Enabled"`
	DogBaseCents       int64            `json:"DogBaseCents"`
	CatExclusiveCents  int64            `json:"CatExclusiveCents"`
	AdditionalCatCents int64            `json:"AdditionalCatCents"`
	AdditionalDogCents int64            `json:"AdditionalDogCents"`
	PuppyCents         int64            `json:"PuppyCents"`
	LongTermCents      int64            `json:"LongTermCents"`
	HolidayCents       int64            `json:"HolidayCents"`
	GroomingCents      int64            `json:"GroomingCents"`
	PickupCents        int64            `json:"PickupCents"`
	Walking60MinCents  int64            `json:"Walking60MinCents"`
}

type SitterPricing struct {
	SitterID string                      `json:"SitterID"`
	Services map[OrderServiceType]ServicePricing `json:"Services"`
}

