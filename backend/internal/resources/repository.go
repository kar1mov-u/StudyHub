package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// type ResourceRepository interface {
// 	GetResourceByID(context.Context, uuid.UUID) (Resource, error)
// 	SaveResource(context.Context, Resource) error
// 	ListResourcesByWeek(context.Context, uuid.UUID) ([]Resource, error)
// }

type ResourceRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewResourceRepositoryPostgres(p *pgxpool.Pool) *ResourceRepositoryPostgres {
	return &ResourceRepositoryPostgres{pool: p}
}

func (r *ResourceRepositoryPostgres) Create(ctx context.Context, resource Resource) error {
	query := `INSERT INTO resources(id, type, hash, url) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, resource.ID, resource.ResourceType, resource.Hash, resource.Url)
	if err != nil {
		return fmt.Errorf("CreateResource err: %w", err)
	}
	return nil
}

func (r *ResourceRepositoryPostgres) CreateUserResource(ctx context.Context, resource Resource) error {
	query := `INSERT INTO resource_owners (resource_id, user_id) VALUES ($1, $2)`
	_, err := r.pool.Exec(ctx, query, resource.ID, resource.UserID)
	if err != nil {
		return fmt.Errorf("CreateResourceOwner err: %w", err)
	}
	return nil
}

func (r *ResourceRepositoryPostgres) CreateWeekResource(ctx context.Context, resource Resource) error {
	query := `INSERT INTO week_resources (resource_id, week_id) VALUES ($1, $2)`
	_, err := r.pool.Exec(ctx, query, resource.ID, resource.WeekID)
	if err != nil {
		return fmt.Errorf("CreateWeekResource err: %w", err)
	}
	return nil
}

// this function will check if the resource with this hash exists, if yes returns id of it, if no it will return false
func (r *ResourceRepositoryPostgres) ResourceExists(ctx context.Context, hash string) (uuid.UUID, bool, error) {
	var id uuid.UUID
	query := `SELECT id FROM resources WHERE hash=$1`
	row := r.pool.QueryRow(ctx, query, hash)
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, false, nil
		}
		return uuid.UUID{}, false, err
	}
	return id, true, nil

}

func (r *ResourceRepositoryPostgres) ListResourcesByWeek(ctx context.Context, id uuid.UUID) ([]Resource, error) {
	query := `SELECT resource_id FROM week_resources WHERE week_id=$1`
	rows, err := r.pool.Query(ctx, query, id)

	if err != nil {
		return []Resource{}, fmt.Errorf("ListResourceByWeek err: %w", err)
	}
	defer rows.Close()

	resources := make([]Resource, 0)
	for rows.Next() {
		var resource Resource
		err := rows.Scan(&resource.ID)
		if err != nil {
			return []Resource{}, fmt.Errorf("ListReousrceByWeek scan err: %w", err)
		}
		resources = append(resources, resource)
	}
	return resources, nil
}
