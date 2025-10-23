package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type Company struct {
	concern.CommonWithIDs
	OrganizationID       string
	Organization         Organization
	UserID               string
	User                 User
	CompanyCode          string
	PhoneNumber          string
	FirstName            string
	LastName             string
	Email                string
	CompanyType          CompanyType
	Address              string
	NIB                  string
	Domisili             string
	NPWP                 string
	PicName              string
	PicPhone             string
	PicEmail             string
	AgreementFee         float64
	CooperationAgreement string
	AktaPerusahaan       string
	NIBFile              string
	TDP                  string
	KTP                  string
	NPWPPerusahaan       string
	SuratKuasa           string
	Status               string
	DomisiliObject       *Regency `gorm:"foreignKey:Domisili"`
	PilarType            PilarType
}

type CompanyType string

const (
	CompanyTypeDKP  CompanyType = "DKP"
	CompanyTypePPIP CompanyType = "PPIP"
)

type PilarType string

const (
	TypePilar    PilarType = "PILAR"
	TypeNonPilar PilarType = "NON PILAR"
)
