package bankcustomer

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	Create(ctx context.Context, dto *entity.BankCustomerDto) (*entity.BankCustomerDto, error)
	FindByID(ctx context.Context, id string) (*entity.BankCustomerDto, error)
	Update(ctx context.Context, dto *entity.BankCustomerDto) (*entity.BankCustomerDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.BankCustomerFindAllRequest) (*pagination.ResultPagination, error)
	FindByCustomer(ctx context.Context, req *entity.BankCustomerFindAllRequest) (*pagination.ResultPagination, error)
	SetDefaultBankCustomer(ctx context.Context, id string) (*entity.BankCustomerDto, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, dto *entity.BankCustomerDto) (*entity.BankCustomerDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	dto.CustomerID = shared.GetUserCredential(ctx).CustomerID
	if dto.IsDefault {
		// If the bank customer is set as default, ensure all others are set to non-default
		err := s.repo.SetAllNonDefault(ctx, dto.CustomerID)
		if err != nil {
			return nil, err
		}
	}
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.BankCustomerDto, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.BankCustomerDto) (*entity.BankCustomerDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	dto.CustomerID = shared.GetUserCredential(ctx).CustomerID
	if dto.IsDefault {
		// If the bank customer is set as default, ensure all others are set to non-default
		err := s.repo.SetAllNonDefault(ctx, dto.CustomerID)
		if err != nil {
			return nil, err
		}
	}
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.BankCustomerFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repo.FindAll(ctx, req)
}

func (s *service) FindByCustomer(ctx context.Context, req *entity.BankCustomerFindAllRequest) (*pagination.ResultPagination, error) {
	req.CustomerID = shared.GetUserCredential(ctx).CustomerID
	return s.repo.FindAll(ctx, req)
}

func (s *service) SetDefaultBankCustomer(ctx context.Context, id string) (*entity.BankCustomerDto, error) {
	dto, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Set all other bank customers to non-default
	err = s.repo.SetAllNonDefault(ctx, dto.CustomerID)
	if err != nil {
		return nil, err
	}

	// Set the specified bank customer as default
	dto.IsDefault = true
	return s.repo.Update(ctx, dto)
}
