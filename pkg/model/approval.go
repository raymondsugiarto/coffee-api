package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type Approval struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	UserIDRequest  string
	UserRequest    *User `gorm:"foreignKey:UserIDRequest"`
	RefID          string
	RefTable       string
	Detail         string
	Type           ApprovalType
	Action         ApprovalAction
	Status         string
	Reason         string
	CreatedDate    time.Time
}

type ApprovalType string

const (
	ApprovalTypeCompany              ApprovalType = "COMPANY"
	ApprovalTypeCustomer             ApprovalType = "CUSTOMER"
	ApprovalTypeClaim                ApprovalType = "CLAIM"
	ApprovalTypeInvestment           ApprovalType = "INVESTMENT"
	ApprovalTicketInvestment         ApprovalType = "TICKET"
	ApprovalTypeBenefitParticipation ApprovalType = "BENEFIT_PARTICIPATION"
)

type ApprovalAction string

const (
	ApprovalActionAdd       ApprovalAction = "ADD"
	ApprovalActionUpdate    ApprovalAction = "UPDATE"
	ApprovalActionDelete    ApprovalAction = "DELETE"
	ApprovalActionSuspend   ApprovalAction = "SUSPEND"
	ApprovalActionUnSuspend ApprovalAction = "UNSUSPEND"
)
