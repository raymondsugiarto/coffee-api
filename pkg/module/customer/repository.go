package customer

import (
	"context"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.CustomerDto, cb func(tx *gorm.DB) error) (*entity.CustomerDto, error)
	FindByReferralCode(ctx context.Context, referralCode string) (*entity.CustomerDto, error)
	FindByEmail(ctx context.Context, email string) (*entity.CustomerDto, error)
	FindByIDWithScope(ctx context.Context, id string, scopes []string) (*entity.CustomerDto, error)
	FindByID(ctx context.Context, id string) (*entity.CustomerDto, error)
	FindByUserID(ctx context.Context, userID string) (*entity.CustomerDto, error)
	FindByCompanyID(ctx context.Context, companyID string) (*entity.CustomerDto, error)
	Update(ctx context.Context, dto *entity.CustomerDto, cb func(tx *gorm.DB) error) (*entity.CustomerDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.CustomerFindAllRequest) (*pagination.ResultPagination, error)
	CountByType(ctx context.Context, companyID *string) (*entity.CustomerCount, error)
	CountByTypeThisMonth(ctx context.Context, companyID *string) (*entity.CustomerCount, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, dto *entity.CustomerDto, cb func(tx *gorm.DB) error) (*entity.CustomerDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	m := dto.ToModel()

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		dto.ID = m.ID
		dto.UserID = m.UserID
		dto.User.ID = m.UserID
		if err := cb(tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return new(entity.CustomerDto).FromModel(m), nil
}

func (r *repository) FindByReferralCode(ctx context.Context, referralCode string) (*entity.CustomerDto, error) {
	organizationID := shared.GetOrganization(ctx).ID
	var m *model.Customer
	if err := r.db.Where("referral_code = ? and organization_id = ?", referralCode, organizationID).
		Preload("Company").
		Preload("PensionBenefitRecipients.CountryBirth").
		Preload("CityOfBirth").
		Preload("Province").
		Preload("MailingProvince").
		Preload("Regency").
		Preload("MailingRegency").
		Preload("District").
		Preload("MailingDistrict").
		Preload("Village").
		Preload("MailingVillage").
		Preload("CountryBirth").
		Preload("Citizen").
		First(&m).Error; err != nil {
		return nil, err
	}
	return new(entity.CustomerDto).FromModel(m), nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*entity.CustomerDto, error) {
	organizationID := shared.GetOrganization(ctx).ID
	var m *model.Customer
	if err := r.db.Where("email = ? and organization_id = ?", email, organizationID).
		Preload("Company").
		Preload("PensionBenefitRecipients.CountryBirth").
		Preload("CityOfBirth").
		Preload("Province").
		Preload("MailingProvince").
		Preload("Regency").
		Preload("MailingRegency").
		Preload("District").
		Preload("MailingDistrict").
		Preload("Village").
		Preload("MailingVillage").
		Preload("CountryBirth").
		Preload("Citizen").
		First(&m).Error; err != nil {
		return nil, err
	}
	return new(entity.CustomerDto).FromModel(m), nil
}

func (r *repository) FindByIDWithScope(ctx context.Context, id string, scopes []string) (*entity.CustomerDto, error) {
	organizationID := shared.GetOrganization(ctx).ID
	var m *model.Customer

	db := r.db.WithContext(ctx).Where("id = ? and organization_id = ?", id, organizationID)
	for _, scope := range scopes {
		switch scope {
		case "complete":
			db = db.Scopes(m.Complete)
		}
	}
	if err := db.
		First(&m).Error; err != nil {
		return nil, err
	}
	return new(entity.CustomerDto).FromModel(m), nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.CustomerDto, error) {
	organizationID := shared.GetOrganization(ctx).ID
	var m *model.Customer
	if err := r.db.Preload("Company").Where("id = ? and organization_id = ?", id, organizationID).
		First(&m).Error; err != nil {
		return nil, err
	}
	return new(entity.CustomerDto).FromModel(m), nil
}

func (r *repository) FindByUserID(ctx context.Context, userID string) (*entity.CustomerDto, error) {
	organizationID := shared.GetOrganization(ctx).ID
	var m *model.Customer
	if err := r.db.Where("user_id = ? and organization_id = ?", userID, organizationID).
		Preload("Company").
		Preload("PensionBenefitRecipients.CountryBirth").
		Preload("CityOfBirth").
		Preload("Province").
		Preload("MailingProvince").
		Preload("Regency").
		Preload("MailingRegency").
		Preload("District").
		Preload("MailingDistrict").
		Preload("Village").
		Preload("MailingVillage").
		Preload("CountryBirth").
		Preload("Citizen").
		First(&m).Error; err != nil {
		return nil, err
	}
	return new(entity.CustomerDto).FromModel(m), nil
}

func (r *repository) FindByCompanyID(ctx context.Context, companyID string) (*entity.CustomerDto, error) {
	organizationID := shared.GetOrganization(ctx).ID
	var m *model.Customer
	if err := r.db.Where("company_id = ? and organization_id = ?", companyID, organizationID).
		Preload("Company").
		Preload("PensionBenefitRecipients.CountryBirth").
		Preload("CityOfBirth").
		Preload("Province").
		Preload("MailingProvince").
		Preload("Regency").
		Preload("MailingRegency").
		Preload("District").
		Preload("MailingDistrict").
		Preload("Village").
		Preload("MailingVillage").
		Preload("CountryBirth").
		Preload("Citizen").
		Find(&m).Error; err != nil {
		return nil, err
	}
	return new(entity.CustomerDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.CustomerDto, cb func(tx *gorm.DB) error) (*entity.CustomerDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Customer{}).Where("id = ?", dto.ID).Updates(dto.ToModel()).Error; err != nil {
			return err
		}
		if err := cb(tx); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload the updated customer from the database
	updated, err := r.FindByIDWithScope(ctx, dto.ID, []string{"complete"})
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	organizationID := shared.GetOrganization(ctx).ID
	err := r.db.Where("id = ? and organization_id = ?", id, organizationID).Delete(&model.Customer{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.CustomerFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Customer = make([]model.Customer, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.Customer{}).
			Joins("Company").
			Preload("Company").
			Preload("PensionBenefitRecipients.CountryBirth").
			Preload("CityOfBirth").
			Preload("Province").
			Preload("MailingProvince").
			Preload("Regency").
			Preload("MailingRegency").
			Preload("District").
			Preload("MailingDistrict").
			Preload("Village").
			Preload("MailingVillage").
			Preload("CountryBirth").
			Preload("Citizen")
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"customer.first_name", `"Company".first_name`, "customer.email", "sim_status", "approval_status"},
		Data:          &m,
		AllowedFields: []string{"customer_id_parent", "company_id", "sim_status"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Customer)
	var data []*entity.CustomerDto = make([]*entity.CustomerDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.CustomerDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) CountByType(ctx context.Context, companyID *string) (*entity.CustomerCount, error) {
	var count entity.CustomerCount

	subquery := r.db.WithContext(ctx).Model(&model.Customer{}).
		Where("sim_status = 'ACTIVE'").
		Group("company_id").
		Select("NULLIF(TRIM(company_id), '') AS company_id, COUNT(*) as total_customer")

	query := r.db.WithContext(ctx).Table("(?) as cg", subquery).
		Joins("LEFT JOIN company ON company.id = cg.company_id").
		Select(`COALESCE(SUM(cg.total_customer) FILTER (WHERE company.company_type = 'DKP'), 0) AS dkp,
			COALESCE(SUM(cg.total_customer) FILTER (WHERE company.company_type = 'PPIP'), 0) AS ppip_corporate,
			COALESCE(SUM(cg.total_customer) FILTER (WHERE cg.company_id IS null), 0) AS ppip_mandiri`)

	if companyID != nil {
		query = query.Where("company.id = ?", companyID)
	}

	err := query.Take(&count).Error
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (r *repository) CountByTypeThisMonth(ctx context.Context, companyID *string) (*entity.CustomerCount, error) {
	var count entity.CustomerCount

	// this month
	thisMonth := time.Now().UTC().Month()
	nextMonth := thisMonth + 1
	if nextMonth > 12 {
		nextMonth = 1
	}
	startDate := time.Date(time.Now().Year(), thisMonth, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(time.Now().Year(), nextMonth, 1, 0, 0, 0, 0, time.UTC)

	subquery := r.db.WithContext(ctx).Model(&model.Customer{}).
		Where("effective_date >= ? AND effective_date < ?", startDate, endDate).
		Group("company_id").
		Select("NULLIF(TRIM(company_id), '') AS company_id, COUNT(*) as total_customer")

	query := r.db.WithContext(ctx).Table("(?) as cg", subquery).
		Joins("LEFT JOIN company ON company.id = cg.company_id").
		Select(`COALESCE(SUM(cg.total_customer) FILTER (WHERE company.company_type = 'DKP'), 0) AS dkp,
			COALESCE(SUM(cg.total_customer) FILTER (WHERE company.company_type = 'PPIP'), 0) AS ppip_corporate,
			COALESCE(SUM(cg.total_customer) FILTER (WHERE cg.company_id IS null), 0) AS ppip_mandiri`)

	if companyID != nil {
		query = query.Where("company.id = ?", companyID)
	}

	err := query.Take(&count).Error
	if err != nil {
		return nil, err
	}
	return &count, nil
}
