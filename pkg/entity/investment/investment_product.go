package entity

import (
	"mime/multipart"

	e "github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type InvestmentProductInputDto struct {
	Code           string                `json:"code"`
	Name           string                `json:"name"`
	Description    string                `json:"description"`
	FixIncome      float64               `json:"fixIncome"`
	StockValue     float64               `json:"stockValue"`
	MixedValue     float64               `json:"mixedValue"`
	MoneyMarket    float64               `json:"moneyMarket"`
	ShariaValue    float64               `json:"shariaValue"`
	FundFactSheet  *multipart.FileHeader `json:"fundFactSheet"`
	Riplay         *multipart.FileHeader `json:"riplay"`
	AdminFee       float64               `json:"adminFee"`
	ManagementFee  float64               `json:"managementFee"`
	FounderFee     float64               `json:"founderFee"`
	CommissionFee  float64               `json:"commissionFee"`
	OrderingNumber int                   `json:"orderingNumber"`
}

func (i *InvestmentProductInputDto) ToDto() *InvestmentProductDto {
	return &InvestmentProductDto{
		Code:           i.Code,
		Name:           i.Name,
		Description:    i.Description,
		FixIncome:      i.FixIncome,
		StockValue:     i.StockValue,
		MixedValue:     i.MixedValue,
		MoneyMarket:    i.MoneyMarket,
		ShariaValue:    i.ShariaValue,
		AdminFee:       i.AdminFee,
		ManagementFee:  i.ManagementFee,
		FounderFee:     i.FounderFee,
		CommissionFee:  i.CommissionFee,
		OrderingNumber: i.OrderingNumber,
	}
}

func (*InvestmentProductDto) TableName() string {
	return "investment_product"
}

type InvestmentProductDto struct {
	ID             string  `json:"id,omitempty"`
	OrganizationID string  `json:"-"`
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	FixIncome      float64 `json:"fixIncome"`
	StockValue     float64 `json:"stockValue"`
	MixedValue     float64 `json:"mixedValue"`
	MoneyMarket    float64 `json:"moneyMarket"`
	ShariaValue    float64 `json:"shariaValue"`
	FundFactSheet  string  `json:"fundFactSheet"`
	Riplay         string  `json:"riplay"`
	AdminFee       float64 `json:"adminFee"`
	ManagementFee  float64 `json:"managementFee"`
	FounderFee     float64 `json:"founderFee"`
	CommissionFee  float64 `json:"commissionFee"`
	Status         string  `json:"status"`
	OrderingNumber int     `json:"orderingNumber"`
	AUM            float64 `json:"aum,omitempty" gorm:"-"`

	NetAssetValueDto *NetAssetValueDto `json:"netAssetValue,omitempty" gorm:"-"`
}

func (d *InvestmentProductDto) ToModel() *model.InvestmentProduct {
	m := &model.InvestmentProduct{
		OrganizationID: d.OrganizationID,
		Code:           d.Code,
		Name:           d.Name,
		Description:    d.Description,
		FixIncome:      d.FixIncome,
		StockValue:     d.StockValue,
		MixedValue:     d.MixedValue,
		MoneyMarket:    d.MoneyMarket,
		ShariaValue:    d.ShariaValue,
		AdminFee:       d.AdminFee,
		ManagementFee:  d.ManagementFee,
		FounderFee:     d.FounderFee,
		CommissionFee:  d.CommissionFee,
		Status:         d.Status,
		OrderingNumber: d.OrderingNumber,
	}
	if d.FundFactSheet != "" {
		m.FundFactSheet = d.FundFactSheet
	}
	if d.Riplay != "" {
		m.Riplay = d.Riplay
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *InvestmentProductDto) FromModel(m *model.InvestmentProduct) *InvestmentProductDto {
	d.ID = m.ID
	d.OrganizationID = m.OrganizationID
	d.Code = m.Code
	d.Name = m.Name
	d.Description = m.Description
	d.FixIncome = m.FixIncome
	d.StockValue = m.StockValue
	d.MixedValue = m.MixedValue
	d.MoneyMarket = m.MoneyMarket
	d.ShariaValue = m.ShariaValue
	d.FundFactSheet = m.FundFactSheet
	d.Riplay = m.Riplay
	d.AdminFee = m.AdminFee
	d.ManagementFee = m.ManagementFee
	d.FounderFee = m.FounderFee
	d.CommissionFee = m.CommissionFee
	d.Status = m.Status
	d.OrderingNumber = m.OrderingNumber
	return d
}

type InvestmentProductSummaryDto struct {
	InvestmentProductID   string  `json:"investmentProductId"`
	InvestmentProductName string  `json:"investmentProductName"`
	InvestAmount          float64 `json:"investAmount"`
	TotalAmount           float64 `json:"totalAmount"`
	Ip                    float64 `json:"ip"`
	Nab                   float64 `json:"nab"`
}

type InvestmentProductFilter struct {
	e.FindAllRequest
	IncludeAum bool
}
