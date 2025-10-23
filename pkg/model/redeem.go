package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type Redeem struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	CustomerID     string
	Customer       *Customer
	RewardID       string
	Reward         *Reward
	RedemptionDate time.Time
	PointsRedeemed int
	RedemptionCode string
	Status         RedeemStatus
	Address        string
	ProvinceID     string
	Province       *Province `gorm:"foreignKey:ProvinceID"`
	RegencyID      string
	Regency        *Regency `gorm:"foreignKey:RegencyID"`
	DistrictID     string
	District       *District `gorm:"foreignKey:DistrictID"`
	VillageID      string
	Village        *Village `gorm:"foreignKey:VillageID"`
	RT             string
	RW             string
	PostalCode     string
	Note           string
}

type RedeemStatus string

const (
	RedeemPending   RedeemStatus = "PENDING"
	RedeemRejected  RedeemStatus = "REJECTED"
	RedeemCompleted RedeemStatus = "COMPLETED"
)
