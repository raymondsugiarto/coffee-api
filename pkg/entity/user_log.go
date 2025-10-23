package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type UserLogDto struct {
	ID               string          `json:"id"`
	OrganizationID   string          `json:"organizationId"`
	CompanyID        string          `json:"companyId"`
	Company          *CompanyDto     `json:"company,omitempty"`
	UserCredentialID string          `json:"userCredential_id"`
	UserCredential   *UserCredential `json:"userCredential,omitempty"`
	RefModule        string          `json:"refModule"`
	RefTable         string          `json:"refTable"`
	RefID            string          `json:"refId"`
	Description      string          `json:"description"`
	CreatedAt        time.Time       `json:"createdAt"`
}

func (dto *UserLogDto) FromModel(m *model.UserLog) *UserLogDto {
	if m == nil {
		return nil
	}
	dto.ID = m.ID
	dto.OrganizationID = m.OrganizationID
	dto.CompanyID = m.CompanyID
	if m.Company != nil {
		dto.Company = &CompanyDto{}
		dto.Company.FromModel(m.Company)
	}
	dto.UserCredentialID = m.UserCredentialID
	if m.UserCredential != nil {
		dto.UserCredential = &UserCredential{}
		dto.UserCredential.Username = m.UserCredential.Username
	}
	dto.RefModule = m.RefModule
	dto.RefTable = m.RefTable
	dto.RefID = m.RefID
	dto.Description = m.Description
	dto.CreatedAt = m.CreatedAt
	return dto
}

func (dto *UserLogDto) ToModel() *model.UserLog {
	m := &model.UserLog{
		OrganizationID:   dto.OrganizationID,
		CompanyID:        dto.CompanyID,
		UserCredentialID: dto.UserCredentialID,
		RefModule:        dto.RefModule,
		RefTable:         dto.RefTable,
		RefID:            dto.RefID,
		Description:      dto.Description,
	}
	if dto.ID != "" {
		m.ID = dto.ID
	}
	return m
}

type UserLogFindAllRequest struct {
	FindAllRequest
	CompanyID *string
	RefModule *string
}

func (dto *UserLogFindAllRequest) GenerateFilter() {
	if dto.CompanyID != nil {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "company_id",
				Op:    "eq",
				Val:   dto.CompanyID,
			},
		)
	}
	if dto.RefModule != nil {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "ref_module",
				Op:    "eq",
				Val:   dto.RefModule,
			},
		)
	}
}
