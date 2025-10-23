package entity

import "github.com/raymondsugiarto/coffee-api/pkg/model"

type UserHasRoleDto struct {
	UserID string
	User   UserDto
	RoleID string
	Role   RoleDto
}

func (u *UserHasRoleDto) ToModel() *model.UserHasRole {
	return &model.UserHasRole{
		UserID: u.UserID,
		RoleID: u.RoleID,
	}
}
