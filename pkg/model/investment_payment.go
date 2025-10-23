package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type InvestmentPayment struct {
	concern.CommonWithIDs
	OrganizationID       string
	Organization         *Organization
	InvestmentID         string
	Investment           *Investment
	PaymentMethod        string
	BankID               string
	Bank                 *Bank
	BankCode             string
	BankName             string
	AccountName          string
	AccountNumber        string
	Amount               float64
	ConfirmationImageUrl string
	PaymentAt            time.Time
	Status               InvestmentPaymentStatus
}

type InvestmentPaymentStatus string

const (
	InvestmentPaymentStatusPending   InvestmentPaymentStatus = "pending"
	InvestmentPaymentStatusConfirmed InvestmentPaymentStatus = "confirmed"
	InvestmentPaymentStatusRejected  InvestmentPaymentStatus = "rejected"
	InvestmentPaymentStatusSuccess   InvestmentPaymentStatus = "success"
	InvestmentPaymentStatusExpired   InvestmentPaymentStatus = "expired"
)
