package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type Article struct {
	concern.CommonWithIDs
	ImageUrl  string
	Title     string
	Slug      string
	Content   string
	Status    string
	CreatedBy string
	UpdatedBy string
}
