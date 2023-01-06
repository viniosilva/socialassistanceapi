package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/viniosilva/socialassistanceapi/internal/exception"
	"github.com/viniosilva/socialassistanceapi/internal/infra"
	"github.com/viniosilva/socialassistanceapi/internal/model"
)

//go:generate mockgen -destination ../../mock/resource_repository_mock.go -package mock . ResourceRepository
type ResourceRepository interface {
	FindAll(ctx context.Context) ([]model.Resource, error)
	FindOneById(ctx context.Context, resourceID int) (*model.Resource, error)
	Create(ctx context.Context, data model.Resource) (*model.Resource, error)
	Update(ctx context.Context, data model.Resource) error
	UpdateQuantity(ctx context.Context, resourceID int, quantity float64) error
}

type ResourceRepositoryImpl struct {
	DB infra.MySQL
}

func (impl *ResourceRepositoryImpl) FindAll(ctx context.Context) ([]model.Resource, error) {
	data := []model.Resource{}

	res, err := impl.DB.DB.Query(`
		SELECT id,
			created_at,
			updated_at,
			name,
			amount,
			measurement,
			quantity
		FROM resources`)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		resource, err := impl.Scan(res)
		if err != nil {
			return nil, err
		}
		data = append(data, *resource)
	}
	return data, nil
}

func (impl *ResourceRepositoryImpl) FindOneById(ctx context.Context, resourceID int) (*model.Resource, error) {
	res, err := impl.DB.DB.QueryContext(ctx, `
		SELECT id,
			created_at,
			updated_at,
			name,
			amount,
			measurement,
			quantity
		FROM resources
		WHERE id = ?
		LIMIT 1 `, resourceID)
	if err != nil {
		return nil, err
	}

	var data *model.Resource
	for res.Next() {
		data, err = impl.Scan(res)
		if err != nil {
			return nil, err
		}
	}

	if data == nil {
		return nil, &exception.NotFoundException{Err: fmt.Errorf("resource %d not found", resourceID)}
	}

	return data, nil
}

func (impl *ResourceRepositoryImpl) Create(ctx context.Context, data model.Resource) (*model.Resource, error) {
	now := time.Now()
	nowMysql := now.Format("2006-01-02T15:04:05")
	res, err := impl.DB.DB.ExecContext(ctx, `
		INSERT INTO resources (created_at, updated_at, name, amount, measurement, quantity)
		VALUES (?, ?, ?, ?, ?, ?)
	`, nowMysql, nowMysql, data.Name, data.Amount, data.Measurement, data.Quantity)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	data.ID = int(id)
	data.CreatedAt = now
	data.UpdatedAt = now

	return &data, nil
}

func (impl *ResourceRepositoryImpl) Update(ctx context.Context, data model.Resource) error {
	fields, values := impl.DB.BuildUpdateData(map[string]interface{}{
		"name":        data.Name,
		"amount":      data.Amount,
		"measurement": data.Measurement,
	})
	if len(fields) == 0 {
		return &exception.EmptyModelException{Err: fmt.Errorf("empty resource model")}
	}

	query := fmt.Sprintf(`
		UPDATE resources
		SET updated_at = ?, %s
		WHERE id = ?
	`, strings.Join(fields, ", "))

	now := time.Now()

	values = append([]interface{}{now.Format("2006-01-02T15:04:05")}, values...)
	values = append(values, data.ID)

	res, err := impl.DB.DB.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return &exception.NotFoundException{Err: fmt.Errorf("resource %d not found", data.ID)}
	}

	return nil
}

func (impl *ResourceRepositoryImpl) UpdateQuantity(ctx context.Context, resourceID int, quantity float64) error {
	query := `
		UPDATE resources
		SET updated_at = ?,
			quantity = ?
		WHERE id = ?
	`

	res, err := impl.DB.DB.ExecContext(ctx, query, time.Now().Format("2006-01-02T15:04:05"), quantity, resourceID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return &exception.NotFoundException{Err: fmt.Errorf("resource %d not found", resourceID)}
	}

	return nil
}

func (impl *ResourceRepositoryImpl) Scan(res *sql.Rows) (*model.Resource, error) {
	var resource = &model.Resource{}
	var createdAt, updatedAt string

	if err := res.Scan(&resource.ID, &createdAt, &updatedAt, &resource.Name,
		&resource.Amount, &resource.Measurement, &resource.Quantity); err != nil {

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
