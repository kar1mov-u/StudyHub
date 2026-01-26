package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewUserRepositoryPostgres(p *pgxpool.Pool) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{pool: p}
}

func (r *UserRepositoryPostgres) Create(ctx context.Context, user User) error {
	query := `
		INSERT INTO users (id, email, first_name, last_name, is_admin, password)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query, user.ID, user.Email, user.FirstName, user.LastName, user.IsAdmin, user.Password)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (User, error) {
	query := `
		SELECT id, email, first_name, last_name, is_admin, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.IsAdmin,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *UserRepositoryPostgres) GetByEmail(ctx context.Context, email string) (uuid.UUID, string, error) {
	var id uuid.UUID
	var password string
	query := `SELECT id, password FROM users WHERE email=$1`
	if err := r.pool.QueryRow(ctx, query, email).Scan(&id, &password); err != nil {
		return id, password, err
	}
	return id, password, nil

}

func (r *UserRepositoryPostgres) IsAdmin(ctx context.Context, id uuid.UUID) (bool, error) {
	var isAdmin bool
	query := `SELECT is_admin FROM users WHERE id=$1`
	err := r.pool.QueryRow(ctx, query, id).Scan(&isAdmin)
	return isAdmin, err
}

func (r *UserRepositoryPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *UserRepositoryPostgres) List(ctx context.Context) ([]User, error) {
	query := `
		SELECT id, email, first_name, last_name, is_admin, password
		FROM users
		ORDER BY email
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.IsAdmin,
			&user.Password,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}
