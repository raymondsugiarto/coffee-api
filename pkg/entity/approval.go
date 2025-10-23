package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type ApprovalDto struct {
	ID             string               `json:"id,omitempty"`
	OrganizationID string               `json:"organizationId,omitempty"`
	UserIDRequest  string               `json:"userIDRequest,omitempty"`
	RefID          string               `json:"refID,omitempty"`
	RefTable       string               `json:"refTable,omitempty"`
	Detail         string               `json:"detail,omitempty"`
	Type           model.ApprovalType   `json:"type,omitempty"`
	Action         model.ApprovalAction `json:"action,omitempty"`
	Status         string               `json:"status,omitempty"`
	Reason         string               `json:"reason,omitempty"`
	CreatedAt      time.Time            `json:"createdAt,omitempty"`
	CreatedDate    time.Time            `json:"createdDate,omitempty"`
	UpdatedAt      time.Time            `json:"updatedAt,omitempty"` // This field is optional and can be used to track when the approval was last updated

	RefData     interface{} `json:"refData,omitempty"` // This can be used to hold additional data related to the approval
	UserRequest *UserDto    `json:"userRequest,omitempty"`
}

func (a *ApprovalDto) FromModel(m *model.Approval) *ApprovalDto {
	a.ID = m.ID
	a.OrganizationID = m.OrganizationID
	a.UserIDRequest = m.UserIDRequest
	a.RefID = m.RefID
	a.RefTable = m.RefTable
	a.Detail = m.Detail
	a.Type = m.Type
	a.Action = m.Action
	a.Status = m.Status
	a.Reason = m.Reason
	a.CreatedAt = m.CreatedAt
	a.UpdatedAt = m.UpdatedAt
	a.CreatedDate = m.CreatedDate
	if m.UserRequest != nil {
		a.UserRequest = (&UserDto{}).FromModel(m.UserRequest)
	}
	return a
}

func (a *ApprovalDto) ToModel() *model.Approval {
	m := &model.Approval{
		OrganizationID: a.OrganizationID,
		UserIDRequest:  a.UserIDRequest,
		RefID:          a.RefID,
		RefTable:       a.RefTable,
		Detail:         a.Detail,
		Type:           a.Type,
		Action:         a.Action,
		Status:         a.Status,
		Reason:         a.Reason,
		CreatedDate:    a.CreatedDate,
	}
	if a.ID != "" {
		m.ID = a.ID
	}
	return m
}

type ApprovalFindAllRequest struct {
	FindAllRequest
	Status    string
	Type      model.ApprovalType
	Types     []model.ApprovalType
	Statuses  []string
	CompanyID string
}

func (r *ApprovalFindAllRequest) GenerateFilter() {
	if r.Status != "" {
		r.FindAllRequest.AddFilter(pagination.FilterItem{
			Field: "status",
			Op:    "eq",
			Val:   r.Status,
		})
	}

	if r.Type != "" {
		r.FindAllRequest.AddFilter(pagination.FilterItem{
			Field: "type",
			Op:    "eq",
			Val:   r.Type,
		})
	}
}

type RejectEmail struct {
	Email       string
	Name        string
	Description string
}

type HasContactInfo interface {
	GetInfo() RejectEmail
}
