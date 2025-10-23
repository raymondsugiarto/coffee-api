package entity

import (
	"context"
	"mime/multipart"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/utils"
)

type SignUpCustomerInputDto struct {
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	Username     string `json:"username"`
	ReferralCode string `json:"referralCode"`
}

func (u *SignUpCustomerInputDto) ToDto() *SignUpCustomerDto {
	return &SignUpCustomerDto{
		Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
		Name:         u.Name,
		Password:     u.Password,
		Username:     u.Username,
		ReferralCode: u.ReferralCode,
	}
}

type SignUpCustomerDto struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	Username     string `json:"username"`
	ReferralCode string `json:"referralCode"`
}

func (d *SignUpCustomerDto) ToCustomerDto(ctx context.Context) *entity.CustomerDto {
	organizationID := shared.GetOrganization(ctx).ID
	hashPassword, _ := utils.HashPassword(d.Password)
	dto := &entity.CustomerDto{
		OrganizationID: organizationID,
		Email:          d.Email,
		PhoneNumber:    d.PhoneNumber,
		FirstName:      d.Name,
	}
	dto.User = &entity.UserDto{
		OrganizationID: organizationID,
		UserType:       entity.CUSTOMER,
		UserCredential: []entity.UserCredentialDto{
			{
				OrganizationID: organizationID,
				Username:       d.Username,
				Password:       hashPassword,
			},
		},
	}
	return dto
}

type SignUpCompanyInputDto struct {
	PhoneNumber    string                `json:"phoneNumber"`
	Email          string                `json:"email"`
	Password       string                `json:"password"`
	CompanyType    string                `json:"companyType"`
	Name           string                `json:"name"`
	Address        string                `json:"address"`
	NIB            string                `json:"nib"`
	Domisili       string                `json:"domisili"`
	NPWP           string                `json:"npwp"`
	PicName        string                `json:"picName"`
	PicPhone       string                `json:"picPhone"`
	PicEmail       string                `json:"picEmail"`
	AktaPerusahaan *multipart.FileHeader `json:"aktaPerusahaan"`
	NIBFile        *multipart.FileHeader `json:"nibFile"`
	TDP            *multipart.FileHeader `json:"tdp"`
	KTP            *multipart.FileHeader `json:"ktp"`
	NPWPPerusahaan *multipart.FileHeader `json:"npwpPerusahaan"`
	SuratKuasa     *multipart.FileHeader `json:"suratKuasa"`
}

func (u *SignUpCompanyInputDto) ToDto() *SignUpCompanyDto {
	return &SignUpCompanyDto{
		PhoneNumber: u.PhoneNumber,
		Email:       u.Email,
		Password:    u.Password,
		CompanyType: u.CompanyType,
		Name:        u.Name,
		Address:     u.Address,
		NIB:         u.NIB,
		Domisili:    u.Domisili,
		NPWP:        u.NPWP,
		PicName:     u.PicName,
		PicPhone:    u.PicPhone,
		PicEmail:    u.PicEmail,
	}
}

type SignUpCompanyDto struct {
	ID                   string `json:"id"`
	CompanyCode          string `json:"companyCode"`
	PhoneNumber          string `json:"phoneNumber"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	CompanyType          string `json:"companyType"`
	Name                 string `json:"name"`
	Address              string `json:"address"`
	NIB                  string `json:"nib"`
	Domisili             string `json:"domisili"`
	NPWP                 string `json:"npwp"`
	PicName              string `json:"picName"`
	PicPhone             string `json:"picPhone"`
	PicEmail             string `json:"picEmail"`
	CooperationAgreement string `json:"cooperationAgreement"`
	AktaPerusahaan       string `json:"aktaPerusahaan"`
	NIBFile              string `json:"nibFile"`
	TDP                  string `json:"tdp"`
	KTP                  string `json:"ktp"`
	NPWPPerusahaan       string `json:"npwpPerusahaan"`
	SuratKuasa           string `json:"suratKuasa"`
}

func (d *SignUpCompanyDto) ToCompanyDto(ctx context.Context) *entity.CompanyDto {
	organizationID := shared.GetOrganization(ctx).ID
	hashPassword, _ := utils.HashPassword(d.Password)
	dto := &entity.CompanyDto{
		OrganizationID:       organizationID,
		CompanyCode:          d.CompanyCode,
		Email:                d.Email,
		PhoneNumber:          d.PhoneNumber,
		FirstName:            d.Name,
		CompanyType:          d.CompanyType,
		Address:              d.Address,
		NIB:                  d.NIB,
		Domisili:             d.Domisili,
		NPWP:                 d.NPWP,
		PicName:              d.PicName,
		PicPhone:             d.PicPhone,
		PicEmail:             d.PicEmail,
		CooperationAgreement: d.CooperationAgreement,
		AktaPerusahaan:       d.AktaPerusahaan,
		NIBFile:              d.NIBFile,
		TDP:                  d.TDP,
		KTP:                  d.KTP,
		NPWPPerusahaan:       d.NPWPPerusahaan,
		SuratKuasa:           d.SuratKuasa,
	}
	dto.User = &entity.UserDto{
		OrganizationID: organizationID,
		UserType:       entity.COMPANY,
		UserCredential: []entity.UserCredentialDto{
			{
				OrganizationID: organizationID,
				Username:       d.Email,
				Password:       hashPassword,
			},
		},
	}
	return dto
}
