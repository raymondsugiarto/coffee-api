package entity

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

const (
	OriginKey       = "x-origin"
	OriginTypeKey   = "x-origin-type" // ADMIN, COMPANY, CUSTOMER
	OrganizationKey = "organization"

	UserContextKey        = "user"
	UserCredentialDataKey = "userCredentialData"

	CompanyKey = "company"
)

type OrganizationData struct {
	ID string `json:"id"`
}

type UserCredentialData struct {
	ID         string `json:"id"`   // user credential id
	UserID     string `json:"uid"`  // user id
	CustomerID string `json:"cid"`  // user id
	AdminID    string `json:"aid"`  // user id
	CompanyID  string `json:"coid"` // user id
}

type OrganizationDto struct {
	concern.CommonWithIDs
	Code   string
	Name   string
	Origin string
}
