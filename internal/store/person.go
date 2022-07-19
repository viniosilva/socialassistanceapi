package store

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

//go:generate mockgen -destination ../../mock/person_store_mock.go -package mock . PersonStore
type PersonStore interface {
	FindAll(ctx context.Context) ([]model.Person, error)
	FindOneById(ctx context.Context, personID int) (*model.Person, error)
	Create(ctx context.Context, person model.Person) (*model.Person, error)
	Update(ctx context.Context, person model.Person) (*model.Person, error)
	Delete(ctx context.Context, personID int) error
}

type personStore struct {
	db configuration.MySQL
}

func NewPersonStore(db configuration.MySQL) PersonStore {
	return &personStore{
		db: db,
	}
}

func (impl *personStore) FindAll(ctx context.Context) ([]model.Person, error) {
	persons := []model.Person{}

	res, err := impl.db.DB.Query(`
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
		person, err := scanPerson(res)
		if err != nil {
			return nil, err
		}

		persons = append(persons, *person)
	}

	return persons, nil
}

func (impl *personStore) FindOneById(ctx context.Context, personID int) (*model.Person, error) {
	res, err := impl.db.DB.QueryContext(ctx, `
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
		person, err = scanPerson(res)
		if err != nil {
			return nil, err
		}
	}

	if person == nil {
		return nil, exception.NewNotFoundException("person")
	}

	return person, nil
}

func (impl *personStore) Create(ctx context.Context, person model.Person) (*model.Person, error) {
	now := time.Now()
	nowMysql := now.Format("2006-01-02T15:04:05")
	res, err := impl.db.DB.ExecContext(ctx, `
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

func (impl *personStore) Update(ctx context.Context, person model.Person) (*model.Person, error) {
	fields, values := impl.db.BuildUpdateData(map[string]interface{}{
		"name": person.Name,
	})
	if len(fields) == 0 {
		return nil, exception.NewEmptyModelException("person")
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

	t, err := impl.db.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	values = append([]interface{}{now.Format("2006-01-02T15:04:05")}, values...)
	values = append(values, person.ID)

	res, err := t.ExecContext(ctx, query, values...)
	if err != nil {
		if err := t.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		if err := t.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if rows == 0 {
		if err := t.Rollback(); err != nil {
			return nil, err
		}
		return nil, exception.NewNotFoundException("person")
	}

	resS, err := t.QueryContext(ctx, `
		SELECT id,
			created_at,
			updated_at,
			address_id,
			name
		FROM persons
		WHERE id = ?
		LIMIT 1
	`, person.ID)
	if err != nil {
		if err := t.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	var p *model.Person
	for resS.Next() {
		p, err = scanPerson(resS)
		if err != nil {
			if err := t.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	if err := t.Commit(); err != nil {
		if err := t.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	return p, nil
}

func (impl *personStore) Delete(ctx context.Context, personID int) error {
	_, err := impl.db.DB.ExecContext(ctx, `
		UPDATE persons
		SET deleted_at = NOW()
		WHERE id = ?
	`, personID)

	return err
}

func scanPerson(res *sql.Rows) (*model.Person, error) {
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
