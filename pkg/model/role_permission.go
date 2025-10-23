package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type RolePermission struct {
	concern.CommonWithIDs
	RoleID       string
	PermissionID string

	Role       Role
	Permission Permission
}
