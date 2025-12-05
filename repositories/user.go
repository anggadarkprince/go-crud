package repositories

import (
	"context"
	"database/sql"

	"github.com/anggadarkprince/crud-employee-go/database"
	"github.com/anggadarkprince/crud-employee-go/models"
	"gitlab.com/tozd/go/errors"
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
        return nil, errors.Errorf("failed to query users: %w", err)
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
			return nil, errors.Errorf("failed to get user rows: %w", err)
		}

		users = append(users, user)
	}

	return &users, nil
}

func (repository *UserRepository) mapResultToUser(row *sql.Row) (*models.User, error) {
	if row.Err() != nil {
		return nil, errors.Errorf("failed to query user: %w", row.Err())
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
		return nil, errors.Errorf("user not found: %w", err)
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

func (repository *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users(name, username, email, password, user_type, status, avatar)
		VALUES(?, ?, ?, ?, ?, ?, ?)
	`
	result, err := repository.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.UserType,
		user.Status,
		user.Avatar,
	)

	if err != nil {
		return nil, errors.Errorf("failed to store user: %w", err)
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Errorf("failed to get last insert id: %w", err)
	}

	return repository.GetById(ctx, int(userId))
}

func (repository *UserRepository) UpdateAccount(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		UPDATE users SET name = ?, username = ?, email = ?, password = ?, avatar = ?
		WHERE id = ?
	`
	_, err := repository.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.Avatar,
		user.Id,
	)

	if err != nil {
		return nil, errors.Errorf("failed to update account: %w", err)
	}

	return repository.GetById(ctx, int(user.Id))
}