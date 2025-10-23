package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type Claim struct {
	concern.CommonWithIDs
	OrganizationID          string
	ParticipantID           string
	Participant             *Participant
	BankName                string
	AccountName             string
	AccountNumber           string
	BankBranch              string
	Amount                  float64
	ApprovalStatus          ApprovalStatus
	ParticipantCard         string
	IdentityCardFile        string
	TaxIdentityFile         string
	FamilyCardFile          string
	DeathCertificateFile    string
	GuardianCertificateFile string
	MedicalCertificateFile  string
	WorkCertificateFile     string
}
