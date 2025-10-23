package entity

import "github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"

type FindAllRequest struct {
	OrganizationData
	pagination.GetListRequest
}
