package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
	"gorm.io/gorm"
)

type Customer struct {
	concern.CommonWithIDs
	OrganizationID           string
	Organization             *Organization
	UserID                   string
	User                     *User
	CustomerIDParent         string
	CustomerParent           *Customer `gorm:"foreignKey:CustomerIDParent"`
	CompanyID                string
	Company                  *Company
	Email                    string
	PhoneNumber              string
	FirstName                string
	LastName                 string
	ReferralCode             string
	SIMNumber                string
	SIMStatus                SIMStatus      `gorm:"default:INACTIVE"`
	ApprovalStatus           ApprovalStatus `gorm:"default:SUBMIT"`
	CustomerID               string
	Nickname                 string
	DateStart                *string
	PlaceOfBirth             string
	DateOfBirth              *string
	CountryOfBirth           string
	MotherName               string
	NormalRetirementAge      *int64
	Citizenship              string
	Sex                      string
	MaritalStatus            string
	Occupation               string
	Position                 string
	SourceOfFunds            string
	AnnualIncome             string
	PurposeOfOpeningAccount  string
	NameOnBankAccount        string
	BankAccountNumber        string
	BankName                 string
	IdentificationNumber     string
	TaxIdentificationNumber  string
	Address                  string
	MailingAddress           string
	OfficeAddress            string
	PhoneOffice              string
	MobilePhone              string
	EmployerPercentage       *float64
	EmployerAmount           *float64
	CustomerPercentage       *float64
	CustomerAmount           *float64
	EffectiveDate            *string
	PaymentMethod            string
	IdentityCardFile         string
	CustomerPhoto            string
	TaxIdentityCardFile      string
	CityOfBirthID            string
	CityOfBirth              *Regency `gorm:"foreignKey:CityOfBirthID"`
	ProvinceID               string
	Province                 *Province `gorm:"foreignKey:ProvinceID"`
	RegencyID                string
	Regency                  *Regency `gorm:"foreignKey:RegencyID"`
	DistrictID               string
	District                 *District `gorm:"foreignKey:DistrictID"`
	VillageID                string
	Village                  *Village `gorm:"foreignKey:VillageID"`
	RT                       string
	RW                       string
	PostalCode               string
	MailingProvinceID        string
	MailingProvince          *Province `gorm:"foreignKey:MailingProvinceID"`
	MailingRegencyID         string
	MailingRegency           *Regency `gorm:"foreignKey:MailingRegencyID"`
	MailingDistrictID        string
	MailingDistrict          *District `gorm:"foreignKey:MailingDistrictID"`
	MailingVillageID         string
	MailingVillage           *Village `gorm:"foreignKey:MailingVillageID"`
	MailingRT                string
	MailingRW                string
	MailingPostalCode        string
	PensionBenefitRecipients []*PensionBenefitRecipient `gorm:"foreignKey:CustomerID"`
	CountryBirth             *Country                   `gorm:"foreignKey:CountryOfBirth"`
	VoluntaryAmount          *float64
	EducationFundAmount      *float64
	Citizen                  *Country `gorm:"foreignKey:Citizenship"`
}

type SIMStatus string
type ApprovalStatus string
type Sex string

const (
	SIMStatusInactive  SIMStatus = "INACTIVE"
	SIMStatusActive    SIMStatus = "ACTIVE"
	SIMStatusUnsuspend SIMStatus = "UNSUSPEND"
	SIMStatusSuspend   SIMStatus = "SUSPEND"
)

const (
	ApprovalStatusSubmit     ApprovalStatus = "SUBMIT"
	ApprovalStatusKyc        ApprovalStatus = "KYC"
	ApprovalStatusKycPending ApprovalStatus = "KYC_PENDING"
	ApprovalStatusDone       ApprovalStatus = "DONE"
	ApprovalStatusApproved   ApprovalStatus = "APPROVED"
	ApprovalStatusRejected   ApprovalStatus = "REJECTED"
)

const (
	SexMale   Sex = "MALE"
	SexFemale Sex = "FEMALE"
)

func (c *Customer) Complete(db *gorm.DB) *gorm.DB {
	return db.Preload("Company").
		Preload("PensionBenefitRecipients.CountryBirth").
		Preload("CityOfBirth").
		Preload("Province").
		Preload("MailingProvince").
		Preload("Regency").
		Preload("MailingRegency").
		Preload("District").
		Preload("MailingDistrict").
		Preload("Village").
		Preload("MailingVillage").
		Preload("CountryBirth").
		Preload("Citizen")
}
