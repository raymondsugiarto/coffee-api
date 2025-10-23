package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

const (
	POINT_REFERRAL float64 = 1.0
)

type CustomerPointDto struct {
	concern.CommonWithIDs
	OrganizationID string            `json:"organizationId"`
	Organization   *OrganizationData `json:"-"`
	CustomerID     string            `json:"customerId"`
	Customer       *CustomerDto      `json:"customer"`
	Point          float64           `json:"point"`
	Direction      string            `json:"direction"`
	Description    string            `json:"description"`
	RefID          string            `json:"refId"`
	RefCode        string            `json:"refCode"`
	RefModule      string            `json:"refModule"`
}

func (dto *CustomerPointDto) ToModel() *model.CustomerPoint {
	r := &model.CustomerPoint{
		OrganizationID: dto.OrganizationID,
		CustomerID:     dto.CustomerID,
		Point:          dto.Point,
		Direction:      dto.Direction,
		Description:    dto.Description,
		RefID:          dto.RefID,
		RefCode:        dto.RefCode,
		RefModule:      dto.RefModule,
	}
	if dto.ID != "" {
		r.ID = dto.ID
	}
	return r
}

func (dto *CustomerPointDto) FromModel(m *model.CustomerPoint) *CustomerPointDto {
	dto.ID = m.ID
	dto.OrganizationID = m.OrganizationID
	dto.CustomerID = m.CustomerID
	dto.Point = m.Point
	dto.Direction = m.Direction
	dto.Description = m.Description
	dto.RefID = m.RefID
	dto.RefCode = m.RefCode
	dto.RefModule = m.RefModule
	return dto
}
