package store

import (
	"context"
	"database/sql"

	"github.com/viniosilva/socialassistanceapi/internal/model"
)

//go:generate mockgen -destination ../../mock/customer_store_mock.go -package mock . CustomerStore
type CustomerStore interface {
	FindAll(ctx context.Context) ([]model.Customer, error)
	FindOneById(ctx context.Context, customerID int) (*model.Customer, error)
	Create(ctx context.Context, customer model.Customer) (*model.Customer, error)
	Update(ctx context.Context, customer model.Customer) (*model.Customer, error)
}

type customerStore struct {
	db *sql.DB
}

func NewCustomerStore(db *sql.DB) CustomerStore {
	return &customerStore{db}
}

func (impl *customerStore) FindAll(ctx context.Context) ([]model.Customer, error) {
	customers := []model.Customer{}

	res, err := impl.db.Query(`
		SELECT id,
			name
		FROM customers;
	`)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var customer model.Customer
		if err = res.Scan(&customer.ID, &customer.Name); err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

func (impl *customerStore) FindOneById(ctx context.Context, customerID int) (*model.Customer, error) {
	var customer = &model.Customer{}

	res, err := impl.db.QueryContext(ctx, `
		SELECT id,
			name
		FROM customers
		WHERE id = ?
		LIMIT 1;
	`, customerID)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		if err = res.Scan(&customer.ID, &customer.Name); err != nil {
			return nil, err
		}
	}

	if customer.ID == 0 {
		return nil, nil
	}
	return customer, nil
}

func (impl *customerStore) Create(ctx context.Context, customer model.Customer) (*model.Customer, error) {
	res, err := impl.db.ExecContext(ctx, "INSERT INTO customers (name) VALUES (?)", customer.Name)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	customer.ID = int(id)

	return &customer, nil
}

func (impl *customerStore) Update(ctx context.Context, customer model.Customer) (*model.Customer, error) {
	res, err := impl.db.ExecContext(ctx, `
		UPDATE customers
		SET name = ?
		WHERE id = ?
	`, customer.Name, customer.ID)

	if err != nil {
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return nil, err
	}

	return &customer, nil
}
