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

//go:generate mockgen -destination ../../mock/address_store_mock.go -package mock . AddressStore
type AddressStore interface {
	FindAll(ctx context.Context) ([]model.Address, error)
	FindOneById(ctx context.Context, addressID int) (*model.Address, error)
	Create(ctx context.Context, address model.Address) (*model.Address, error)
	Update(ctx context.Context, address model.Address) (*model.Address, error)
	Delete(ctx context.Context, addressID int) error
}

type addressStore struct {
	db *sql.DB
}

func NewAddressStore(db *sql.DB) AddressStore {
	return &addressStore{db}
}

func (impl *addressStore) FindAll(ctx context.Context) ([]model.Address, error) {
	addresses := []model.Address{}

	res, err := impl.db.Query(`
		SELECT id,
			created_at,
			updated_at,
			country,
			state,
			city,
			neighborhood,
			street,
			number,
			complement,
			zipcode
		FROM addresses
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		address, err := scanAddress(res)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, *address)
	}

	return addresses, nil
}

func (impl *addressStore) FindOneById(ctx context.Context, addressID int) (*model.Address, error) {
	res, err := impl.db.QueryContext(ctx, `
		SELECT id,
			created_at,
			updated_at,
			country,
			state,
			city,
			neighborhood,
			street,
			number,
			complement,
			zipcode
		FROM addresses
		WHERE id = ?
		LIMIT 1
	`, addressID)
	if err != nil {
		return nil, err
	}

	var address *model.Address
	for res.Next() {
		address, err = scanAddress(res)
		if err != nil {
			return nil, err
		}
	}

	return address, nil
}

func (impl *addressStore) Create(ctx context.Context, address model.Address) (*model.Address, error) {
	now := time.Now()
	nowMysql := now.Format("2006-01-02T15:04:05")
	res, err := impl.db.ExecContext(ctx, `
		INSERT INTO addresses (created_at, updated_at, country,
			state, city, neighborhood, street, number, complement, zipcode)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, nowMysql, nowMysql, address.Country, address.State, address.City, address.Neighborhood,
		address.Street, address.Number, address.Complement, address.Zipcode)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	address.ID = int(id)
	address.CreatedAt = now
	address.UpdatedAt = now

	return &address, nil
}

func (impl *addressStore) Update(ctx context.Context, address model.Address) (*model.Address, error) {
	fields, values := getNotEmptyAddressFields(address)
	if len(fields) == 0 {
		return nil, exception.NewEmptyModelException("address")
	}

	query := fmt.Sprintf(`
		UPDATE addresses
		SET updated_at = ?, %s
		WHERE id = ?
	`, strings.Join(fields, ", "))

	now := time.Now()
	t, err := impl.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	values = append([]interface{}{now.Format("2006-01-02T15:04:05")}, values...)
	values = append(values, address.ID)

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
			country,
			state,
			city,
			neighborhood,
			street,
			number,
			complement,
			zipcode
		FROM addresses
		WHERE id = ?
		LIMIT 1
	`, address.ID)
	if err != nil {
		return nil, err
	}

	var a *model.Address
	for resS.Next() {
		a, err = scanAddress(resS)
		if err != nil {
			return nil, err
		}
	}

	if err := t.Commit(); err != nil {
		return nil, err
	}

	return a, nil
}

func (impl *addressStore) Delete(ctx context.Context, addressID int) error {
	_, err := impl.db.ExecContext(ctx, `
		UPDATE addresses
		SET deleted_at = NOW()
		WHERE id = ?
	`, addressID)

	return err
}

func scanAddress(res *sql.Rows) (*model.Address, error) {
	var address = &model.Address{}
	var createdAt, updatedAt string

	if err := res.Scan(&address.ID, &createdAt, &updatedAt, &address.Country,
		&address.State, &address.City, &address.Neighborhood, &address.Street,
		&address.Number, &address.Complement, &address.Zipcode); err != nil {
		return nil, err
	}

	t, err := time.Parse("2006-01-02T15:04:05", strings.Replace(createdAt, " ", "T", 1))
	if err != nil {
		return nil, err
	}
	address.CreatedAt = t

	t, err = time.Parse("2006-01-02T15:04:05", strings.Replace(updatedAt, " ", "T", 1))
	if err != nil {
		return nil, err
	}
	address.UpdatedAt = t

	return address, nil
}

func getNotEmptyAddressFields(address model.Address) ([]string, []interface{}) {
	fields := []string{}
	values := []interface{}{}

	all := map[string]string{
		"country":      address.Country,
		"state":        address.State,
		"city":         address.City,
		"neighborhood": address.Neighborhood,
		"street":       address.Street,
		"number":       address.Number,
		"complement":   address.Complement,
		"zipcode":      address.Zipcode,
	}

	for field, value := range all {
		if value != "" {
			fields = append(fields, field+" = ?")
			values = append(values, value)
		}
	}

	return fields, values
}
