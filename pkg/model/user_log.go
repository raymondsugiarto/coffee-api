package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type UserLog struct {
	concern.CommonWithIDs
	OrganizationID   string
	CompanyID        string
	Company          *Company
	UserCredentialID string
	UserCredential   *UserCredential
	RefModule        string
	RefTable         string
	RefID            string
	Description      string
}
