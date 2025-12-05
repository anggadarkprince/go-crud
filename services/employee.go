package services

import (
	"context"
	"database/sql"

	"github.com/anggadarkprince/crud-employee-go/dto"
	"github.com/anggadarkprince/crud-employee-go/models"
	"github.com/anggadarkprince/crud-employee-go/repositories"
	"github.com/anggadarkprince/crud-employee-go/utilities"
)

type EmployeeService struct {
	employeeRepository *repositories.EmployeeRepository
	employeeAllowanceRepository *repositories.EmployeeAllowanceRepository
	db *sql.DB
}

func NewEmployeeService(
	employeeRepository *repositories.EmployeeRepository,
	employeeAllowanceRepository *repositories.EmployeeAllowanceRepository,
	db *sql.DB,
) *EmployeeService {
	return &EmployeeService{
		employeeRepository: employeeRepository,
		employeeAllowanceRepository: employeeAllowanceRepository,
		db: db,
	}
}

func (service *EmployeeService) GetAll(ctx context.Context) (*[]models.Employee, error) {
	return service.employeeRepository.GetAll(ctx)
}

func (service *EmployeeService) GetById(ctx context.Context, id int) (*models.Employee, error) {
	return service.employeeRepository.GetById(ctx, id)
}

func (service *EmployeeService) Store(ctx context.Context, data *dto.CreateEmployeeRequest) (*models.Employee, error) {
	hiredDate, err := utilities.StringToDate(data.HiredDate)
	
	if err != nil {
		return nil, err
	}

	tx, err := service.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	employeeModel := &models.Employee{
        Name: data.Name,
        Email: sql.NullString{String: data.Email, Valid: data.Email != ""},
        TaxNumber: sql.NullString{String: data.TaxNumber, Valid: data.TaxNumber != ""},
        Gender: sql.NullString{String: data.Gender, Valid: data.Gender != ""},
        Address: sql.NullString{String: data.Address, Valid: data.Address != ""},
        Status: sql.NullString{String: data.Status, Valid: data.Status != ""},
		HiredDate: hiredDate,
    }

	employeeRepository := service.employeeRepository.WithTx(tx)
	employee, err := employeeRepository.Store(ctx, employeeModel)
	if err != nil {
		return nil, err
	}
	
	employeeAllowanceRepository := service.employeeAllowanceRepository.WithTx(tx)
	_, err = employeeAllowanceRepository.StoreMany(ctx, employee.Id, data.Allowances)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return employee, nil
}

func (service *EmployeeService) Update(ctx context.Context, data *dto.UpdateEmployeeRequest) (*models.Employee, error) {
	hiredDate, err := utilities.StringToDate(data.HiredDate)
	
	if err != nil {
		return nil, err
	}

	tx, err := service.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	employeeModel := &models.Employee{
		Id: data.Id,
        Name: data.Name,
        Email: sql.NullString{String: data.Email, Valid: data.Email != ""},
        TaxNumber: sql.NullString{String: data.TaxNumber, Valid: data.TaxNumber != ""},
        Gender: sql.NullString{String: data.Gender, Valid: data.Gender != ""},
        Address: sql.NullString{String: data.Address, Valid: data.Address != ""},
        Status: sql.NullString{String: data.Status, Valid: data.Status != ""},
		HiredDate: hiredDate,
    }

	employeeRepository := service.employeeRepository.WithTx(tx)
	employee, err := employeeRepository.Update(ctx, employeeModel)
	if err != nil {
		return nil, err
	}
	
	employeeAllowanceRepository := service.employeeAllowanceRepository.WithTx(tx)
	err = employeeAllowanceRepository.DestroyByEmployeeId(ctx, employee.Id)
	if err != nil {
		return nil, err
	}
	_, err = employeeAllowanceRepository.StoreMany(ctx, employee.Id, data.Allowances)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return employee, nil
}

func (service *EmployeeService) Destroy(ctx context.Context, id int) error {
	tx, err := service.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	employeeRepository := service.employeeRepository.WithTx(tx)
	err = employeeRepository.Destroy(ctx, id)
	if err != nil {
		return err
	}
	
	employeeAllowanceRepository := service.employeeAllowanceRepository.WithTx(tx)
	err = employeeAllowanceRepository.DestroyByEmployeeId(ctx, id)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}