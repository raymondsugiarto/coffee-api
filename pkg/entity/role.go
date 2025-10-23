package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type RoleInputDto struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	RoleIDParent *string `json:"roleIdParent"`
}

func (r *RoleInputDto) ToDto() *RoleDto {
	return &RoleDto{
		ID:           r.ID,
		Name:         r.Name,
		RoleIDParent: r.RoleIDParent,
	}
}

type RoleDto struct {
	ID             string   `json:"id,omitempty"`
	OrganizationID string   `json:"-"`
	Name           string   `json:"name,omitempty"`
	RoleIDParent   *string  `json:"roleIdParent"`
	RoleParent     *RoleDto `json:"roleParent,omitempty"`
}

func (r *RoleDto) FromModel(m *model.Role) *RoleDto {
	r.ID = m.ID
	r.Name = m.Name
	r.RoleIDParent = m.RoleIDParent

	if m.RoleParent != nil {
		r.RoleParent = new(RoleDto).FromModel(m.RoleParent)
	}
	return r
}

func (r *RoleDto) ToModel() *model.Role {
	m := &model.Role{
		OrganizationID: r.OrganizationID,
		Name:           r.Name,
		RoleIDParent:   r.RoleIDParent,
	}
	if r.ID != "" {
		m.ID = r.ID
	}
	return m
}
