package driver

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Service interface {
	FindAllDrivers(ctx context.Context, req *entity.DriverFindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (s *service) FindAllDrivers(ctx context.Context, req *entity.DriverFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Admin = make([]model.Admin, 0)
	tbl := pagination.NewTable(s.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := s.db.Model(&model.Admin{}).Where("admin_type = ?", "EMPLOYEE")
		if req.Query != "" {
			like := "%" + req.Query + "%"
			q = q.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?", like, like, like)
		}
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"first_name", "last_name", "email"},
		Data:          &m,
		AllowedFields: []string{"first_name", "created_at"},
	})
	if err != nil {
		return nil, err
	}
	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Admin)
	data := make([]*entity.AdminDto, 0, len(*results))
	for i := range *results {
		d := &entity.AdminDto{}
		d.FromModel(&(*results)[i])
		data = append(data, d)
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
