package fee_setting

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	UpsertConfig(ctx context.Context, feeSetting *entity.FeeSettingDto) error
	GetConfig(ctx context.Context) (*entity.FeeSettingDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) UpsertConfig(ctx context.Context, feeSetting *entity.FeeSettingDto) error {
	var existingFeeSetting model.FeeSetting

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&existingFeeSetting).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				m := feeSetting.ToModel()
				return tx.Create(m).Error
			}
			return err
		}

		return tx.Model(&existingFeeSetting).
			Select("admin_fee", "operational_fee"). // sesuaikan dengan nama kolom di DB
			Updates(feeSetting.ToModel()).Error
	})

}

func (r *repository) GetConfig(ctx context.Context) (*entity.FeeSettingDto, error) {
	var m model.FeeSetting
	err := r.db.First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.FeeSettingDto).FromModel(&m), nil
}
