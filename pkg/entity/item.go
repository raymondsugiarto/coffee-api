package entity

import (
	"strings"

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
	ID             string           `json:"id"`
	OrganizationID string           `json:"-"`
	Organization   *OrganizationDto `json:"-"`
	CategoryID     string           `json:"categoryId"`
	Category       *ItemCategoryDto `json:"category,omitempty"`
	ParentID       string           `json:"parentId"`
	Parent         *ItemDto         `json:"parent,omitempty"`
	CompanyID      string           `json:"companyId"`
	Company        *CompanyDto      `json:"company"`
	Code           string           `json:"code"`
	SKU            string           `json:"sku"`
	Name           string           `json:"name" validate:"required,min=1,max=255"`
	Price          float64          `json:"price" validate:"gte=0"`
	CostPrice      float64          `json:"costPrice" validate:"gte=0"`
	Commision      float64          `json:"commision" validate:"gte=0"`
	IsActive       bool             `json:"isActive"`
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
		CategoryID:     d.CategoryID,
		ParentID:       d.ParentID,
		Code:           d.Code,
		SKU:            d.SKU,
		Name:           d.Name,
		Price:          d.Price,
		CostPrice:      d.CostPrice,
		Commision:      d.Commision,
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
	AdminID        string `json:"adminId"`
	CompanyID      string
	MyEmployeeItem bool
	Session        model.SessionType
	CategoryID     string
	IsActive       *bool
	Query          string
	ParentID       string   // <-- "parent_id" filter (exact). Empty = top-level.
	ParentIDs      []string // <-- "parent_id IN (...)" filter. Empty = no restriction.
}

func (r *ItemFindAllRequest) GenerateFilter() {
	if r.CategoryID != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "category_id", Op: "eq", Val: r.CategoryID})
	}
	if r.IsActive != nil {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "is_active", Op: "eq", Val: *r.IsActive})
	}
	if r.ParentID != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "parent_id", Op: "eq", Val: r.ParentID})
	}
}

type ItemCategoryFindAllRequest struct {
	FindAllRequest
	Query string
}

func (r *ItemCategoryFindAllRequest) GenerateFilter() {
	if q := strings.TrimSpace(r.Query); q != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "name", Op: "ilike", Val: "%" + q + "%"})
	}
}
