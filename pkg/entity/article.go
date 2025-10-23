package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type ArticleStatus string

const (
	DRAFT     ArticleStatus = "DRAFT"
	PUBLISHED ArticleStatus = "PUBLISHED"
	ARCHIVED  ArticleStatus = "ARCHIVED"
)

type ArticleDto struct {
	ID        string        `json:"id"`
	ImageUrl  string        `json:"imageUrl"`
	Title     string        `json:"title"`
	Slug      string        `json:"slug"`
	Content   string        `json:"content"`
	Status    ArticleStatus `json:"status"`
	CreatedBy string        `json:"createdBy"`
	UpdatedBy string        `json:"updatedBy"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

type ArticleInputDto struct {
	ImageUrl string        `json:"imageUrl"`
	Title    string        `json:"title" validate:"required"`
	Slug     string        `json:"slug" validate:"required"`
	Content  string        `json:"content" validate:"required"`
	Status   ArticleStatus `json:"status" validate:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
}

type ArticleUpdateStatusDto struct {
	Status string `json:"status" validate:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
	ID     string `json:"id" validate:"required"`
}

type ArticleFindAllRequest struct {
	FindAllRequest
	Status ArticleStatus `query:"status"`
}

func (r *ArticleFindAllRequest) GenerateFilter() {
	if r.Status != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{
			Field: "status",
			Op:    "eq",
			Val:   r.Status,
		})
	}
}

func (dto *ArticleInputDto) ToDto() *ArticleDto {
	return &ArticleDto{
		ImageUrl: dto.ImageUrl,
		Title:    dto.Title,
		Slug:     dto.Slug,
		Content:  dto.Content,
		Status:   dto.Status,
	}
}

func (dto *ArticleDto) ToModel() *model.Article {
	return &model.Article{
		CommonWithIDs: concern.CommonWithIDs{
			ID: dto.ID,
		},
		ImageUrl:  dto.ImageUrl,
		Title:     dto.Title,
		Slug:      dto.Slug,
		Content:   dto.Content,
		Status:    string(dto.Status),
		CreatedBy: dto.CreatedBy,
		UpdatedBy: dto.UpdatedBy,
	}
}

func (dto *ArticleDto) FromModel(m *model.Article) *ArticleDto {
	return &ArticleDto{
		ID:        m.ID,
		ImageUrl:  m.ImageUrl,
		Title:     m.Title,
		Slug:      m.Slug,
		Content:   m.Content,
		Status:    ArticleStatus(m.Status),
		CreatedBy: m.CreatedBy,
		UpdatedBy: m.UpdatedBy,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (dto *ArticleUpdateStatusDto) ToDto() *ArticleDto {
	return &ArticleDto{
		ID:     dto.ID,
		Status: ArticleStatus(dto.Status),
	}
}
