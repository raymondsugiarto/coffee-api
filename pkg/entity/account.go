package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

// AccountDto is the wire shape used by the chart-of-accounts CRUD.
//
// `code` is required and unique-per-organization. `name` is a
// human label ("Kas", "Piutang Karyawan"). The
// `(organizationId, code)` uniqueness is enforced at the DB layer
// via a UNIQUE index; the service layer surfaces an error so the
// caller can react.
type AccountDto struct {
	ID             string `json:"id"`
	OrganizationID string `json:"-"`
	Name           string `json:"name"           validate:"required,min=1,max=255"`
	Code           string `json:"code"           validate:"required,min=1,max=64"`
}

func NewAccountDtoFromModel(m *model.Account) *AccountDto {
	if m == nil {
		return nil
	}
	return &AccountDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Name:           m.Name,
		Code:           m.Code,
	}
}

func (d *AccountDto) ToModel() *model.Account {
	m := &model.Account{
		OrganizationID: d.OrganizationID,
		Name:           d.Name,
		Code:           d.Code,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

// AccountFindAllRequest powers GET /api/accounts. Filtering is
// optional; an empty request returns every row visible to the
// caller's organization (NULL organization_id rows are treated
// as global seed).
type AccountFindAllRequest struct {
	FindAllRequest
	Code string
	Name string
}

func (r *AccountFindAllRequest) GenerateFilter() {
	if r.Code != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "code", Op: "eq", Val: r.Code})
	}
	if r.Name != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "name", Op: "ilike", Val: "%" + r.Name + "%"})
	}
}
