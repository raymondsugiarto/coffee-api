package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type ParticipantDto struct {
	ID             string                  `json:"id"`
	OrganizationID string                  `json:"-"`
	CustomerID     string                  `json:"customerId"`
	Customer       *CustomerDto            `json:"customer"`
	Code           string                  `json:"code"`
	Type           model.InvestmentType    `json:"type"`
	Status         model.ParticipantStatus `json:"status"`
	Balance        float64                 `json:"balance"`
}

func (dto *ParticipantDto) ToModel() *model.Participant {
	m := &model.Participant{
		OrganizationID: dto.OrganizationID,
		CustomerID:     dto.CustomerID,
		Code:           dto.Code,
		Type:           dto.Type,
		Status:         dto.Status,
	}
	if dto.ID != "" {
		m.ID = dto.ID
	}
	return m
}

func (dto *ParticipantDto) FromModel(m *model.Participant) *ParticipantDto {
	dto.ID = m.ID
	dto.OrganizationID = m.OrganizationID
	dto.CustomerID = m.CustomerID
	dto.Code = m.Code
	dto.Type = m.Type
	dto.Status = m.Status
	if m.Customer != nil {
		dto.Customer = new(CustomerDto).FromModel(m.Customer)
	}
	return dto
}

type ParticipantFindIDRequest struct {
	FindAllRequest
	Include string
}

type ParticipantFindAllRequest struct {
	FindAllRequest
	CompanyID    *string
	CustomerID   string
	Status       model.ParticipantStatus
	InvestmentAt time.Time
	CalculateAll bool
	PaidEmployee bool
	UsePeriod    bool // if true, use the month and year of InvestmentAt, otherwise use the full date
}

func (r *ParticipantFindAllRequest) GenerateFilter() {
	if r.CustomerID != "" {
		r.AddFilter(pagination.FilterItem{
			Field: "customer_id",
			Op:    "eq",
			Val:   r.CustomerID,
		})
	}
	if r.Status != "" {
		r.AddFilter(pagination.FilterItem{
			Field: "participant.status",
			Op:    "eq",
			Val:   r.Status,
		})
	}
}
