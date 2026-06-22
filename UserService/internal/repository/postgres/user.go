package postgres

import (
	"fmt"
	"user-service/internal/errorz"

	sq "github.com/Masterminds/squirrel"

	"context"
	"database/sql"
	"errors"

	"user-service/internal/models"

	"github.com/jmoiron/sqlx"
)

const (
	usersTableName = "users"

	usersIDColumn           = "id"
	usersEmailColumn        = "email"
	usersPasswordHashColumn = "password_hash"
	usersGradeColumn        = "grade"
	usersCreatedAtColumn    = "created_at"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create создаёт нового пользователя.
func (r *UserRepository) Create(ctx context.Context, user *models.User) (string, error) {
	const op = "Postgres.UserRepository.Create"

	builder := sq.Insert(usersTableName).
		Columns(usersEmailColumn, usersPasswordHashColumn, usersGradeColumn).
		Values(user.Email, user.Password, user.Grade).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var id string
	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

// GetByID получает пользователя по его ID.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	const op = "Postgres.UserRepository.GetByID"

	builder := sq.Select(usersIDColumn, usersEmailColumn, usersPasswordHashColumn, usersGradeColumn, usersCreatedAtColumn).
		From(usersTableName).
		Where(sq.Eq{usersIDColumn: id}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user models.User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorz.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

// GetByEmail получает пользователя по его почте.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "Postgres.UserRepository.GetByEmail"

	builder := sq.Select(usersIDColumn, usersEmailColumn, usersPasswordHashColumn, usersGradeColumn, usersCreatedAtColumn).
		From(usersTableName).
		Where(sq.Eq{usersEmailColumn: email}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user models.User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorz.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

// Update обновляет данные пользователя.
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	const op = "Postgres.UserRepository.Update"

	builder := sq.Update(usersTableName).
		Set(usersPasswordHashColumn, user.Password).
		Set(usersGradeColumn, user.Grade).
		Where(sq.Eq{usersIDColumn: user.ID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, errorz.ErrUserNotFound)
	}

	return nil
}

// Delete удаляет пользователя по его ID.
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	const op = "Postgres.UserRepository.Delete"

	builder := sq.Delete(usersTableName).
		Where(sq.Eq{usersIDColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, errorz.ErrUserNotFound)
	}

	return nil
}
