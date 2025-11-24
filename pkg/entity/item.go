package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type ItemDto struct {
	ID             string           `json:"id"`
	OrganizationID string           `json:"-"`
	Organization   *OrganizationDto `json:"-"`
	CompanyID      string           `json:"companyId"`
	Company        *CompanyDto      `json:"company"`
	Code           string           `json:"code"`
	Name           string           `json:"name"`
	Price          float64          `json:"price"`
}

func NewItemDtoFromModel(m *model.Item) *ItemDto {
	if m == nil {
		return nil
	}
	return &ItemDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Code:           m.Code,
		Name:           m.Name,
		Price:          m.Price,
	}
}

func (d *ItemDto) ToModel() *model.Item {
	m := &model.Item{
		OrganizationID: d.OrganizationID,
		Code:           d.Code,
		Name:           d.Name,
		Price:          d.Price,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

type ItemFindAllRequest struct {
	FindAllRequest
	UserID         string
	CompanyID      string
	MyEmployeeItem bool
}

func (r *ItemFindAllRequest) GenerateFilter() {
}
