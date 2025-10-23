package entity

import (
	"mime/multipart"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type RewardInputDto struct {
	Name        string                `json:"name" form:"name" validate:"required"`
	Points      int                   `json:"points" form:"points" validate:"required"`
	Stock       int                   `json:"stock" form:"stock"`
	ExpiredAt   time.Time             `json:"expiredAt" form:"expiredAt"`
	Image       *multipart.FileHeader `json:"image,omitempty" form:"image"`
	Status      string                `json:"status" form:"status"`
	Description string                `json:"description" form:"description" validate:"required"`
	Code        string                `json:"code" form:"code" validate:"required"`
}

func (d *RewardInputDto) ToDto() *RewardDto {
	return &RewardDto{
		Name:        d.Name,
		Points:      d.Points,
		Stock:       d.Stock,
		ExpiredAt:   d.ExpiredAt,
		Status:      d.Status,
		Description: d.Description,
		Code:        d.Code,
	}
}

type RewardDto struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Points      int       `json:"points"`
	Stock       int       `json:"stock"`
	ImageUrl    string    `json:"imageUrl"`
	ExpiredAt   time.Time `json:"expiredAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
}

func (dto *RewardDto) ToModel() *model.Reward {
	return &model.Reward{
		CommonWithIDs: concern.CommonWithIDs{
			ID: dto.ID,
		},
		Name:        dto.Name,
		Points:      dto.Points,
		Stock:       dto.Stock,
		ExpiredAt:   dto.ExpiredAt,
		ImageUrl:    dto.ImageUrl,
		Status:      model.RewardStatus(dto.Status),
		Description: dto.Description,
		Code:        dto.Code,
	}
}

func (dto *RewardDto) FromModel(m *model.Reward) *RewardDto {
	dto.ID = m.ID
	dto.Name = m.Name
	dto.Points = m.Points
	dto.Stock = m.Stock
	dto.ImageUrl = m.ImageUrl
	dto.ExpiredAt = m.ExpiredAt
	dto.CreatedAt = m.CreatedAt
	dto.UpdatedAt = m.UpdatedAt
	dto.Status = string(m.Status)
	dto.Description = m.Description
	dto.Code = m.Code
	return dto
}
