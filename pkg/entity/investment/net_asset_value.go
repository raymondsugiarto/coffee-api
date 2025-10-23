package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type NetAssetValueBatchInputDto struct {
	Items []*NetAssetValueInputDto `json:"items"`
}

func (i *NetAssetValueBatchInputDto) ToDto() *NetAssetValueBatchDto {
	items := make([]*NetAssetValueDto, 0)
	for _, item := range i.Items {
		items = append(items, item.ToDto())
	}
	return &NetAssetValueBatchDto{
		Items: items,
	}
}

type NetAssetValueInputDto struct {
	InvestmentProductID string  `json:"investmentProductId,omitempty"`
	Amount              float64 `json:"amount,omitempty"`
}

func (i *NetAssetValueInputDto) ToDto() *NetAssetValueDto {
	return &NetAssetValueDto{
		InvestmentProductID: i.InvestmentProductID,
		Amount:              i.Amount,
		CreatedDate:         time.Now().Local().UTC(),
	}
}

type NetAssetValueBatchDto struct {
	Items []*NetAssetValueDto `json:"items"`
}

func (d *NetAssetValueBatchDto) ToModel() []*model.NetAssetValue {
	items := make([]*model.NetAssetValue, 0)
	for _, item := range d.Items {
		items = append(items, item.ToModel())
	}
	return items
}

type NetAssetValueDto struct {
	ID                  string                `json:"id"`
	OrganizationID      string                `json:"organizationId,omitempty"`
	InvestmentProductID string                `json:"investmentProductId,omitempty"`
	InvestmentProduct   *InvestmentProductDto `json:"investmentProduct,omitempty"`
	Amount              float64               `json:"amount,omitempty"`
	CreatedDate         time.Time             `json:"createdDate,omitempty"`
}

func (d *NetAssetValueDto) ToModel() *model.NetAssetValue {
	m := &model.NetAssetValue{
		OrganizationID:      d.OrganizationID,
		InvestmentProductID: d.InvestmentProductID,
		Amount:              d.Amount,
		CreatedDate:         d.CreatedDate,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *NetAssetValueDto) FromModel(m *model.NetAssetValue) *NetAssetValueDto {
	d.ID = m.ID
	d.OrganizationID = m.OrganizationID
	d.InvestmentProductID = m.InvestmentProductID
	d.Amount = m.Amount
	d.CreatedDate = m.CreatedDate

	if m.InvestmentProduct != nil {
		d.InvestmentProduct = new(InvestmentProductDto).FromModel(m.InvestmentProduct)
	}

	return d
}

type NetAssetValueFindAllRequest struct {
	entity.FindAllRequest
	FormCreate bool
}
