package services

import (
	"context"

	"github.com/anggadarkprince/crud-employee-go/models"
	"github.com/anggadarkprince/crud-employee-go/repositories"
)

type EmployeeAllowanceService struct {
	employeeAllowanceRepository *repositories.EmployeeAllowanceRepository
}

func NewEmployeeAllowanceService(
	employeeAllowanceRepository *repositories.EmployeeAllowanceRepository,
) *EmployeeAllowanceService {
	return &EmployeeAllowanceService{
		employeeAllowanceRepository: employeeAllowanceRepository,
	}
}

func (service *EmployeeAllowanceService) GetByEmployeeId(ctx context.Context, employeeId int) (*[]models.EmployeeAllowance, error) {
	return service.employeeAllowanceRepository.GetByEmployeeId(ctx, employeeId)
}