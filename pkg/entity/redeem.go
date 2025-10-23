package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity/concern"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type RedeemInputDto struct {
	OrganizationID string `json:"organizationId"`
	CustomerID     string `json:"customerId"`
	RewardID       string `json:"rewardId" validate:"required"`
	Address        string `json:"address" validate:"required"`
	ProvinceID     string `json:"provinceId" validate:"required"`
	RegencyID      string `json:"regencyId" validate:"required"`
	DistrictID     string `json:"districtId" validate:"required"`
	VillageID      string `json:"villageId" validate:"required"`
	RT             string `json:"rt"`
	RW             string `json:"rw"`
	PostalCode     string `json:"postalCode"`
	Note           string `json:"note"`
}

func (input *RedeemInputDto) ToDto() *RedeemDto {
	return &RedeemDto{
		OrganizationID: input.OrganizationID,
		CustomerID:     input.CustomerID,
		RewardID:       input.RewardID,
		Address:        input.Address,
		ProvinceID:     input.ProvinceID,
		RegencyID:      input.RegencyID,
		DistrictID:     input.DistrictID,
		VillageID:      input.VillageID,
		RT:             input.RT,
		RW:             input.RW,
		PostalCode:     input.PostalCode,
		Note:           input.Note,
		Status:         model.RedeemPending,
	}
}

type RedeemDto struct {
	concern.CommonWithID
	OrganizationID string             `json:"organizationId"`
	CustomerID     string             `json:"customerId"`
	Customer       *CustomerDto       `json:"customer,omitempty"`
	RewardID       string             `json:"rewardId"`
	Reward         *RewardDto         `json:"reward,omitempty"`
	RedemptionDate time.Time          `json:"redemptionDate"`
	PointsRedeemed int                `json:"pointsRedeemed"`
	RedemptionCode string             `json:"redemptionCode"`
	Status         model.RedeemStatus `json:"status"`
	Address        string             `json:"address"`
	ProvinceID     string             `json:"provinceId"`
	Province       *ProvinceDto       `json:"province,omitempty"`
	RegencyID      string             `json:"regencyId"`
	Regency        *RegencyDto        `json:"regency,omitempty"`
	DistrictID     string             `json:"districtId"`
	District       *DistrictDto       `json:"district,omitempty"`
	VillageID      string             `json:"villageId"`
	Village        *VillageDto        `json:"village,omitempty"`
	RT             string             `json:"rt"`
	RW             string             `json:"rw"`
	PostalCode     string             `json:"postalCode"`
	Note           string             `json:"note"`
}

type UpdateRedeemStatusDto struct {
	Status model.RedeemStatus `json:"status" validate:"required,oneof=PENDING REJECTED COMPLETED"`
}

type RedeemFindAllRequest struct {
	FindAllRequest
	CustomerID string `query:"customerId"`
	Status     string `query:"status"`
}

func (dto *RedeemDto) ToModel() *model.Redeem {
	r := &model.Redeem{
		OrganizationID: dto.OrganizationID,
		CustomerID:     dto.CustomerID,
		RewardID:       dto.RewardID,
		RedemptionDate: dto.RedemptionDate,
		PointsRedeemed: dto.PointsRedeemed,
		RedemptionCode: dto.RedemptionCode,
		Status:         dto.Status,
		Address:        dto.Address,
		ProvinceID:     dto.ProvinceID,
		RegencyID:      dto.RegencyID,
		DistrictID:     dto.DistrictID,
		VillageID:      dto.VillageID,
		RT:             dto.RT,
		RW:             dto.RW,
		PostalCode:     dto.PostalCode,
		Note:           dto.Note,
	}
	if dto.ID != "" {
		r.ID = dto.ID
	}
	return r
}

func (dto *RedeemDto) FromModel(m *model.Redeem) *RedeemDto {
	dto.ID = m.ID
	dto.OrganizationID = m.OrganizationID
	dto.CustomerID = m.CustomerID
	dto.RewardID = m.RewardID
	dto.RedemptionDate = m.RedemptionDate
	dto.PointsRedeemed = m.PointsRedeemed
	dto.RedemptionCode = m.RedemptionCode
	dto.Status = m.Status
	dto.Address = m.Address
	dto.ProvinceID = m.ProvinceID
	dto.RegencyID = m.RegencyID
	dto.DistrictID = m.DistrictID
	dto.VillageID = m.VillageID
	dto.RT = m.RT
	dto.RW = m.RW
	dto.PostalCode = m.PostalCode
	dto.Note = m.Note
	dto.CreatedAt = m.CreatedAt
	dto.UpdatedAt = m.UpdatedAt

	if m.Customer != nil {
		dto.Customer = new(CustomerDto).FromModel(m.Customer)
	}

	if m.Reward != nil {
		dto.Reward = new(RewardDto).FromModel(m.Reward)
	}

	if m.Province != nil {
		dto.Province = new(ProvinceDto).FromModel(m.Province)
	}

	if m.Regency != nil {
		dto.Regency = new(RegencyDto).FromModel(m.Regency)
	}

	if m.District != nil {
		dto.District = new(DistrictDto).FromModel(m.District)
	}

	if m.Village != nil {
		dto.Village = new(VillageDto).FromModel(m.Village)
	}

	return dto
}
