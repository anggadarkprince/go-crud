package repositories

import (
	"context"
	"database/sql"

	"gitlab.com/tozd/go/errors"

	"github.com/anggadarkprince/crud-employee-go/database"
	"github.com/anggadarkprince/crud-employee-go/models"
)

type EmployeeRepository struct {
	db database.Transaction
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) WithTx(tx *sql.Tx) *EmployeeRepository {
    return &EmployeeRepository{
        db: tx,
    }
}

func (repository *EmployeeRepository) GetAll(ctx context.Context) (*[]models.Employee, error) {
	query := `
		SELECT 
			id, name, email, tax_number, gender, hired_date, address, status, 
			(SELECT COUNT(*) FROM employee_allowances WHERE employee_id = employees.id) AS total_allowance 
		FROM employees
		ORDER BY id DESC
	`
	rows, err := repository.db.QueryContext(ctx, query)
	
	if err != nil {
        return nil, errors.Errorf("failed to query employees: %w", err)
    }
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var employee models.Employee

		err = rows.Scan(
			&employee.Id,
			&employee.Name,
			&employee.Email,
			&employee.TaxNumber,
			&employee.Gender,
			&employee.HiredDate,
			&employee.Address,
			&employee.Status,
			&employee.TotalAllowance,
		)
		if err != nil {
			return nil, errors.Errorf("failed to get employee rows: %w", err)
		}

		employees = append(employees, employee)
	}

	return &employees, nil
}

func (repository *EmployeeRepository) GetById(ctx context.Context, employeeId int) (*models.Employee, error) {
	query := `
		SELECT id, name, email, tax_number, gender, hired_date, address, status
		FROM employees WHERE id = ?
	`;
	row := repository.db.QueryRowContext(ctx, query, employeeId)
	if row.Err() != nil {
		return nil, errors.Errorf("failed to query employee id=%d: %w", employeeId, row.Err())
	}
	var employee models.Employee
	err := row.Scan(
		&employee.Id,
		&employee.Name,
		&employee.Email,
		&employee.TaxNumber,
		&employee.Gender,
		&employee.HiredDate,
		&employee.Address,
		&employee.Status,
	)
	if err != nil {
		return nil, errors.Errorf("employee not found id=%d: %w", employeeId, err)
	}

	return &employee, nil
}

func (repository *EmployeeRepository) Store(ctx context.Context, employee *models.Employee) (*models.Employee, error) {
	query := `
		INSERT INTO employees(name, email, tax_number, gender, hired_date, address, status)
		VALUES(?, ?, ?, ?, ?, ?, ?)
	`
	result, err := repository.db.ExecContext(
		ctx,
		query,
		employee.Name,
		employee.Email,
		employee.TaxNumber,
		employee.Gender,
		employee.HiredDate,
		employee.Address,
		employee.Status,
	)

	if err != nil {
		return nil, errors.Errorf("failed to store employee: %w", err)
	}

	employeeId, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Errorf("failed to get last id: %w", err)
	}

	return repository.GetById(ctx, int(employeeId))
}

func (repository *EmployeeRepository) Update(ctx context.Context, employee *models.Employee) (*models.Employee, error) {
	query := `
		UPDATE employees 
		SET name = ?, email = ?, tax_number = ?, gender = ?, hired_date = ?, address = ?, status = ? 
		WHERE id = ?
	`
	_, err := repository.db.ExecContext(
		ctx,
		query,
		employee.Name,
		employee.Email,
		employee.TaxNumber,
		employee.Gender,
		employee.HiredDate,
		employee.Address,
		employee.Status,
		employee.Id,
	)

	if err != nil {
		return nil, errors.Errorf("failed to update employee id=%d: %w", employee.Id, err)
	}

	return repository.GetById(ctx, employee.Id)
}

func (repository *EmployeeRepository) Destroy(ctx context.Context, employeeId int) (int64, error) {
	query := `DELETE FROM employees WHERE id = ?`
	result, err := repository.db.ExecContext(ctx, query, employeeId)
	if err != nil {
		return 0, errors.Errorf("failed to delete employee id=%d: %w", employeeId, err)
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Errorf("failed to get rows affected: %w", err)
	}
	return rowAffected, nil
}