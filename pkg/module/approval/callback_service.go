package approval

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"gorm.io/gorm"
)

type CallbackService interface {
	ConfirmationApprovalCallback(ctx context.Context, req *entity.ApprovalDto, tx *gorm.DB) (context.Context, error)
	FindByID(ctx context.Context, id string) (interface{}, error)
	NotifyApprovalCallback(ctx context.Context, req *entity.ApprovalDto) error
}
