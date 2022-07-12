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

//go:generate mockgen -destination ../../mock/address_repository_mock.go -package mock . AddressRepository
type AddressRepository interface {
	FindAll(ctx context.Context) ([]model.Address, error)
	FindOneById(ctx context.Context, addressID int) (*model.Address, error)
	Create(ctx context.Context, address model.Address) (*model.Address, error)
	Update(ctx context.Context, address model.Address) error
	Delete(ctx context.Context, addressID int) error
}

type AddressRepositoryImpl struct {
	DB configuration.MySQL
}

func (impl *AddressRepositoryImpl) FindAll(ctx context.Context) ([]model.Address, error) {
	addresses := []model.Address{}

	res, err := impl.DB.DB.Query(`
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
		address, err := impl.ScanAddress(res)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, *address)
	}

	return addresses, nil
}

func (impl *AddressRepositoryImpl) FindOneById(ctx context.Context, addressID int) (*model.Address, error) {
	res, err := impl.DB.DB.QueryContext(ctx, `
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
		address, err = impl.ScanAddress(res)
		if err != nil {
			return nil, err
		}
	}

	if address == nil {
		return nil, &exception.NotFoundException{Err: fmt.Errorf("address %d not found", addressID)}
	}

	return address, nil
}

func (impl *AddressRepositoryImpl) Create(ctx context.Context, address model.Address) (*model.Address, error) {
	now := time.Now()
	nowMysql := now.Format("2006-01-02T15:04:05")
	res, err := impl.DB.DB.ExecContext(ctx, `
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

func (impl *AddressRepositoryImpl) Update(ctx context.Context, address model.Address) error {
	fields, values := impl.DB.BuildUpdateData(map[string]interface{}{
		"country":      address.Country,
		"state":        address.State,
		"city":         address.City,
		"neighborhood": address.Neighborhood,
		"street":       address.Street,
		"number":       address.Number,
		"complement":   address.Complement,
		"zipcode":      address.Zipcode,
	})
	if len(fields) == 0 {
		return &exception.EmptyModelException{Err: fmt.Errorf("empty address model")}
	}

	query := fmt.Sprintf(`
		UPDATE addresses
		SET updated_at = ?, %s
		WHERE id = ?
	`, strings.Join(fields, ", "))

	now := time.Now()

	values = append([]interface{}{now.Format("2006-01-02T15:04:05")}, values...)
	values = append(values, address.ID)

	res, err := impl.DB.DB.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return &exception.NotFoundException{Err: fmt.Errorf("address %d not found", address.ID)}
	}

	return nil
}

func (impl *AddressRepositoryImpl) Delete(ctx context.Context, addressID int) error {
	_, err := impl.DB.DB.ExecContext(ctx, `
		UPDATE addresses
		SET deleted_at = NOW()
		WHERE id = ?
	`, addressID)

	return err
}

func (impl *AddressRepositoryImpl) ScanAddress(res *sql.Rows) (*model.Address, error) {
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
