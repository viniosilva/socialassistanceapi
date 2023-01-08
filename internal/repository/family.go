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

//go:generate mockgen -destination ../../mock/family_repository_mock.go -package mock . FamilyRepository
type FamilyRepository interface {
	FindAll(ctx context.Context) ([]model.Family, error)
	FindOneById(ctx context.Context, familyID int) (*model.Family, error)
	Create(ctx context.Context, data model.Family) (*model.Family, error)
	Update(ctx context.Context, data model.Family) error
	Delete(ctx context.Context, data int) error
}

type FamilyRepositoryImpl struct {
	DB infra.MySQL
}

func (impl *FamilyRepositoryImpl) FindAll(ctx context.Context) ([]model.Family, error) {
	data := []model.Family{}

	res, err := impl.DB.DB.Query(`
		SELECT id,
			created_at,
			updated_at,
			name,
			country,
			state,
			city,
			neighborhood,
			street,
			number,
			complement,
			zipcode
		FROM families
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		d, err := impl.Scan(res)
		if err != nil {
			return nil, err
		}

		data = append(data, *d)
	}

	return data, nil
}

func (impl *FamilyRepositoryImpl) FindOneById(ctx context.Context, familyID int) (*model.Family, error) {
	res, err := impl.DB.DB.QueryContext(ctx, `
		SELECT id,
			created_at,
			updated_at,
			name,
			country,
			state,
			city,
			neighborhood,
			street,
			number,
			complement,
			zipcode
		FROM families
		WHERE id = ?
		LIMIT 1
	`, familyID)
	if err != nil {
		return nil, err
	}

	var data *model.Family
	for res.Next() {
		data, err = impl.Scan(res)
		if err != nil {
			return nil, err
		}
	}

	if data == nil {
		return nil, &exception.NotFoundException{Err: fmt.Errorf("family %d not found", familyID)}
	}

	return data, nil
}

func (impl *FamilyRepositoryImpl) Create(ctx context.Context, data model.Family) (*model.Family, error) {
	now := time.Now()
	nowMysql := now.Format("2006-01-02T15:04:05")
	res, err := impl.DB.DB.ExecContext(ctx, `
		INSERT INTO families (created_at, updated_at, name, country,
			state, city, neighborhood, street, number, complement, zipcode)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, nowMysql, nowMysql, data.Name, data.Country, data.State, data.City,
		data.Neighborhood, data.Street, data.Number, data.Complement, data.Zipcode)
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

func (impl *FamilyRepositoryImpl) Update(ctx context.Context, data model.Family) error {
	fields, values := impl.DB.BuildUpdateData(map[string]interface{}{
		"name":         data.Name,
		"country":      data.Country,
		"state":        data.State,
		"city":         data.City,
		"neighborhood": data.Neighborhood,
		"street":       data.Street,
		"number":       data.Number,
		"complement":   data.Complement,
		"zipcode":      data.Zipcode,
	})
	if len(fields) == 0 {
		return &exception.EmptyModelException{Err: fmt.Errorf("empty family model")}
	}

	query := fmt.Sprintf(`
		UPDATE families
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
		return &exception.NotFoundException{Err: fmt.Errorf("family %d not found", data.ID)}
	}

	return nil
}

func (impl *FamilyRepositoryImpl) Delete(ctx context.Context, familyID int) error {
	_, err := impl.DB.DB.ExecContext(ctx, `
		UPDATE families
		SET deleted_at = NOW()
		WHERE id = ?
	`, familyID)

	return err
}

func (impl *FamilyRepositoryImpl) Scan(res *sql.Rows) (*model.Family, error) {
	var data = &model.Family{}
	var createdAt, updatedAt string

	if err := res.Scan(&data.ID, &createdAt, &updatedAt, &data.Name, &data.Country,
		&data.State, &data.City, &data.Neighborhood, &data.Street, &data.Number,
		&data.Complement, &data.Zipcode); err != nil {
		return nil, err
	}

	t, err := time.Parse("2006-01-02T15:04:05", strings.Replace(createdAt, " ", "T", 1))
	if err != nil {
		return nil, err
	}
	data.CreatedAt = t

	t, err = time.Parse("2006-01-02T15:04:05", strings.Replace(updatedAt, " ", "T", 1))
	if err != nil {
		return nil, err
	}
	data.UpdatedAt = t

	return data, nil
}
