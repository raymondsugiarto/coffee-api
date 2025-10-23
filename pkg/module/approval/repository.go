package approval

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

const EMPTY = 0

type Repository interface {
	Create(ctx context.Context, dto *entity.ApprovalDto, tx *gorm.DB) (*entity.ApprovalDto, error)
	Get(ctx context.Context, id string) (*entity.ApprovalDto, error)
	GetByRefID(ctx context.Context, refID, approvalType string) (*entity.ApprovalDto, error)
	Update(ctx context.Context, dto *entity.ApprovalDto, tx *gorm.DB) (*entity.ApprovalDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.ApprovalFindAllRequest) (*pagination.ResultPagination, error)
	Confirmation(ctx context.Context, req *entity.ApprovalDto, cb func(*gorm.DB) error) (*entity.ApprovalDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Confirmation(ctx context.Context, req *entity.ApprovalDto, cb func(*gorm.DB) error) (*entity.ApprovalDto, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.Approval{}).Where("id = ?", req.ID).Updates(req.ToModel()).Error
		if err != nil {
			return err
		}
		err = cb(tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *repository) Create(ctx context.Context, dto *entity.ApprovalDto, tx *gorm.DB) (*entity.ApprovalDto, error) {
	db := r.db
	if tx != nil {
		db = tx
	}
	m := dto.ToModel()
	err := db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.ApprovalDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.ApprovalDto, error) {
	var m *model.Approval
	err := r.db.Where("id = ?", id).
		Preload("UserRequest").First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.ApprovalDto).FromModel(m), nil
}

func (r *repository) GetByRefID(ctx context.Context, refID, approvalType string) (*entity.ApprovalDto, error) {
	var m *model.Approval
	err := r.db.Where("ref_id = ? AND type = ?", refID, approvalType).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.ApprovalDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.ApprovalDto, tx *gorm.DB) (*entity.ApprovalDto, error) {
	db := r.db
	if tx != nil {
		db = tx
	}
	err := db.Updates(dto.ToModel()).Where("id = ? ", dto.ID).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.Approval{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.ApprovalFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Approval = make([]model.Approval, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		query := r.db.Model(&model.Approval{}).Preload("UserRequest").
			Joins("LEFT JOIN customer ON customer.user_id = approval.user_id_request")

		if len(req.Types) > EMPTY {
			query = query.Where("type IN ?", req.Types)
			if req.CompanyID != "" {
				requestRegistration := false
				for _, approvalType := range req.Types {
					if approvalType == model.ApprovalTypeCompany {
						requestRegistration = true
						break
					}
				}
				if requestRegistration {
					query = query.Joins("LEFT JOIN customer ON customer.id = approval.ref_id AND ref_table = 'customer'")
					query = query.Where("(type = 'COMPANY' AND ref_id = ?) OR (type = 'CUSTOMER' AND customer.id is not null AND customer.company_id = ?)", req.CompanyID, req.CompanyID)
				}
			}
		}

		if len(req.Statuses) > EMPTY {
			query = query.Where("status IN ?", req.Statuses)
		}

		return query
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"detail", "customer.first_name", "customer.last_name"},
		Data:          &m,
		AllowedFields: []string{"type", "status", "detail", "customer.first_name", "customer.last_name"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Approval)
	var data []*entity.ApprovalDto = make([]*entity.ApprovalDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.ApprovalDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
