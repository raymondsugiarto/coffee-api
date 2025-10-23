package portfolio

// import (
// 	"context"

// 	e "github.com/raymondsugiarto/coffee-api/pkg/entity"
// 	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/customer"
// 	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
// )

// type Service interface {
// 	Create(ctx context.Context, dto *entity.PortfolioDto) (*entity.PortfolioDto, error)
// 	Get(ctx context.Context, id string) (*entity.PortfolioDto, error)
// 	Update(ctx context.Context, dto *entity.PortfolioDto) (*entity.PortfolioDto, error)
// 	Delete(ctx context.Context, id string) error
// 	FindAll(ctx context.Context, req *e.FindAllRequest) (*pagination.ResultPagination, error)
// }

// type service struct {
// 	repo Repository
// }

// func NewService(repo Repository) Service {
// 	return &service{repo}
// }

// func (s *service) Create(ctx context.Context, dto *entity.PortfolioDto) (*entity.PortfolioDto, error) {
// 	m, err := s.repo.Create(ctx, dto)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return m, nil
// }
// func (s *service) Get(ctx context.Context, id string) (*entity.PortfolioDto, error) {
// 	m, err := s.repo.Get(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return m, nil
// }

// func (s *service) GetByParticipantIdAndInvestmentProductId(ctx context.Context, participantId string, investmentProductId string) (*entity.PortfolioDto, error) {
// 	m, err := s.repo.GetByParticipantIdAndInvestmentProductId(ctx, participantId, investmentProductId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return m, nil
// }

// func (s *service) Update(ctx context.Context, dto *entity.PortfolioDto) (*entity.PortfolioDto, error) {
// 	m, err := s.repo.Update(ctx, dto)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return m, nil
// }

// func (s *service) Delete(ctx context.Context, id string) error {
// 	err := s.repo.Delete(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s *service) FindAll(ctx context.Context, req *e.FindAllRequest) (*pagination.ResultPagination, error) {
// 	m, err := s.repo.FindAll(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return m, nil
// }
