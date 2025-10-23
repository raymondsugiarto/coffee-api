package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type Role struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	Name           string
	RoleIDParent   *string
	RoleParent     *Role `gorm:"foreignKey:RoleIDParent"`

	RolePermissions []RolePermission
}
