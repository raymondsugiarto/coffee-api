package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type UserIdentityVerification struct {
	concern.CommonWithIDs
	UserID         string
	User           User
	OrganizationID string
	Organization   Organization
	IdentityFor    IdentityFor
	IdentityType   string
	UserIdentity   string
	UniqueCode     string
	TryCount       int
	ExpiredAt      time.Time
	Status         UserIdentityVerificationStatus
}

type UserIdentityVerificationStatus string

const (
	USER_IDENTITY_VERIFICATION_REQUEST UserType = "REQUEST"
	USER_IDENTITY_VERIFICATION_SUCCESS UserType = "SUCCESS"
)

type IdentityFor string

const (
	EMAIL_VERIFICATION       IdentityFor = "EMAIL_VERIFICATION"
	PHONE_VERIFICATION       IdentityFor = "PHONE_VERIFICATION"
	FORGOT_PASSWORD_CUSTOMER IdentityFor = "FORGOT_PASSWORD_CUSTOMER"
	FORGOT_PASSWORD_COMPANY  IdentityFor = "FORGOT_PASSWORD_COMPANY"
)
