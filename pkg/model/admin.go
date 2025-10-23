package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type Admin struct {
	concern.CommonWithIDs
	AdminType       string
	UserID          string
	User            *User
	PhoneNumber     string
	Email           string
	FirstName       string
	LastName        string
	ProfileImageUrl string
	OrganizationID  string
	CompanyID       string
	Company         *Company
}
