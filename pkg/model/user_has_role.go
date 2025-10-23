package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type UserHasRole struct {
	concern.CommonWithIDs
	UserID string
	User   User
	RoleID string
	Role   Role
}
