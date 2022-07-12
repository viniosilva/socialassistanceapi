package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/viniosilva/socialassistanceapi/internal/configuration"
	"github.com/viniosilva/socialassistanceapi/internal/exception"
	"github.com/viniosilva/socialassistanceapi/internal/model"
)

//go:generate mockgen -destination ../../mock/person_repository_mock.go -package mock . PersonRepository
type PersonRepository interface {
	FindAll(ctx context.Context) ([]model.Person, error)
	FindOneById(ctx context.Context, personID int) (*model.Person, error)
	Create(ctx context.Context, person model.Person) (*model.Person, error)
	Update(ctx context.Context, person model.Person) error
	Delete(ctx context.Context, personID int) error
}

type PersonRepositoryImpl struct {
	DB configuration.MySQL
}

func (impl *PersonRepositoryImpl) FindAll(ctx context.Context) ([]model.Person, error) {
	persons := []model.Person{}

	res, err := impl.DB.DB.Query(`
		SELECT id,
			created_at,
			updated_at,
			address_id,
			name
		FROM persons
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		person, err := impl.ScanPerson(res)
		if err != nil {
			return nil, err
		}

		persons = append(persons, *person)
	}

	return persons, nil
}

func (impl *PersonRepositoryImpl) FindOneById(ctx context.Context, personID int) (*model.Person, error) {
	res, err := impl.DB.DB.QueryContext(ctx, `
		SELECT id,
			created_at,
			updated_at,
			address_id,
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
		person, err = impl.ScanPerson(res)
		if err != nil {
			return nil, err
		}
	}

	if person == nil {
		return nil, &exception.NotFoundException{Err: fmt.Errorf("person %d not found", personID)}
	}

	return person, nil
}

func (impl *PersonRepositoryImpl) Create(ctx context.Context, person model.Person) (*model.Person, error) {
	now := time.Now()
	nowMysql := now.Format("2006-01-02T15:04:05")
	res, err := impl.DB.DB.ExecContext(ctx, `
		INSERT INTO persons (created_at, updated_at, address_id, name)
		VALUES (?, ?, ?, ?)
	`, nowMysql, nowMysql, person.AddressID, person.Name)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	person.ID = int(id)
	person.CreatedAt = now
	person.UpdatedAt = now

	return &person, nil
}

func (impl *PersonRepositoryImpl) Update(ctx context.Context, person model.Person) error {
	fields, values := impl.DB.BuildUpdateData(map[string]interface{}{
		"name": person.Name,
	})
	if len(fields) == 0 {
		return &exception.EmptyModelException{Err: fmt.Errorf("empty person model")}
	}

	if person.AddressID > 0 {
		fields = append(fields, "address_id")
		values = append(values, person.AddressID)
	}

	query := fmt.Sprintf(`
		UPDATE persons
		SET updated_at = ?, %s
		WHERE id = ?
	`, strings.Join(fields, ", "))

	now := time.Now()
	values = append([]interface{}{now.Format("2006-01-02T15:04:05")}, values...)
	values = append(values, person.ID)

	res, err := impl.DB.DB.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return &exception.NotFoundException{Err: fmt.Errorf("person %d not found", person.ID)}
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

func (impl *PersonRepositoryImpl) ScanPerson(res *sql.Rows) (*model.Person, error) {
	var person = &model.Person{}
	var createdAt, updatedAt string

	if err := res.Scan(&person.ID, &createdAt, &updatedAt, &person.AddressID, &person.Name); err != nil {
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
