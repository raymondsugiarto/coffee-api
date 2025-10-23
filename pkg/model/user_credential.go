package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type UserCredential struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   Organization
	UserID         string
	User           User
	Username       string
	Password       string
}
