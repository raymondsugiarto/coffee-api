package pensionbenefitrecipient

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, recipient *entity.PensionBenefitRecipientDto, tx *gorm.DB) (*entity.PensionBenefitRecipientDto, error)
	FindAll(ctx context.Context, req *entity.RecipientFindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.PensionBenefitRecipientDto, error)
	Update(ctx context.Context, recipient *entity.PensionBenefitRecipientDto, tx *gorm.DB) (*entity.PensionBenefitRecipientDto, error)
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, customerID string, recipients []*entity.PensionBenefitRecipientDto) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, recipient *entity.PensionBenefitRecipientDto, tx *gorm.DB) (*entity.PensionBenefitRecipientDto, error) {
	db := r.db
	if tx != nil {
		db = tx
	}
	recipientModel := recipient.ToModel()
	err := db.Create(recipientModel).Error
	if err != nil {
		return nil, err
	}
	return new(entity.PensionBenefitRecipientDto).FromModel(recipientModel), nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.RecipientFindAllRequest) (*pagination.ResultPagination, error) {
	var recipients []model.PensionBenefitRecipient = make([]model.PensionBenefitRecipient, 0)
	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.PensionBenefitRecipient{}).Preload("CountryBirth")
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &recipients,
		AllowedFields: []string{"customer_id"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.PensionBenefitRecipient)
	var data []*entity.PensionBenefitRecipientDto = make([]*entity.PensionBenefitRecipientDto, 0)
	for _, recipient := range *results {
		data = append(data, new(entity.PensionBenefitRecipientDto).FromModel(&recipient))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.PensionBenefitRecipientDto, error) {
	var recipient *model.PensionBenefitRecipient
	if err := r.db.Where("id = ?", id).Preload("CountryBirth").First(&recipient).Error; err != nil {
		return nil, err
	}
	return new(entity.PensionBenefitRecipientDto).FromModel(recipient), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.PensionBenefitRecipientDto, tx *gorm.DB) (*entity.PensionBenefitRecipientDto, error) {
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
	err := r.db.Where("id = ?", id).Delete(&model.PensionBenefitRecipient{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) BatchDelete(ctx context.Context, customerID string, recipients []*entity.PensionBenefitRecipientDto) error {
	recipientModels := make([]*model.PensionBenefitRecipient, len(recipients))
	for i, recipient := range recipients {
		recipientModels[i] = recipient.ToModel()
	}
	err := r.db.Where("customer_id = ?", customerID).Delete(&recipientModels).Error
	if err != nil {
		return err
	}
	return nil
}
