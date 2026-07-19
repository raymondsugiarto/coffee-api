package payroll

import (
	"context"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	Simulate(ctx context.Context, req *entity.SimulatePayrollRequest) (*entity.SimulatePayrollResultDto, error)
	Save(ctx context.Context, req *entity.SavePayrollRequest) (*entity.EmployeeSalaryDto, error)
	FindAll(ctx context.Context, req *entity.EmployeeSalaryFindAllRequest) (*pagination.ResultPagination, error)
	FindOne(ctx context.Context, id string) (*entity.EmployeeSalaryDto, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Simulate(
	ctx context.Context,
	req *entity.SimulatePayrollRequest,
) (*entity.SimulatePayrollResultDto, error) {
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, err
	}
	return s.repo.Simulate(ctx, req.AdminIDEmployee, start, end)
}

// Save turns the operator-approved simulation into one persisted
// employee_salary header row + N employee_salary_component rows.
// The frontend submits the same SimulatePayrollResultDto it just
// rendered (minus the session evidence), plus the operator-entered
// TotalCashReceipt, and the service flattens it into the schema.
func (s *service) Save(
	ctx context.Context,
	req *entity.SavePayrollRequest,
) (*entity.EmployeeSalaryDto, error) {
	orgID := shared.GetOrganization(ctx).ID

	header := &entity.EmployeeSalaryDto{
		OrganizationID:     orgID,
		AdminIDEmployee:    req.AdminIDEmployee,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		TotalMealAllowance: req.TotalMealAllowance,
		TotalAttendance:    req.TotalAttendance,
		TotalCommission:    req.TotalCommission,
		TotalBonusTarget:   req.TotalBonusTarget,
		TotalSalary:        req.TotalSalary,
		TotalCashReceipt:   req.TotalCashReceipt,
		RemainingSalary:    req.TotalSalary - req.TotalCashReceipt,
		Components:         req.Components,
	}
	return s.repo.Save(ctx, header)
}

func (s *service) FindAll(
	ctx context.Context,
	req *entity.EmployeeSalaryFindAllRequest,
) (*pagination.ResultPagination, error) {
	return s.repo.FindAll(ctx, req)
}

func (s *service) FindOne(
	ctx context.Context,
	id string,
) (*entity.EmployeeSalaryDto, error) {
	return s.repo.FindOne(ctx, id)
}
