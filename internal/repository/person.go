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

//go:generate mockgen -destination ../../mock/person_repository_mock.go -package mock . PersonRepository
type PersonRepository interface {
	FindAll(ctx context.Context) ([]model.Person, error)
	FindOneById(ctx context.Context, personID int) (*model.Person, error)
	Create(ctx context.Context, data model.Person) (*model.Person, error)
	Update(ctx context.Context, data model.Person) error
	Delete(ctx context.Context, personID int) error
}

type PersonRepositoryImpl struct {
	DB infra.MySQL
}

func (impl *PersonRepositoryImpl) FindAll(ctx context.Context) ([]model.Person, error) {
	data := []model.Person{}

	res, err := impl.DB.DB.Query(`
		SELECT id,
			created_at,
			updated_at,
			family_id,
			name
		FROM persons
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		person, err := impl.Scan(res)
		if err != nil {
			return nil, err
		}

		data = append(data, *person)
	}

	return data, nil
}

func (impl *PersonRepositoryImpl) FindOneById(ctx context.Context, personID int) (*model.Person, error) {
	res, err := impl.DB.DB.QueryContext(ctx, `
		SELECT id,
			created_at,
			updated_at,
			family_id,
			name
		FROM persons
		WHERE id = ?
		LIMIT 1
	`, personID)
	if err != nil {
		return nil, err
	}

	var person *model.Person
	for res.Next() {
		person, err = impl.Scan(res)
		if err != nil {
			return nil, err
		}
	}

	if person == nil {
		return nil, &exception.NotFoundException{Err: fmt.Errorf("person %d not found", personID)}
	}

	return person, nil
}

func (impl *PersonRepositoryImpl) Create(ctx context.Context, data model.Person) (*model.Person, error) {
	now := time.Now()
	nowMysql := now.Format("2006-01-02T15:04:05")
	res, err := impl.DB.DB.ExecContext(ctx, `
		INSERT INTO persons (created_at, updated_at, family_id, name)
		VALUES (?, ?, ?, ?)
	`, nowMysql, nowMysql, data.FamilyID, data.Name)
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

func (impl *PersonRepositoryImpl) Update(ctx context.Context, data model.Person) error {
	fields, values := impl.DB.BuildUpdateData(map[string]interface{}{
		"name": data.Name,
	})
	if len(fields) == 0 {
		return &exception.EmptyModelException{Err: fmt.Errorf("empty person model")}
	}

	if data.FamilyID > 0 {
		fields = append(fields, "family_id")
		values = append(values, data.FamilyID)
	}

	query := fmt.Sprintf(`
		UPDATE persons
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
		return &exception.NotFoundException{Err: fmt.Errorf("person %d not found", data.ID)}
	}

	return nil
}

func (impl *PersonRepositoryImpl) Delete(ctx context.Context, personID int) error {
	_, err := impl.DB.DB.ExecContext(ctx, `
		UPDATE persons
		SET deleted_at = NOW()
		WHERE id = ?
	`, personID)

	return err
}

func (impl *PersonRepositoryImpl) Scan(res *sql.Rows) (*model.Person, error) {
	var person = &model.Person{}
	var createdAt, updatedAt string

	if err := res.Scan(&person.ID, &createdAt, &updatedAt, &person.FamilyID, &person.Name); err != nil {
		return nil, err
	}

	t, err := time.Parse("2006-01-02T15:04:05", strings.Replace(createdAt, " ", "T", 1))
	if err != nil {
		return nil, err
	}
	person.CreatedAt = t

	t, err = time.Parse("2006-01-02T15:04:05", strings.Replace(updatedAt, " ", "T", 1))
	if err != nil {
		return nil, err
	}
	person.UpdatedAt = t

	return person, nil
}
