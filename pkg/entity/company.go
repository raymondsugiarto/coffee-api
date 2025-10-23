package entity

import (
	"mime/multipart"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type CompanyInputDto struct {
	CompanyType          string                `json:"companyType"`
	FirstName            string                `json:"firstName"`
	Address              string                `json:"address"`
	Email                string                `json:"email"`
	PhoneNumber          string                `json:"phoneNumber"`
	NIB                  string                `json:"nib"`
	Domisili             string                `json:"domisili"`
	NPWP                 string                `json:"npwp"`
	PicName              string                `json:"picName"`
	PicPhone             string                `json:"picPhone"`
	PicEmail             string                `json:"picEmail"`
	AgreementFee         float64               `json:"agreementFee"`
	CooperationAgreement *multipart.FileHeader `json:"cooperationAgreement"`
	NPWPPerusahaan       *multipart.FileHeader `json:"npwpPerusahaan"`
	PilarType            string                `json:"pilarType"`
}

func (u *CompanyInputDto) ToDto() *CompanyDto {
	return &CompanyDto{
		CompanyType:  u.CompanyType,
		FirstName:    u.FirstName,
		Address:      u.Address,
		Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
		NIB:          u.NIB,
		Domisili:     u.Domisili,
		NPWP:         u.NPWP,
		PicName:      u.PicName,
		PicPhone:     u.PicPhone,
		PicEmail:     u.PicEmail,
		AgreementFee: u.AgreementFee,
		PilarType:    u.PilarType,
	}
}

type CompanyDto struct {
	ID                   string    `json:"id"`
	UserID               string    `json:"userId"`
	OrganizationID       string    `json:"organizationId"`
	CompanyCode          string    `json:"companyCode"`
	PhoneNumber          string    `json:"phoneNumber"`
	Email                string    `json:"email"`
	FirstName            string    `json:"firstName"`
	LastName             string    `json:"lastName"`
	CompanyType          string    `json:"companyType"`
	Address              string    `json:"address"`
	NIB                  string    `json:"nib"`
	Domisili             string    `json:"domisili"`
	NPWP                 string    `json:"npwp"`
	PicName              string    `json:"picName"`
	PicPhone             string    `json:"picPhone"`
	PicEmail             string    `json:"picEmail"`
	AgreementFee         float64   `json:"agreementFee"`
	CooperationAgreement string    `json:"cooperationAgreement"`
	AktaPerusahaan       string    `json:"aktaPerusahaan"`
	NIBFile              string    `json:"nibFile"`
	TDP                  string    `json:"tdp"`
	KTP                  string    `json:"ktp"`
	NPWPPerusahaan       string    `json:"npwpPerusahaan"`
	SuratKuasa           string    `json:"suratKuasa"`
	Status               string    `json:"status"`
	PilarType            string    `json:"pilarType"`
	CreatedAt            time.Time `json:"createdAt"`

	User           *UserDto    `json:"user"`
	DomisiliObject *RegencyDto `json:"domisiliObject"`

	Approval *ApprovalDto `json:"approval,omitempty"`
}

func (c *CompanyDto) FromModel(m *model.Company) *CompanyDto {
	c.ID = m.ID
	c.UserID = m.UserID
	c.OrganizationID = m.OrganizationID
	c.CompanyCode = m.CompanyCode
	c.PhoneNumber = m.PhoneNumber
	c.Email = m.Email
	c.FirstName = m.FirstName
	c.LastName = m.LastName
	c.CompanyType = string(m.CompanyType)
	c.Address = m.Address
	c.NIB = m.NIB
	c.Domisili = m.Domisili
	c.NPWP = m.NPWP
	c.PicName = m.PicName
	c.PicPhone = m.PicPhone
	c.PicEmail = m.PicEmail
	c.AgreementFee = m.AgreementFee
	c.CooperationAgreement = m.CooperationAgreement
	c.AktaPerusahaan = m.AktaPerusahaan
	c.NIBFile = m.NIBFile
	c.TDP = m.TDP
	c.KTP = m.KTP
	c.NPWPPerusahaan = m.NPWPPerusahaan
	c.SuratKuasa = m.SuratKuasa
	c.Status = m.Status
	c.User = &UserDto{
		ID:       m.User.ID,
		UserType: UserType(m.User.UserType),
	}
	if m.DomisiliObject != nil {
		c.DomisiliObject = (&RegencyDto{}).FromModel(m.DomisiliObject)
	}
	c.PilarType = string(m.PilarType)
	c.CreatedAt = m.CreatedAt
	return c
}

func (c *CompanyDto) ToModel() *model.Company {
	m := &model.Company{
		UserID:         c.UserID,
		OrganizationID: c.OrganizationID,
		CompanyCode:    c.CompanyCode,
		PhoneNumber:    c.PhoneNumber,
		Email:          c.Email,
		FirstName:      c.FirstName,
		LastName:       c.LastName,
		CompanyType:    model.CompanyType(c.CompanyType),
		Address:        c.Address,
		NIB:            c.NIB,
		Domisili:       c.Domisili,
		NPWP:           c.NPWP,
		PicName:        c.PicName,
		PicPhone:       c.PicPhone,
		PicEmail:       c.PicEmail,
		AgreementFee:   c.AgreementFee,
		PilarType:      model.PilarType(c.PilarType),
	}
	if c.CooperationAgreement != "" {
		m.CooperationAgreement = c.CooperationAgreement
	}
	if c.AktaPerusahaan != "" {
		m.AktaPerusahaan = c.AktaPerusahaan
	}
	if c.NIBFile != "" {
		m.NIBFile = c.NIBFile
	}
	if c.TDP != "" {
		m.TDP = c.TDP
	}
	if c.KTP != "" {
		m.KTP = c.KTP
	}
	if c.NPWPPerusahaan != "" {
		m.NPWPPerusahaan = c.NPWPPerusahaan
	}
	if c.SuratKuasa != "" {
		m.SuratKuasa = c.SuratKuasa
	}
	if c.ID != "" {
		m.ID = c.ID
	}
	if c.User != nil {
		m.User = *c.User.ToModel()
	}
	if c.Status != "" {
		m.Status = c.Status
	}
	return m
}

func (c *CompanyDto) ToAdminDto() *AdminDto {
	return &AdminDto{
		AdminType:   "COMPANY",
		UserID:      c.UserID,
		PhoneNumber: c.PhoneNumber,
		Email:       c.Email,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		CompanyID:   c.ID,
	}
}

func (c *CompanyDto) ToSubmitApprovalDto(uid string) *ApprovalDto {
	return &ApprovalDto{
		OrganizationID: c.OrganizationID,
		UserIDRequest:  uid,
		RefID:          c.ID,
		RefTable:       "company",
		Detail:         "Pendaftaran Perusahaan Baru [" + c.FirstName + "]",
		Type:           "COMPANY",
		Action:         "ADD",
		Status:         "SUBMIT",
		Reason:         "New Registration",
	}
}

type CountCompanyPerType struct {
	CompanyType string `json:"companyType"`
	Count       int64  `json:"count"`
}

func (i *CompanyDto) GetInfo() RejectEmail {
	return RejectEmail{
		Email:       i.Email,
		Name:        i.CompanyCode,
		Description: "Pengajuan Perusahaan " + i.FirstName + "Dengan Kode Perusahaan " + i.CompanyCode,
	}
}
