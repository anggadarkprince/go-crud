package repositories

import (
	"context"
	"database/sql"

	"github.com/anggadarkprince/crud-employee-go/database"
	"github.com/anggadarkprince/crud-employee-go/models"
	"gitlab.com/tozd/go/errors"
)

type EmployeeAllowanceRepository struct {
	db database.Transaction
}

func NewEmployeeAllowanceRepository(db *sql.DB) *EmployeeAllowanceRepository {
	return &EmployeeAllowanceRepository{db: db}
}

func (r *EmployeeAllowanceRepository) GetByEmployeeId(ctx context.Context, employeeId int) (*[]models.EmployeeAllowance, error) {
	query := `SELECT id, employee_id, allowance FROM employee_allowances WHERE employee_id = ?`;
	rows, err := r.db.QueryContext(ctx, query, employeeId)
	if err != nil {
		return nil, errors.Errorf("failed to query employee allowance by employee id=%d: %w", employeeId, err)
	}
	defer rows.Close()

	var employeeAllowances []models.EmployeeAllowance
	for rows.Next() {
		var employeeAllowance models.EmployeeAllowance

		err = rows.Scan(
			&employeeAllowance.Id,
			&employeeAllowance.EmployeeId,
			&employeeAllowance.Allowance,
		)
		if err != nil {
			return nil, errors.Errorf("allowance by employee not found id=%d: %w", employeeId, err)
		}

		employeeAllowances = append(employeeAllowances, employeeAllowance)
	}
	return &employeeAllowances, nil
}

func (r *EmployeeAllowanceRepository) WithTx(tx *sql.Tx) *EmployeeAllowanceRepository {
    return &EmployeeAllowanceRepository{
        db: tx,
    }
}

func (repository *EmployeeAllowanceRepository) GetById(ctx context.Context, id int) (*models.EmployeeAllowance, error) {
	query := `SELECT id, employee_id, allowance FROM employee_allowances WHERE id = ?`
	row := repository.db.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, errors.Errorf("failed to query employee allowance id=%d: %w", id, row.Err())
	}
	var employeeAllowance models.EmployeeAllowance
	err := row.Scan(
		&employeeAllowance.Id,
		&employeeAllowance.EmployeeId,
		&employeeAllowance.Allowance,
	)
	if err != nil {
		return nil, errors.Errorf("employee allowance not found id=%d: %w", id, err)
	}
	return &employeeAllowance, nil
}

func (repository *EmployeeAllowanceRepository) Store(ctx context.Context, employeeAllowance *models.EmployeeAllowance) (*models.EmployeeAllowance, error) {
	query := `INSERT INTO employee_allowances(employee_id, allowance) VALUES(?, ?)`
	result, err := repository.db.ExecContext(ctx, query)
	if err != nil {
		return nil, errors.Errorf("failed to store employee allowance: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Errorf("failed to get last id: %w", err)
	}
	employeeAllowance.Id = int(id)
	return employeeAllowance, nil
}

func (repository *EmployeeAllowanceRepository) StoreMany(ctx context.Context, employeeId int, allowances []string) (*[]models.EmployeeAllowance, error) {
	query := `INSERT INTO employee_allowances(employee_id, allowance) VALUES(?, ?)`
	statement, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Errorf("failed to prepare statement: %w", err)
	}
	defer statement.Close()

	var employeeAllowances []models.EmployeeAllowance
	for _, item := range allowances {
		result, err := statement.Exec(employeeId, item)
		if err != nil {
			return nil, errors.Errorf("failed to store employee allowance: %w", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return nil, errors.Errorf("failed to get last id: %w", err)
		}
		employeeAllowances = append(employeeAllowances, models.EmployeeAllowance{
			Id: int(id),
			EmployeeId: employeeId,
			Allowance: item,
		})
	}
	return &employeeAllowances, nil
}

func (repository *EmployeeAllowanceRepository) Update(ctx context.Context, employeeAllowance *models.EmployeeAllowance) (*models.EmployeeAllowance, error) {
	query := `
		UPDATE employee_allowances 
		SET employee_id = ?, allowance = ? 
		WHERE id = ?
	`
	_, err := repository.db.ExecContext(
		ctx,
		query,
		employeeAllowance.EmployeeId,
		employeeAllowance.Allowance,
		employeeAllowance.Id,
	)
	if err != nil {
		return nil, errors.Errorf("failed to update allowance id=%d: %w", employeeAllowance.Id, err)
	}
	
	return repository.GetById(ctx, employeeAllowance.Id)
}

func (repository *EmployeeAllowanceRepository) DestroyByEmployeeId(ctx context.Context, employeeId int) (int64, error) {
	query := `DELETE FROM employee_allowances WHERE employee_id = ?`
	result, err := repository.db.ExecContext(ctx, query, employeeId)
	if err != nil {
		return 0, errors.Errorf("failed to delete allowance by employee id=%d: %w", employeeId, err)
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Errorf("failed to get rows affected: %w", err)
	}
	return rowAffected, nil
}