package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type Permission struct {
	concern.CommonWithIDs
	Code            string
	RolePermissions []RolePermission
}
