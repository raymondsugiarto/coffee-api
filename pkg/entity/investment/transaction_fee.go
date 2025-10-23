package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type TransactionFeeDto struct {
	concern.CommonWithIDs
	OrganizationID      string
	CompanyID           string
	Company             *entity.CompanyDto
	ParticipantID       string
	Participant         *entity.ParticipantDto
	InvestmentProductID string
	InvestmentProduct   *InvestmentProductDto
	Type                model.InvestmentType
	TransactionDate     time.Time
	OperationFee        float64
	Nav                 float64
	Ip                  float64
	PortfolioAmount     float64
}

func (tf *TransactionFeeDto) ToModel() *model.TransactionFee {
	m := &model.TransactionFee{
		OrganizationID:      tf.OrganizationID,
		CompanyID:           tf.CompanyID,
		Company:             tf.Company.ToModel(),
		ParticipantID:       tf.ParticipantID,
		Participant:         tf.Participant.ToModel(),
		InvestmentProductID: tf.InvestmentProductID,
		InvestmentProduct:   tf.InvestmentProduct.ToModel(),
		Type:                tf.Type,
		TransactionDate:     tf.TransactionDate,
		OperationFee:        tf.OperationFee,
		Nav:                 tf.Nav,
		Ip:                  tf.Ip,
		PortfolioAmount:     tf.PortfolioAmount,
	}
	if tf.ID != "" {
		m.ID = tf.ID
	}

	return m
}

func (tf *TransactionFeeDto) FromModel(m *model.TransactionFee) *TransactionFeeDto {
	t := &TransactionFeeDto{
		CommonWithIDs: concern.CommonWithIDs{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: m.DeletedAt,
		},
		OrganizationID:      m.OrganizationID,
		CompanyID:           m.CompanyID,
		ParticipantID:       m.ParticipantID,
		InvestmentProductID: m.InvestmentProductID,
		Type:                m.Type,
		TransactionDate:     m.TransactionDate,
		OperationFee:        m.OperationFee,
		Nav:                 m.Nav,
		Ip:                  m.Ip,
		PortfolioAmount:     m.PortfolioAmount,
	}
	if m.Company != nil {
		t.Company = new(entity.CompanyDto).FromModel(m.Company)
	}
	if m.Participant != nil {
		t.Participant = new(entity.ParticipantDto).FromModel(m.Participant)
	}
	if m.InvestmentProduct != nil {
		t.InvestmentProduct = new(InvestmentProductDto).FromModel(m.InvestmentProduct)
	}
	return t
}
