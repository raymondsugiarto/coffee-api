package ticket

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	e "github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/approval"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	"github.com/raymondsugiarto/coffee-api/pkg/module/user"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, dto *entity.TicketDto) (*entity.TicketDto, error)
	FindByID(ctx context.Context, id string) (*entity.TicketDto, error)
	Update(ctx context.Context, dto *entity.TicketDto) (*entity.TicketDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *e.FindAllRequest) (*pagination.ResultPagination, error)

	ConfirmationApprovalCallback(ctx context.Context, req *entity.TicketDto, tx *gorm.DB) (context.Context, error)
}

type service struct {
	repo            Repository
	approvalService approval.Service
	customerService customer.Service
	userService     user.Service
}

func NewService(repo Repository, approvalService approval.Service, customerService customer.Service, userService user.Service) Service {
	return &service{repo, approvalService, customerService, userService}
}

func (s *service) Create(ctx context.Context, dto *entity.TicketDto) (*entity.TicketDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.CreateWithTx(ctx, dto, func(tx *gorm.DB) error {
		if _, err := s.approvalService.CreateWithTx(ctx, dto.ToApprovalSubmitDto(dto.ID), tx); err != nil {
			return err
		}
		return nil
	})
}

func (s *service) ConfirmationApprovalCallback(ctx context.Context, dto *entity.TicketDto, tx *gorm.DB) (context.Context, error) {
	ticket, err := s.repo.Get(ctx, dto.ID)
	if err != nil {
		return ctx, err
	}

	ticket.Status = dto.Status

	_, err = s.repo.Update(ctx, ticket)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.TicketDto, error) {
	ticket, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := s.userService.FindByID(ctx, ticket.UserID)
	if err != nil {
		return nil, err
	}

	if user.UserType == entity.CUSTOMER {
		customer, err := s.customerService.FindByUserID(ctx, ticket.UserID)
		if err == nil {
			ticket.Customer = customer
		}
	}

	return ticket, nil
}

func (s *service) Update(ctx context.Context, dto *entity.TicketDto) (*entity.TicketDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *e.FindAllRequest) (*pagination.ResultPagination, error) {
	res, err := s.repo.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}

	tickets := res.Data.(*[]model.Ticket)
	var data []*entity.TicketDto

	for _, t := range *tickets {
		ticketDto := new(entity.TicketDto).FromModel(&t)

		user, err := s.userService.FindByID(ctx, ticketDto.UserID)
		if err != nil {
			return nil, err
		}

		if user.UserType == entity.CUSTOMER {
			customer, err := s.customerService.FindByUserID(ctx, ticketDto.UserID)
			if err == nil {
				ticketDto.Customer = customer
			}
		}

		data = append(data, ticketDto)
	}

	return &pagination.ResultPagination{
		Data:        data,
		Page:        res.Page,
		Count:       res.Count,
		RowsPerPage: res.RowsPerPage,
		TotalPages:  res.TotalPages,
	}, nil

}
