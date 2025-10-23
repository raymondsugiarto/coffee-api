package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type BankCustomer struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	CustomerID     string
	Customer       *Customer
	BankCode       string
	BankName       string
	AccountName    string
	AccountNumber  string
	IsDefault      bool
}
