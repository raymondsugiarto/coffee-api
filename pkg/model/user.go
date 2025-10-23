package model

import (
	"errors"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type User struct {
	concern.CommonWithIDs
	OrganizationID          string
	Organization            Organization
	UserType                UserType
	PhoneVerificationStatus IdentityStatus
	EmailVerificationStatus IdentityStatus
	UserCredential          []*UserCredential
	UserHasRole             []*UserHasRole
}

type UserType string

const (
	COMPANY  UserType = "COMPANY"
	EMPLOYEE UserType = "EMPLOYEE"
	ADMIN    UserType = "ADMIN"
)

// Function to convert string to UserType
func StringToUserType(s string) (UserType, error) {
	switch s {
	case "admin":
		return ADMIN, nil
	case "company":
		return COMPANY, nil
	case "employee":
		return EMPLOYEE, nil
	default:
		return "", errors.New("enumParseError")
	}
}

type IdentityStatus string

const (
	UNVERIFIED IdentityStatus = "UNVERIFIED"
	VERIFIED   IdentityStatus = "VERIFIED"
)
