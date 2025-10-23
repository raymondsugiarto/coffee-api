package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type TicketApprovalInputDto struct {
	OrganizationID string
	Status         string
}

func (dto *TicketApprovalInputDto) ToDto() *TicketDto {
	return &TicketDto{
		OrganizationID: dto.OrganizationID,
		Status:         model.TicketStatus(dto.Status),
	}
}

type TicketInputDto struct {
	OrganizationID string
	UserID         string
	Title          string
	Message        string
}

func (dto *TicketInputDto) ToDto() *TicketDto {
	return &TicketDto{
		OrganizationID: dto.OrganizationID,
		UserID:         dto.UserID,
		Title:          dto.Title,
		Message:        dto.Message,
		Status:         model.TicketStatusSubmitted,
	}
}

type TicketDto struct {
	ID             string             `json:"id"`
	OrganizationID string             `json:"organizationId"`
	Organization   *OrganizationData  `json:"-"`
	UserID         string             `json:"userId"`
	User           *UserDto           `json:"user,omitempty"`
	Title          string             `json:"title"`
	Message        string             `json:"message"`
	Status         model.TicketStatus `json:"status"`
	Customer       *CustomerDto       `json:"customer"`
	CreatedAt      time.Time          `json:"createdAt"`
}

func (dto *TicketDto) ToApprovalSubmitDto(id string) *ApprovalDto {
	return &ApprovalDto{
		OrganizationID: dto.OrganizationID,
		UserIDRequest:  id,
		RefID:          dto.ID,
		RefTable:       "ticket",
		Detail:         "Ticket: " + dto.Message, // TBD
		Type:           "TICKET",
		Action:         "ADD",
		Status:         "SUBMIT",
		Reason:         "New Ticket",
	}
}

func (dto *TicketDto) ToModel() *model.Ticket {
	m := &model.Ticket{
		OrganizationID: dto.OrganizationID,
		UserID:         dto.UserID,
		Title:          dto.Title,
		Message:        dto.Message,
		Status:         string(dto.Status),
	}
	if dto.ID != "" {
		m.ID = dto.ID
	}
	return m
}

func (dto *TicketDto) FromModel(m *model.Ticket) *TicketDto {
	dto.ID = m.ID
	dto.OrganizationID = m.OrganizationID
	dto.UserID = m.UserID
	dto.Title = m.Title
	dto.Message = m.Message
	dto.Status = model.TicketStatus(m.Status)
	if m.User != nil {
		dto.User = (new(UserDto)).FromModel(m.User)
	}
	dto.CreatedAt = m.CreatedAt
	return dto
}

func (i *TicketDto) GetInfo() RejectEmail {
	return RejectEmail{
		Email:       i.Customer.Email,
		Name:        i.Customer.ID,
		Description: "Pengajuan Tiket Dengan Kode Nasabah " + i.Customer.ID,
	}
}
