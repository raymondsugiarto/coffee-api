package entity

import (
	"strings"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type CompanyDto struct {
	ID             string          `json:"id"`
	OrganizationID string          `json:"-"`
	Organization   OrganizationDto `json:"-"`
	PhoneNumber    string          `json:"phoneNumber"`
	Name           string          `json:"name"`
}

func NewCompanyDtoFromModel(m *model.Company) *CompanyDto {
	if m == nil {
		return nil
	}
	d := &CompanyDto{
		ID:   m.ID,
		Name: m.Name,
	}
	return d
}

// CompanyFindAllRequest powers GET /api/companies. Mirrors the
// ItemCategoryFindAllRequest pattern: org-scope + free-text query
// on name. The handler fills OrganizationData from the request
// context, the service layer passes it through to the repository,
// and the repository applies the WHERE clause.
type CompanyFindAllRequest struct {
	FindAllRequest
	Query string
}

func (r *CompanyFindAllRequest) GenerateFilter() {
	if q := strings.TrimSpace(r.Query); q != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "name", Op: "ilike", Val: "%" + q + "%"})
	}
}
