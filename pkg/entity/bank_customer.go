package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type BankCustomerInputDto struct {
	BankCode      string
	BankName      string
	AccountName   string
	AccountNumber string
	CustomerID    string
	IsDefault     bool
}

func (d *BankCustomerInputDto) ToDto() *BankCustomerDto {
	return &BankCustomerDto{
		BankCode:      d.BankCode,
		BankName:      d.BankName,
		AccountName:   d.AccountName,
		AccountNumber: d.AccountNumber,
		CustomerID:    d.CustomerID,
		IsDefault:     d.IsDefault,
	}
}

type BankCustomerDto struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	CustomerID     string `json:"customerId"`
	BankCode       string `json:"bankCode"`
	BankName       string `json:"bankName"`
	AccountName    string `json:"accountName"`
	AccountNumber  string `json:"accountNumber"`
	IsDefault      bool   `json:"isDefault"`
}

func (d *BankCustomerDto) ToModel() *model.BankCustomer {
	m := &model.BankCustomer{
		OrganizationID: d.OrganizationID,
		CustomerID:     d.CustomerID,
		BankCode:       d.BankCode,
		BankName:       d.BankName,
		AccountName:    d.AccountName,
		AccountNumber:  d.AccountNumber,
		IsDefault:      d.IsDefault,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *BankCustomerDto) FromModel(m *model.BankCustomer) *BankCustomerDto {
	return &BankCustomerDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		CustomerID:     m.CustomerID,
		BankCode:       m.BankCode,
		BankName:       m.BankName,
		AccountName:    m.AccountName,
		AccountNumber:  m.AccountNumber,
		IsDefault:      m.IsDefault,
	}
}

type BankCustomerFindAllRequest struct {
	FindAllRequest
	CustomerID string
	IsDefault  bool
}

func (r *BankCustomerFindAllRequest) GenerateFilter() {
	if r.CustomerID != "" {
		r.FindAllRequest.AddFilter(pagination.FilterItem{
			Field: "customer_id",
			Op:    "eq",
			Val:   r.CustomerID,
		})
	}
	if r.IsDefault {
		r.FindAllRequest.AddFilter(pagination.FilterItem{
			Field: "is_default",
			Op:    "eq",
			Val:   r.IsDefault,
		})
	}
}
