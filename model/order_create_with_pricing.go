package model

import "demo/util/common"

type CreateOrderWithPricingReq struct {
	UserID          string        `json:"userId"`
	OwnerName       string        `json:"ownerName"`
	SitterID        string        `json:"SitterID"`
	SitterName      string        `json:"sitterName"`
	Contact         string        `json:"Contact"`
	AlternativeContact string     `json:"AlternativeContact"`
	Note            string        `json:"Note"`
	Type            int64         `json:"type"`
	FromDate        common.Date   `json:"FromDate"`
	ToDate          common.Date   `json:"ToDate"`
	PetList         []*OrderPet   `json:"PetList"`
	NeedGrooming    bool          `json:"NeedGrooming"`
	NeedPickup      bool          `json:"NeedPickup"`
	DurationMinutes int           `json:"DurationMinutes"`
}
