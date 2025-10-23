package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type PensionBenefitRecipient struct {
	concern.CommonWithIDs
	Name                 string
	Relationship         string
	DateOfBirth          string
	CountryOfBirth       string
	CountryBirth         *Country `gorm:"foreignKey:CountryOfBirth"`
	IdentificationNumber string
	PhoneNumber          string
	CustomerID           string
}
