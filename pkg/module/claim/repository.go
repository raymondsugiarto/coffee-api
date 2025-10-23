package claim

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	FindByID(ctx context.Context, id string) (*entity.ClaimDto, error)
	Create(ctx context.Context, dto *entity.ClaimDto, cb func(tx *gorm.DB) error) (*entity.ClaimDto, error)
	Update(ctx context.Context, dto *entity.ClaimDto, tx *gorm.DB) (*entity.ClaimDto, error)
	FindAll(ctx context.Context, req *entity.ClaimFindAllRequest) (*pagination.ResultPagination, error)
	FindByCompanyID(ctx context.Context, companyID string) ([]*entity.ClaimDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.ClaimDto, error) {
	var m *model.Claim
	if err := r.db.Where("id = ?", id).
		Preload("Participant.Customer.Company").First(&m).Error; err != nil {
		return nil, err
	}

	return new(entity.ClaimDto).FromModel(m), nil
}

func (r *repository) Create(ctx context.Context, dto *entity.ClaimDto, cb func(tx *gorm.DB) error) (*entity.ClaimDto, error) {
	m := dto.ToModel()

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		dto.ID = m.ID
		if err := cb(tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return new(entity.ClaimDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.ClaimDto, tx *gorm.DB) (*entity.ClaimDto, error) {
	if tx == nil {
		tx = r.db
	}

	m := dto.ToModel()

	if err := tx.Model(m).Where("id = ?", dto.ID).Updates(m).Error; err != nil {
		return nil, err
	}

	return new(entity.ClaimDto).FromModel(m), nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.ClaimFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Claim = make([]model.Claim, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.Claim{}).
			Preload("Participant")

		if req.CompanyID != nil {
			q = q.Joins("JOIN participant ON participant.id = claim.participant_id").
				Joins("JOIN customer ON customer.id = participant.customer_id").
				Joins("JOIN company ON company.id = customer.company_id").
				Preload("Participant.Customer.Company").
				Where("company.id = ?", req.CompanyID)
		}

		if req.CustomerID != "" {
			q = q.Joins("JOIN participant ON participant.id = claim.participant_id").
				Joins("JOIN customer ON customer.id = participant.customer_id").
				Where("customer.id = ?", req.CustomerID)
		}

		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"Participant.Customer.first_name", "Participant.code", "amount", "approval_status"},
		Data:          &m,
		AllowedFields: []string{"participant_id", "status"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Claim)
	var data []*entity.ClaimDto = make([]*entity.ClaimDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.ClaimDto).FromModel(&m))
	}

	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) FindByCompanyID(ctx context.Context, companyID string) ([]*entity.ClaimDto, error) {
	var m []model.Claim = make([]model.Claim, 0)
	err := r.db.WithContext(ctx).Joins("JOIN participant ON participant.id = claim.participant_id").
		Joins("JOIN customer ON customer.id = participant.customer_id").
		Where("customer.company_id = ?", companyID).
		Preload("Participant").
		Find(&m).Error
	if err != nil {
		return nil, err
	}

	var data []*entity.ClaimDto = make([]*entity.ClaimDto, 0)
	for _, m := range m {
		data = append(data, new(entity.ClaimDto).FromModel(&m))
	}
	return data, nil
}
