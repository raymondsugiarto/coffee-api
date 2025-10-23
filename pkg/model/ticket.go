package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type Ticket struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	UserID         string
	User           *User
	Title          string
	Message        string
	Status         string
}

type TicketStatus string

const (
	TicketStatusPending   TicketStatus = "PENDING"
	TicketStatusApproved  TicketStatus = "APPROVED"
	TicketStatusCompleted TicketStatus = "COMPLETED"
	TicketStatusRejected  TicketStatus = "REJECTED"
	TicketStatusCancelled TicketStatus = "CANCELLED"
	TicketStatusSubmitted TicketStatus = "SUBMIT"
)
