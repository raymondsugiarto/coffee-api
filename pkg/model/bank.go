package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type Bank struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	BankCode       string
	BankName       string
	AccountName    string
	AccountNumber  string
}
