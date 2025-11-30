package orderitem

import (
	"context"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	Count(ctx context.Context, req *entity.OrderFindAllRequest) ([]entity.OrderItemPerItemCountDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Count(ctx context.Context, req *entity.OrderFindAllRequest) ([]entity.OrderItemPerItemCountDto, error) {
	startAt := req.OrderDate + " 17:00:00"
	// Parse the date and subtract 1 day
	if t, err := time.Parse("2006-01-02 15:04:05", startAt); err == nil {
		startAt = t.AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
	}
	var m []entity.OrderItemPerItemCountDto
	subquery := r.db.Model(&model.OrderItem{}).
		Joins(`JOIN "order" ON "order".id = order_item.order_id`).
		Where("order_at >= ? AND order_at < ? AND admin_id = ? AND company_id = ?", startAt, endAt, req.AdminID, req.CompanyID).
		Select("order_item.item_id, sum(qty) as total_qty, sum(qty * price) as total_price").
		Group("order_item.item_id")
	err := r.db.Table("(?) as sub", subquery).
		Joins(`JOIN item ON "item".id = sub.item_id`).
		Select("sub.total_qty, sub.total_price, item.name as item_name").
		Find(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}
