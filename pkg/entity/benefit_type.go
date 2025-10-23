package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type BenefitTypeDto struct {
	ID                      string    `json:"id"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	MinimumTimePeriodMonths int       `json:"minimumTimePeriodMonths,omitempty"`
	MinimumContribution     float64   `json:"minimumContribution"`
	Status                  string    `json:"status"`
	CreatedAt               time.Time `json:"createdAt"`
	UpdatedAt               time.Time `json:"updatedAt"`
}

func (dto *BenefitTypeDto) FromModel(m *model.BenefitType) *BenefitTypeDto {
	return &BenefitTypeDto{
		ID:                      m.ID,
		Name:                    m.Name,
		Description:             m.Description,
		MinimumTimePeriodMonths: m.MinimumTimePeriodMonths,
		MinimumContribution:     m.MinimumContribution,
		Status:                  string(m.Status),
		CreatedAt:               m.CreatedAt,
		UpdatedAt:               m.UpdatedAt,
	}
}

func (dto *BenefitTypeDto) ToModel() *model.BenefitType {
	m := &model.BenefitType{
		Name:        dto.Name,
		Description: dto.Description,

		MinimumContribution: dto.MinimumContribution,
		Status:              model.BenefitTypeStatus(dto.Status),
	}
	if dto.ID != "" {
		m.ID = dto.ID
	}
	return m
}

type CreateBenefitTypeRequest struct {
	Name                string  `json:"name" validate:"required"`
	Description         string  `json:"description"`
	MinimumContribution float64 `json:"minimumContribution" validate:"required,min=0"`
	Status              string  `json:"status" validate:"required,oneof=ACTIVE INACTIVE"`
}

func (req *CreateBenefitTypeRequest) ToDto() *BenefitTypeDto {
	return &BenefitTypeDto{
		Name:                req.Name,
		Description:         req.Description,
		MinimumContribution: req.MinimumContribution,
		Status:              req.Status,
	}
}

type UpdateBenefitTypeRequest struct {
	Name                    string  `json:"name"`
	Description             string  `json:"description"`
	MinimumTimePeriodMonths int     `json:"minimumTimePeriodMonths" validate:"min=1"`
	MinimumContribution     float64 `json:"minimumContribution" validate:"min=0"`
	Status                  string  `json:"status" validate:"oneof=ACTIVE INACTIVE"`
}

func (req *UpdateBenefitTypeRequest) ToDto() *BenefitTypeDto {
	return &BenefitTypeDto{
		Name:                    req.Name,
		Description:             req.Description,
		MinimumTimePeriodMonths: req.MinimumTimePeriodMonths,
		MinimumContribution:     req.MinimumContribution,
		Status:                  req.Status,
	}
}

type BenefitTypeFilter struct {
	Name   string `json:"name" query:"name"`
	Status string `json:"status" query:"status"`
	Limit  int    `json:"limit" query:"limit"`
	Offset int    `json:"offset" query:"offset"`
}

type BenefitTypeFindAllRequest struct {
	FindAllRequest
	Status string `json:"status" query:"status"`
}

func (r *BenefitTypeFindAllRequest) GenerateFilter() {
	if r.Status != "" {
		r.FindAllRequest.AddFilter(pagination.FilterItem{
			Field: "status",
			Op:    "eq",
			Val:   r.Status,
		})
	}
}
