package repositories

import (
	"context"
	"database/sql"

	"github.com/anggadarkprince/crud-employee-go/database"
	"github.com/anggadarkprince/crud-employee-go/models"
)

type UserRepository struct {
	db database.Transaction
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) WithTx(tx *sql.Tx) *UserRepository {
    return &UserRepository{
        db: tx,
    }
}

func (repository *UserRepository) GetAll(ctx context.Context) (*[]models.User, error) {
	query := `
		SELECT id, name, username, email, password, user_type, status, avatar
		FROM users
		ORDER BY id DESC
	`
	rows, err := repository.db.QueryContext(ctx, query)
	
	if err != nil {
        return nil, err
    }
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User

		err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.UserType,
			&user.Status,
			&user.Avatar,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return &users, nil
}

func (repository *UserRepository) mapResultToUser(row *sql.Row) (*models.User, error) {
	if row.Err() != nil {
		return nil, row.Err()
	}

	var user models.User
	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.UserType,
		&user.Status,
		&user.Avatar,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepository) GetById(ctx context.Context, userId int) (*models.User, error) {
	query := `
		SELECT id, name, username, email, password, user_type, status, avatar
		FROM users WHERE id = ?
	`;
	row := repository.db.QueryRowContext(ctx, query, userId)

	return repository.mapResultToUser(row)
}

func (repository *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, username, email, password, user_type, status, avatar
		FROM users WHERE email = ?
	`;
	row := repository.db.QueryRowContext(ctx, query, email)

	return repository.mapResultToUser(row)
}

func (repository *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, name, username, email, password, user_type, status, avatar
		FROM users WHERE username = ?
	`;
	row := repository.db.QueryRowContext(ctx, query, username)
	
	return repository.mapResultToUser(row)
}