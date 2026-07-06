package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type ItemCategoryDto struct {
	ID             string `json:"id"`
	OrganizationID string `json:"-"`
	Name           string `json:"name" validate:"required,min=1,max=255"`
}

func NewItemCategoryDtoFromModel(m *model.ItemCategory) *ItemCategoryDto {
	if m == nil {
		return nil
	}
	return &ItemCategoryDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Name:           m.Name,
	}
}

func (d *ItemCategoryDto) ToModel() *model.ItemCategory {
	m := &model.ItemCategory{
		OrganizationID: d.OrganizationID,
		Name:           d.Name,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

type ItemDto struct {
	ID             string          `json:"id"`
	OrganizationID string          `json:"-"`
	Organization   *OrganizationDto `json:"-"`
	CategoryID     string          `json:"categoryId"`
	Category       *ItemCategoryDto `json:"category,omitempty"`
	CompanyID      string          `json:"companyId"`
	Company        *CompanyDto     `json:"company"`
	Code           string          `json:"code"`
	SKU            string          `json:"sku"`
	Name           string          `json:"name" validate:"required,min=1,max=255"`
	Price          float64         `json:"sellingPrice" validate:"gte=0"`
	CostPrice      float64         `json:"costPrice" validate:"gte=0"`
	IsActive       bool            `json:"isActive"`
}

func NewItemDtoFromModel(m *model.Item) *ItemDto {
	if m == nil {
		return nil
	}
	d := &ItemDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		CategoryID:     m.CategoryID,
		Code:           m.Code,
		SKU:            m.SKU,
		Name:           m.Name,
		Price:          m.Price,
		CostPrice:      m.CostPrice,
		IsActive:       m.IsActive,
	}
	if m.Category != nil {
		d.Category = NewItemCategoryDtoFromModel(m.Category)
	}
	return d
}

func (d *ItemDto) ToModel() *model.Item {
	m := &model.Item{
		OrganizationID: d.OrganizationID,
		CategoryID:     d.CategoryID,
		Code:           d.Code,
		SKU:            d.SKU,
		Name:           d.Name,
		Price:          d.Price,
		CostPrice:      d.CostPrice,
		IsActive:       d.IsActive,
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
	CategoryID     string
	IsActive       *bool
	Query          string
	MyEmployeeItem bool
}

func (r *ItemFindAllRequest) GenerateFilter() {
	if r.CategoryID != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "category_id", Op: "eq", Val: r.CategoryID})
	}
	if r.IsActive != nil {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "is_active", Op: "eq", Val: *r.IsActive})
	}
}
