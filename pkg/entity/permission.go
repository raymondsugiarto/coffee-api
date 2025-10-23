package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type PermissionDto struct {
	ID             string               `json:"id"`
	Code           string               `json:"code"`
	RolePermission []*RolePermissionDto `json:"rolePermission"`
}

func (p *PermissionDto) FromModel(m *model.Permission) *PermissionDto {
	p.ID = m.ID
	p.Code = m.Code

	if len(m.RolePermissions) > 0 {
		p.RolePermission = make([]*RolePermissionDto, 0)
		for _, rp := range m.RolePermissions {
			p.RolePermission = append(p.RolePermission, new(RolePermissionDto).FromModel(&rp))
		}
	}
	return p
}

type PermissionFindAllRequest struct {
	FindAllRequest
	UserID        string
	RoleID        string
	AllPermission bool
}

func (p *PermissionFindAllRequest) GenerateFilter() {
}

type RolePermissionDto struct {
	ID           string         `json:"id"`
	RoleID       string         `json:"roleId"`
	PermissionID string         `json:"permissionId"`
	Role         *RoleDto       `json:"role"`
	Permission   *PermissionDto `json:"permission"`
}

func (r *RolePermissionDto) FromModel(m *model.RolePermission) *RolePermissionDto {
	r.ID = m.ID
	r.RoleID = m.RoleID
	r.PermissionID = m.PermissionID
	r.Role = new(RoleDto).FromModel(&m.Role)
	r.Permission = new(PermissionDto).FromModel(&m.Permission)
	return r
}

func (r *RolePermissionDto) ToModel() *model.RolePermission {
	m := &model.RolePermission{
		RoleID:       r.RoleID,
		PermissionID: r.PermissionID,
	}
	if r.ID != "" {
		m.ID = r.ID
	}

	return m
}
