package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type BenefitParticipationDto struct {
	ID             string                           `json:"id"`
	OrganizationID string                           `json:"organizationId,omitempty"`
	CustomerID     string                           `json:"customerId"`
	Customer       *CustomerDto                     `json:"customer,omitempty"`
	Participant    *ParticipantDto                  `json:"participant,omitempty"`
	ParticipantID  string                           `json:"participantId"`
	Status         model.BenefitParticipationStatus `json:"status"`

	ExternalDplkName                *string  `json:"externalDplkName,omitempty"`
	ExternalDplkParticipantNumber   *string  `json:"externalDplkParticipantNumber,omitempty"`
	ExternalDplkMonthlyContribution *float64 `json:"externalDplkMonthlyContribution,omitempty"`
	HasBpjsPensionProgram           *bool    `json:"hasBpjsPensionProgram,omitempty"`

	Details      []*BenefitParticipationDetailDto `json:"details,omitempty"`
	InvestmentID string                           `json:"investmentId,omitempty"`
	CreatedAt    time.Time                        `json:"createdAt"`
	UpdatedAt    time.Time                        `json:"updatedAt"`
}

func (dto *BenefitParticipationDto) FromModel(m *model.BenefitParticipation) *BenefitParticipationDto {
	result := &BenefitParticipationDto{
		ID:                              m.ID,
		CustomerID:                      m.CustomerID,
		ParticipantID:                   m.ParticipantID,
		Status:                          m.Status,
		ExternalDplkName:                m.ExternalDplkName,
		ExternalDplkParticipantNumber:   m.ExternalDplkParticipantNumber,
		ExternalDplkMonthlyContribution: m.ExternalDplkMonthlyContribution,
		HasBpjsPensionProgram:           m.HasBpjsPensionProgram,
		InvestmentID:                    m.InvestmentID,
		CreatedAt:                       m.CreatedAt,
		UpdatedAt:                       m.UpdatedAt,
	}

	if len(m.Details) > 0 {
		result.Details = make([]*BenefitParticipationDetailDto, 0)
		for _, detail := range m.Details {
			result.Details = append(result.Details, new(BenefitParticipationDetailDto).FromModel(detail))
		}
	}

	return result
}

func (dto *BenefitParticipationDto) ToModel() *model.BenefitParticipation {
	m := &model.BenefitParticipation{
		CommonWithIDs: concern.CommonWithIDs{
			ID: dto.ID,
		},
		CustomerID:                      dto.CustomerID,
		ParticipantID:                   dto.ParticipantID,
		Status:                          dto.Status,
		ExternalDplkName:                dto.ExternalDplkName,
		ExternalDplkParticipantNumber:   dto.ExternalDplkParticipantNumber,
		ExternalDplkMonthlyContribution: dto.ExternalDplkMonthlyContribution,
		HasBpjsPensionProgram:           dto.HasBpjsPensionProgram,
		InvestmentID:                    dto.InvestmentID,
		Details:                         make([]*model.BenefitParticipationDetail, 0),
	}

	if len(dto.Details) > 0 {
		for _, detail := range dto.Details {
			m.Details = append(m.Details, detail.ToModel())
		}
	}

	return m
}

type BenefitParticipationDetailDto struct {
	ID                      string                           `json:"id"`
	BenefitParticipationID  string                           `json:"benefitParticipationId"`
	BenefitTypeID           string                           `json:"benefitTypeId"`
	BenefitType             *BenefitTypeDto                  `json:"benefitType,omitempty"`
	TimePeriodMonths        int                              `json:"timePeriodMonths"`
	PlannedWithdrawalMonths int                              `json:"plannedWithdrawalMonths"`
	MonthlyContribution     float64                          `json:"monthlyContribution"`
	Status                  model.BenefitParticipationStatus `json:"status"`
	CreatedAt               time.Time                        `json:"createdAt"`
	UpdatedAt               time.Time                        `json:"updatedAt"`
}

func (dto *BenefitParticipationDetailDto) FromModel(m *model.BenefitParticipationDetail) *BenefitParticipationDetailDto {
	return &BenefitParticipationDetailDto{
		ID:                      m.ID,
		BenefitParticipationID:  m.BenefitParticipationID,
		BenefitTypeID:           m.BenefitTypeID,
		BenefitType:             new(BenefitTypeDto).FromModel(&m.BenefitType),
		TimePeriodMonths:        m.TimePeriodMonths,
		PlannedWithdrawalMonths: m.PlannedWithdrawalMonths,
		MonthlyContribution:     m.MonthlyContribution,
		Status:                  m.Status,
		CreatedAt:               m.CreatedAt,
		UpdatedAt:               m.UpdatedAt,
	}
}

func (dto *BenefitParticipationDetailDto) ToModel() *model.BenefitParticipationDetail {
	return &model.BenefitParticipationDetail{
		CommonWithIDs: concern.CommonWithIDs{
			ID: dto.ID,
		},
		BenefitParticipationID:  dto.BenefitParticipationID,
		BenefitTypeID:           dto.BenefitTypeID,
		TimePeriodMonths:        dto.TimePeriodMonths,
		PlannedWithdrawalMonths: dto.PlannedWithdrawalMonths,
		MonthlyContribution:     dto.MonthlyContribution,
		Status:                  dto.Status,
	}
}

func (i *BenefitParticipationDto) GetInfo() RejectEmail {
	return RejectEmail{
		Email:       i.Customer.Email,
		Name:        i.Participant.Code,
		Description: "Pengajuan Manfaat Peserta Dengan Kode " + i.Participant.Code,
	}
}
