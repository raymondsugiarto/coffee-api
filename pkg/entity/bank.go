package entity

import "github.com/raymondsugiarto/coffee-api/pkg/model"

type BankInputDto struct {
	BankCode      string
	BankName      string
	AccountName   string
	AccountNumber string
}

func (d *BankInputDto) ToDto() *BankDto {
	return &BankDto{
		BankCode:      d.BankCode,
		BankName:      d.BankName,
		AccountName:   d.AccountName,
		AccountNumber: d.AccountNumber,
	}
}

type BankDto struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	BankCode       string `json:"bankCode"`
	BankName       string `json:"bankName"`
	AccountName    string `json:"accountName"`
	AccountNumber  string `json:"accountNumber"`
}

func (d *BankDto) ToModel() *model.Bank {
	m := &model.Bank{
		OrganizationID: d.OrganizationID,
		BankCode:       d.BankCode,
		BankName:       d.BankName,
		AccountName:    d.AccountName,
		AccountNumber:  d.AccountNumber,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *BankDto) FromModel(m *model.Bank) *BankDto {
	return &BankDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		BankCode:       m.BankCode,
		BankName:       m.BankName,
		AccountName:    m.AccountName,
		AccountNumber:  m.AccountNumber,
	}
}
