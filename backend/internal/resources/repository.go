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

func (r *ResourceRepositoryPostgres) ListResourcesByWeek(ctx context.Context, weekID uuid.UUID) ([]ResourceWithUser, error) {
	query := `SELECT r.id, r.name, r.type, r.storage_object_id, r.external_url, o.user_id, r.created_at, u.first_name FROM week_resources w JOIN resources r ON w.resource_id=r.id JOIN resource_owners o ON o.resource_id=w.resource_id JOIN users u ON o.user_id=u.id WHERE week_id=$1;`

	rows, err := r.pool.Query(ctx, query, weekID)
	if err != nil {
		return []ResourceWithUser{}, fmt.Errorf("ListResourceWeek query :%w", err)
	}
	defer rows.Close()
	resources := make([]ResourceWithUser, 0)
	for rows.Next() {
		var resource ResourceWithUser
		err := rows.Scan(&resource.ID, &resource.Name, &resource.ResourceType, &resource.ObjectID, &resource.ExternalLink, &resource.UserID, &resource.CreatedAt, &resource.UserName)
		if err != nil {
			return []ResourceWithUser{}, fmt.Errorf("ListResourceWeek scan :%w", err)
		}
		resources = append(resources, resource)
	}

	return resources, err
}

func (r *ResourceRepositoryPostgres) ListUserResources(ctx context.Context, userID uuid.UUID) ([]UserResources, error) {
	query := `SELECT ro.resource_id, w.id, w.number, ro.user_id, m.name, mr.semester, mr.year, r.storage_object_id,
	r.external_url, r.type, r.name, r.created_at FROM resource_owners ro JOIN week_resources wr ON ro.resource_id=wr.resource_id 
	JOIN resources r ON r.id=ro.resource_id JOIN weeks w ON wr.week_id=w.id JOIN module_runs mr ON w.module_run_id=mr.id 
	JOIN modules m ON mr.module_id=m.id WHERE ro.user_id=$1;`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return []UserResources{}, fmt.Errorf("ListUserResources err: %w", err)
	}
	defer rows.Close()
	resources := make([]UserResources, 0)
	for rows.Next() {
		var resource UserResources
		err = rows.Scan(&resource.ID, &resource.WeekID, &resource.WeekNumber, &resource.UserID, &resource.ModuleName, &resource.Semester, &resource.Year, &resource.ObjectID, &resource.ExternalLink, &resource.ResourceType, &resource.Name, &resource.CreatedAt)
		if err != nil {
			return []UserResources{}, fmt.Errorf("ListUserResources scan err: %w", err)
		}
		resources = append(resources, resource)
	}
	return resources, err
}

func (r *ResourceRepositoryPostgres) LinkExistsInWeek(ctx context.Context, resource Resource) (bool, error) {
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

func (r *ResourceRepositoryPostgres) FileExistsInWeek(ctx context.Context, hash string, weekID uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM week_resources w JOIN resources r ON w.resource_id=r.id JOIN storage_objects s ON r.storage_object_id=s.id WHERE s.hash=$1 AND w.week_id=$2;`
	var i int

	err := r.pool.QueryRow(ctx, query, hash, weekID).Scan(&i)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *ResourceRepositoryPostgres) ListOrphanObjects(ctx context.Context) ([]uuid.UUID, error) {
	query := `SELECT so.id FROM storage_objects so LEFT JOIN resources r ON so.id=r.storage_object_id WHERE r.id IS NULL`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return []uuid.UUID{}, err
	}
	defer rows.Close()
	ids := make([]uuid.UUID, 0)
	for rows.Next() {
		var id uuid.UUID
		err := rows.Scan(&id)
		if err != nil {
			return []uuid.UUID{}, err
		}
		ids = append(ids, id)

	}
	return ids, nil
}

func (r *ResourceRepositoryPostgres) DeleteStorageObjecst(ctx context.Context, ids []uuid.UUID) error {
	batch := pgx.Batch{}
	query := `DELETE FROM storage_objects WHERE id=$1`
	for _, id := range ids {
		batch.Queue(query, id)
	}
	err := r.pool.SendBatch(ctx, &batch).Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *ResourceRepositoryPostgres) DeleteResource(ctx context.Context, userID, resourceID uuid.UUID) error {
	query := `DELETE FROM resources WHERE id IN (SELECT resource_id FROM resource_owners WHERE resource_id=$1 AND user_id =$2)`

	_, err := r.pool.Exec(ctx, query, resourceID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}
	return nil

}
