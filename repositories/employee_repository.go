package repositories

import (
	"context"
	"database/sql"

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
        return nil, err
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
			return nil, err
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
		return nil, row.Err()
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
		return nil, err
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
		return nil, err
	}

	employeeId, err := result.LastInsertId()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return repository.GetById(ctx, employee.Id)
}

func (repository *EmployeeRepository) Destroy(ctx context.Context, employeeId int) error {
	query := `DELETE FROM employees WHERE id = ?`
	_, err := repository.db.ExecContext(ctx, query, employeeId)
	return err
}