package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResourceRepositoryPostgres struct {
	pool *pgxpool.Pool
}

// type r interface {
// 	CreateFileResource(ctx context.Context, resource Resource) error
// 	CreateStorageObject(ctx context.Context, object storageObject) error
// 	CreateUserResource(ctx context.Context, resource Resource) error
// 	CreateWeekResource(ctx context.Context, resource Resource) error
// 	ObjectExists(ctx context.Context, hash string) (uuid.UUID, bool, error)
// }

func NewResourceRepositoryPostgres(p *pgxpool.Pool) *ResourceRepositoryPostgres {
	return &ResourceRepositoryPostgres{pool: p}
}

func (r *ResourceRepositoryPostgres) CreateFileResource(ctx context.Context, resource Resource) error {
	query := `INSERT INTO resources(id,name, type, storage_object_id) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, resource.ID, resource.Name, resource.ResourceType, resource.ObjectID)
	if err != nil {
		return fmt.Errorf("CreateFileResource err: %w", err)
	}
	return nil
}

func (r *ResourceRepositoryPostgres) CreateLinkResource(ctx context.Context, resource Resource) error {
	query := `INSERT INTO resources(id,name, type, external_url) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, resource.ID, resource.Name, resource.ResourceType, resource.ExternalLink)
	if err != nil {
		return fmt.Errorf("CreateLinkResource err: %w", err)
	}
	return nil
}

func (r *ResourceRepositoryPostgres) CreateStorageObject(ctx context.Context, object storageObject) error {
	query := `INSERT INTO storage_objects(id, hash, url) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, object.ID, object.Hash, object.URL)
	if err != nil {
		return fmt.Errorf("CreatStorageObject err: %w", err)
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
func (r *ResourceRepositoryPostgres) ObjectExists(ctx context.Context, hash string) (uuid.UUID, bool, error) {
	var id uuid.UUID
	query := `SELECT id FROM storage_objects WHERE hash=$1`
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

// for now just list name of the resources
func (r *ResourceRepositoryPostgres) ListResourcesByWeek(ctx context.Context, weekID uuid.UUID) ([]Resource, error) {
	query := `SELECT r.id, r.name, r.type, r.storage_object_id, r.external_url, o.user_id, r.created_at  FROM week_resources w JOIN resources r ON w.resource_id=r.id JOIN resource_owners o ON o.resource_id=w.resource_id WHERE week_id=$1;`

	rows, err := r.pool.Query(ctx, query, weekID)
	if err != nil {
		return []Resource{}, fmt.Errorf("ListResourceWeek query :%w", err)
	}
	defer rows.Close()
	resources := make([]Resource, 0)
	for rows.Next() {
		var resource Resource
		err := rows.Scan(&resource.ID, &resource.Name, &resource.ResourceType, &resource.ObjectID, &resource.ExternalLink, &resource.UserID, &resource.CreatedAt)
		if err != nil {
			return []Resource{}, fmt.Errorf("ListResourceWeek scan :%w", err)
		}
		resources = append(resources, resource)
	}

	return resources, err
}

func (r *ResourceRepositoryPostgres) LinkExists(ctx context.Context, resource Resource) (bool, error) {
	var link string
	query := `SELECT r.external_url FROM week_resources w JOIN resources r ON w.resource_id=r.id WHERE w.week_id=$1 AND r.external_url=$2`

	err := r.pool.QueryRow(ctx, query, resource.WeekID, *resource.ExternalLink).Scan(&link)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

// func (r *ResourceRepositoryPostgres) CreateLinkResource(ctx context.Context, resource Resource) (bool, error) {
// 	query := `IN`
// }
