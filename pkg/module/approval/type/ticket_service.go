package approvaltype

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/ticket"
	"gorm.io/gorm"
)

type TicketService interface {
	ConfirmationApprovalCallback(ctx context.Context, req *entity.ApprovalDto, tx *gorm.DB) (context.Context, error)
	FindByID(ctx context.Context, id string) (interface{}, error)
	NotifyApprovalCallback(ctx context.Context, req *entity.ApprovalDto) error
}

type ticketService struct {
	ticketService ticket.Service
}

func NewTicketService(is ticket.Service) TicketService {
	return &ticketService{
		ticketService: is,
	}
}

func (s *ticketService) ConfirmationApprovalCallback(ctx context.Context, req *entity.ApprovalDto, tx *gorm.DB) (context.Context, error) {
	ticketStatusMap := map[model.ApprovalStatus]model.TicketStatus{
		model.ApprovalStatusApproved: model.TicketStatusApproved,
		model.ApprovalStatusRejected: model.TicketStatusRejected,
	}

	return s.ticketService.ConfirmationApprovalCallback(ctx, &entity.TicketDto{
		ID:     req.RefID,
		Status: ticketStatusMap[model.ApprovalStatus(req.Status)],
	}, tx)
}

func (s *ticketService) FindByID(ctx context.Context, id string) (interface{}, error) {
	result, err := s.ticketService.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ticketService) NotifyApprovalCallback(ctx context.Context, req *entity.ApprovalDto) error {
	return nil
}
