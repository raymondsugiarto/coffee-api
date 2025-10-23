package entity

import (
	"mime/multipart"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type ClaimInputDto struct {
	ParticipantID           string                `json:"participantId"`
	BankName                string                `json:"bankName"`
	AccountName             string                `json:"accountName"`
	AccountNumber           string                `json:"accountNumber"`
	BankBranch              string                `json:"bankBranch"`
	Amount                  float64               `json:"amount"`
	ParticipantCard         *multipart.FileHeader `json:"participantCard"`
	IdentityCardFile        *multipart.FileHeader `json:"identityCardFile"`
	TaxIdentityFile         *multipart.FileHeader `json:"taxIdentityFile"`
	FamilyCardFile          *multipart.FileHeader `json:"familyCardFile"`
	DeathCertificateFile    *multipart.FileHeader `json:"deathCertificateFile"`
	GuardianCertificateFile *multipart.FileHeader `json:"guardianCertificateFile"`
	MedicalCertificateFile  *multipart.FileHeader `json:"medicalCertificateFile"`
	WorkCertificateFile     *multipart.FileHeader `json:"workCertificateFile"`
}

func (i *ClaimInputDto) ToDto() *ClaimDto {
	return &ClaimDto{
		ParticipantID: i.ParticipantID,
		BankName:      i.BankName,
		AccountName:   i.AccountName,
		AccountNumber: i.AccountNumber,
		BankBranch:    i.BankBranch,
		Amount:        i.Amount,
	}
}

type ClaimDto struct {
	ID                      string               `json:"id"`
	OrganizationID          string               `json:"organizationID"`
	ParticipantID           string               `json:"participantId"`
	Participant             *ParticipantDto      `json:"participant"`
	BankName                string               `json:"bankName"`
	AccountName             string               `json:"accountName"`
	AccountNumber           string               `json:"accountNumber"`
	BankBranch              string               `json:"bankBranch"`
	Amount                  float64              `json:"amount"`
	ApprovalStatus          model.ApprovalStatus `json:"approvalStatus"`
	ParticipantCard         string               `json:"participantCard"`
	IdentityCardFile        string               `json:"identityCardFile"`
	TaxIdentityFile         string               `json:"taxIdentityFile"`
	FamilyCardFile          string               `json:"familyCardFile"`
	DeathCertificateFile    string               `json:"deathCertificateFile"`
	GuardianCertificateFile string               `json:"guardianCertificateFile"`
	MedicalCertificateFile  string               `json:"medicalCertificateFile"`
	WorkCertificateFile     string               `json:"workCertificateFile"`
	CreatedAt               time.Time            `json:"createdAt"`
	UpdatedAt               time.Time            `json:"updatedAt"`
}

func (dto *ClaimDto) ToModel() *model.Claim {
	m := &model.Claim{
		OrganizationID:          dto.OrganizationID,
		ParticipantID:           dto.ParticipantID,
		BankName:                dto.BankName,
		AccountName:             dto.AccountName,
		AccountNumber:           dto.AccountNumber,
		BankBranch:              dto.BankBranch,
		Amount:                  dto.Amount,
		ParticipantCard:         dto.ParticipantCard,
		IdentityCardFile:        dto.IdentityCardFile,
		TaxIdentityFile:         dto.TaxIdentityFile,
		FamilyCardFile:          dto.FamilyCardFile,
		DeathCertificateFile:    dto.DeathCertificateFile,
		GuardianCertificateFile: dto.GuardianCertificateFile,
		MedicalCertificateFile:  dto.MedicalCertificateFile,
		WorkCertificateFile:     dto.WorkCertificateFile,
		ApprovalStatus:          dto.ApprovalStatus,
	}
	if dto.ID != "" {
		m.ID = dto.ID
	}
	return m
}

func (dto *ClaimDto) FromModel(m *model.Claim) *ClaimDto {
	dto.ID = m.ID
	dto.OrganizationID = m.OrganizationID
	dto.ParticipantID = m.ParticipantID
	dto.BankName = m.BankName
	dto.AccountName = m.AccountName
	dto.AccountNumber = m.AccountNumber
	dto.BankBranch = m.BankBranch
	dto.Amount = m.Amount
	dto.ParticipantCard = m.ParticipantCard
	dto.IdentityCardFile = m.IdentityCardFile
	dto.TaxIdentityFile = m.TaxIdentityFile
	dto.FamilyCardFile = m.FamilyCardFile
	dto.DeathCertificateFile = m.DeathCertificateFile
	dto.GuardianCertificateFile = m.GuardianCertificateFile
	dto.MedicalCertificateFile = m.MedicalCertificateFile
	dto.WorkCertificateFile = m.WorkCertificateFile
	dto.ApprovalStatus = m.ApprovalStatus
	dto.CreatedAt = m.CreatedAt
	dto.UpdatedAt = m.UpdatedAt
	if m.Participant != nil {
		dto.Participant = (&ParticipantDto{}).FromModel(m.Participant)
	}
	return dto
}

func (dto *ClaimDto) ToApprovalSubmitDto(uid string) *ApprovalDto {
	return &ApprovalDto{
		OrganizationID: dto.OrganizationID,
		UserIDRequest:  uid,
		RefID:          dto.ID,
		RefTable:       "claim",
		Detail:         "Pengajuan Manfaat [" + dto.Participant.Code + " - " + dto.AccountName + "]",
		Type:           "CLAIM",
		Action:         "ADD",
		Status:         "SUBMIT",
		Reason:         "New Claim",
	}
}

type ClaimFindAllRequest struct {
	FindAllRequest
	CustomerID     string
	CompanyID      *string
	ParticipantID  string
	ApprovalStatus model.ApprovalStatus
}

func (r *ClaimFindAllRequest) GenerateFilter() {
	if r.ParticipantID != "" {
		r.AddFilter(pagination.FilterItem{
			Field: "participant_id",
			Op:    "eq",
			Val:   r.ParticipantID,
		})
	}
	if r.ApprovalStatus != "" {
		r.AddFilter(pagination.FilterItem{
			Field: "approval_status",
			Op:    "eq",
			Val:   r.ApprovalStatus,
		})
	}
}

func (i *ClaimDto) GetInfo() RejectEmail {
	return RejectEmail{
		Email:       i.Participant.Customer.Email,
		Name:        i.Participant.Code,
		Description: "Pengajuan Manfaat Dengan Kode " + i.Participant.Code,
	}
}
