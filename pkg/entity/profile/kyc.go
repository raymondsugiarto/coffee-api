package entity

import (
	"mime/multipart"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
)

type UpdateCustomerKYCDto struct {
	Sex                           string                `json:"sex,omitempty"`
	Occupation                    string                `json:"occupation,omitempty"`
	Position                      string                `json:"position,omitempty"`
	SourceOfFunds                 string                `json:"sourceOfFunds,omitempty"`
	AnnualIncome                  string                `json:"annualIncome,omitempty"`
	PurposeOfAccount              string                `json:"purposeOfAccount,omitempty"`
	IdentificationNumber          string                `json:"identificationNumber,omitempty"`
	TaxIdentificationNumber       string                `json:"taxIdentificationNumber,omitempty"`
	Name                          string                `json:"name,omitempty"`
	Address                       string                `json:"address,omitempty"`
	MailingAddress                string                `json:"mailingAddress,omitempty"`
	PhoneNumber                   string                `json:"phoneNumber,omitempty"`
	MobilePhone                   string                `json:"mobilePhone,omitempty"`
	PlaceOfBirth                  string                `json:"placeOfBirth,omitempty"`
	DateOfBirth                   *string               `json:"dateOfBirth,omitempty"`
	CountryOfBirth                string                `json:"countryOfBirth,omitempty"`
	MotherName                    string                `json:"motherName,omitempty"`
	Citizenship                   string                `json:"citizenship,omitempty"`
	MaritalStatus                 string                `json:"maritalStatus,omitempty"`
	BankName                      string                `json:"bankName,omitempty"`
	NameOnBankAccount             string                `json:"nameOnBankAccount,omitempty"`
	BankAccountNumber             string                `json:"bankAccountNumber,omitempty"`
	IdentityCardFile              *multipart.FileHeader `json:"identityCardFile,omitempty"`
	CustomerPhoto                 *multipart.FileHeader `json:"customerPhoto,omitempty"`
	RecipientName                 string                `json:"recipientName,omitempty"`
	RecipientRelationship         string                `json:"recipientRelationship,omitempty"`
	RecipientDateOfBirth          string                `json:"recipientDateOfBirth,omitempty"`
	RecipientCountryOfBirth       string                `json:"recipientCountryOfBirth,omitempty"`
	RecipientIdentificationNumber string                `json:"recipientIdentificationNumber,omitempty"`
	RecipientPhoneNumber          string                `json:"recipientPhoneNumber,omitempty"`
	TaxIdentityCardFile           *multipart.FileHeader `json:"taxIdentityCardFile,omitempty"`
	CityOfBirthID                 string                `json:"cityOfBirthID,omitempty"`
	ProvinceID                    string                `json:"provinceID,omitempty"`
	RegencyID                     string                `json:"regencyID,omitempty"`
	DistrictID                    string                `json:"districtID,omitempty"`
	VillageID                     string                `json:"villageID,omitempty"`
	RT                            string                `json:"rt,omitempty"`
	RW                            string                `json:"rw,omitempty"`
	PostalCode                    string                `json:"postalCode,omitempty"`
	MailingProvinceID             string                `json:"mailingProvinceID,omitempty"`
	MailingRegencyID              string                `json:"mailingRegencyID,omitempty"`
	MailingDistrictID             string                `json:"mailingDistrictID,omitempty"`
	MailingVillageID              string                `json:"mailingVillageID,omitempty"`
	MailingRT                     string                `json:"mailingRT,omitempty"`
	MailingRW                     string                `json:"mailingRW,omitempty"`
	MailingPostalCode             string                `json:"mailingPostalCode,omitempty"`
}

func (i *UpdateCustomerKYCDto) ToDto() *entity.CustomerDto {
	dto := &entity.CustomerDto{
		Sex:                     i.Sex,
		Occupation:              i.Occupation,
		Position:                i.Position,
		SourceOfFunds:           i.SourceOfFunds,
		AnnualIncome:            i.AnnualIncome,
		PurposeOfAccount:        i.PurposeOfAccount,
		IdentificationNumber:    i.IdentificationNumber,
		TaxIdentificationNumber: i.TaxIdentificationNumber,
		FirstName:               i.Name,
		Address:                 i.Address,
		MailingAddress:          i.MailingAddress,
		PhoneNumber:             i.PhoneNumber,
		MobilePhone:             i.MobilePhone,
		PlaceOfBirth:            i.PlaceOfBirth,
		DateOfBirth:             i.DateOfBirth,
		CountryOfBirth:          i.CountryOfBirth,
		MotherName:              i.MotherName,
		Citizenship:             i.Citizenship,
		MaritalStatus:           i.MaritalStatus,
		BankName:                i.BankName,
		NameOnBankAccount:       i.NameOnBankAccount,
		BankAccountNumber:       i.BankAccountNumber,
		CityOfBirthID:           i.CityOfBirthID,
		ProvinceID:              i.ProvinceID,
		RegencyID:               i.RegencyID,
		DistrictID:              i.DistrictID,
		VillageID:               i.VillageID,
		RT:                      i.RT,
		RW:                      i.RW,
		PostalCode:              i.PostalCode,
		MailingProvinceID:       i.MailingProvinceID,
		MailingRegencyID:        i.MailingRegencyID,
		MailingDistrictID:       i.MailingDistrictID,
		MailingVillageID:        i.MailingVillageID,
		MailingRT:               i.MailingRT,
		MailingRW:               i.MailingRW,
		MailingPostalCode:       i.MailingPostalCode,
	}
	dto.PensionBenefitRecipients = []*entity.PensionBenefitRecipientDto{
		{
			Name:                 i.RecipientName,
			Relationship:         i.RecipientRelationship,
			DateOfBirth:          i.RecipientDateOfBirth,
			IdentificationNumber: i.RecipientIdentificationNumber,
			PhoneNumber:          i.RecipientPhoneNumber,
		},
	}

	return dto
}
