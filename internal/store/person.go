package store

import (
	"context"
	"database/sql"
	"strings"
	"time"

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
	db *sql.DB
}

func NewPersonStore(db *sql.DB) PersonStore {
	return &personStore{db}
}

func (impl *personStore) FindAll(ctx context.Context) ([]model.Person, error) {
	people := []model.Person{}

	res, err := impl.db.Query(`
		SELECT id,
			created_at,
			updated_at,
			name
		FROM people
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

		people = append(people, *person)
	}

	return people, nil
}

func (impl *personStore) FindOneById(ctx context.Context, personID int) (*model.Person, error) {
	res, err := impl.db.QueryContext(ctx, `
		SELECT id,
			created_at,
			updated_at,
			name
		FROM people
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

	return person, nil
}

func (impl *personStore) Create(ctx context.Context, person model.Person) (*model.Person, error) {
	now := time.Now()
	nowMysql := now.Format("2006-01-02T15:04:05")
	res, err := impl.db.ExecContext(ctx, `
		INSERT INTO people (created_at, updated_at, name)
		VALUES (?, ?, ?)
	`, nowMysql, nowMysql, person.Name)
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
	now := time.Now()
	t, err := impl.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	res, err := t.ExecContext(ctx, `
		UPDATE people
		SET name = ?,
			updated_at = ?
		WHERE id = ?
	`, person.Name, now.Format("2006-01-02T15:04:05"), person.ID)

	if err != nil {
		t.Rollback()
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		t.Rollback()
		return nil, err
	}

	var createdAt string
	impl.db.QueryRowContext(ctx, `
		SELECT created_at
		FROM people
		WHERE id = ?
		LIMIT 1
	`, person.ID).Scan(&createdAt)

	if err := t.Commit(); err != nil {
		return nil, err
	}

	c, err := time.Parse("2006-01-02T15:04:05", strings.Replace(createdAt, " ", "T", 1))
	if err != nil {
		return nil, err
	}

	person.CreatedAt = c
	person.UpdatedAt = now
	return &person, nil
}

func (impl *personStore) Delete(ctx context.Context, personID int) error {
	_, err := impl.db.ExecContext(ctx, `
		UPDATE people
		SET deleted_at = NOW()
		WHERE id = ?
	`, personID)

	return err
}

func scanPerson(res *sql.Rows) (*model.Person, error) {
	var person = &model.Person{}
	var createdAt, updatedAt string

	if err := res.Scan(&person.ID, &createdAt, &updatedAt, &person.Name); err != nil {
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
