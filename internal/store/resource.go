package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/viniosilva/socialassistanceapi/internal/exception"
	"github.com/viniosilva/socialassistanceapi/internal/model"
)

//go:generate mockgen -destination ../../mock/resource_store_mock.go -package mock . ResourceStore
type ResourceStore interface {
	FindAll(ctx context.Context) ([]model.Resource, error)
	FindOneById(ctx context.Context, resourceID int) (*model.Resource, error)
	Create(ctx context.Context, resource model.Resource) (*model.Resource, error)
	Update(ctx context.Context, resource model.Resource) (*model.Resource, error)
}

type resourceStore struct {
	db *sql.DB
}

func NewResourceStore(db *sql.DB) ResourceStore {
	return &resourceStore{db}
}

func (iml *resourceStore) FindAll(ctx context.Context) ([]model.Resource, error) {
	resources := []model.Resource{}

	res, err := iml.db.Query(`
		SELECT id,
			created_at,
			updated_at,
			name,
			amount,
			measurement
		FROM resources`)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		resource, err := scanResource(res)
		if err != nil {
			return nil, err
		}
		resources = append(resources, *resource)
	}
	return resources, nil
}

func (impl *resourceStore) FindOneById(ctx context.Context, resourceID int) (*model.Resource, error) {
	res, err := impl.db.QueryContext(ctx, `
		SELECT id,
			created_at,
			updated_at,
			name,
			amount,
			measurement
		FROM resources
		WHERE id = ?
		LIMIT 1 `, resourceID)
	if err != nil {
		return nil, err
	}

	var resource *model.Resource
	for res.Next() {
		resource, err = scanResource(res)
		if err != nil {
			return nil, err
		}
	}

	return resource, nil
}

func (impl *resourceStore) Create(ctx context.Context, resource model.Resource) (*model.Resource, error) {
	now := time.Now()
	nowMysql := now.Format("2006-01-02T15:04:05")
	res, err := impl.db.ExecContext(ctx, `
		INSERT INTO resources (created_at, updated_at, name, amount, measurement)
		VALUES (?, ?, ?, ?, ?)
	`, nowMysql, nowMysql, resource.Name, resource.Amount, resource.Measurement)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	resource.ID = int(id)
	resource.CreatedAt = now
	resource.UpdatedAt = now

	return &resource, nil
}

func (impl *resourceStore) Update(ctx context.Context, resource model.Resource) (*model.Resource, error) {
	fields, values := getNotEmptyResourceFields(map[string]interface{}{
		"name":        resource.Name,
		"measurement": resource.Measurement,
	})
	if len(fields) == 0 {
		return nil, exception.NewEmptyModelException("resource")
	}

	query := fmt.Sprintf(`
		UPDATE resources
		SET updated_at = ?, %s
		WHERE id = ?
	`, strings.Join(fields, ", "))

	now := time.Now()
	t, err := impl.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	values = append([]interface{}{now.Format("2006-01-02T15:04:05")}, values...)
	values = append(values, resource.ID)

	res, err := t.ExecContext(ctx, query, values...)
	if err != nil {
		if err = t.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		if err = t.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	resS, err := t.QueryContext(ctx, `
		SELECT id,
			created_at,
			updated_at,
			name,
			amount,
			measurement
		FROM resources
		WHERE id = ?
		LIMIT 1
	`, resource.ID)
	if err != nil {
		return nil, err
	}

	var r *model.Resource
	for resS.Next() {
		r, err = scanResource(resS)
		if err != nil {
			return nil, err
		}
	}

	if err := t.Commit(); err != nil {
		return nil, err
	}

	return r, nil
}

func scanResource(res *sql.Rows) (*model.Resource, error) {
	var resource = &model.Resource{}
	var createdAt, updatedAt string

	if err := res.Scan(&resource.ID, &createdAt, &updatedAt, &resource.Name, &resource.Amount, &resource.Measurement); err != nil {
		return nil, err
	}

	t, err := time.Parse("2006-01-02T15:04:05", strings.Replace(createdAt, " ", "T", 1))
	if err != nil {
		return nil, err
	}
	resource.CreatedAt = t

	t, err = time.Parse("2006-01-02T15:04:05", strings.Replace(updatedAt, " ", "T", 1))
	if err != nil {
		return nil, err
	}
	resource.UpdatedAt = t

	return resource, nil
}

func getNotEmptyResourceFields(data map[string]interface{}) ([]string, []interface{}) {
	fields := []string{}
	values := []interface{}{}

	for field, value := range data {
		if value != "" {
			fields = append(fields, field+" = ?")
			values = append(values, value)
		}
	}

	return fields, values
}
