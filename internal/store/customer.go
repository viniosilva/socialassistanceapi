package store

import (
	"context"
	"database/sql"
	"time"

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
			created_at,
			updated_at,
			name
		FROM customers
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		customer, err := scanCustomer(res)
		if err != nil {
			return nil, err
		}

		customers = append(customers, *customer)
	}

	return customers, nil
}

func (impl *customerStore) FindOneById(ctx context.Context, customerID int) (*model.Customer, error) {
	res, err := impl.db.QueryContext(ctx, `
		SELECT id,
			created_at,
			updated_at,
			name
		FROM customers
		WHERE id = ?
		LIMIT 1
	`, customerID)
	if err != nil {
		return nil, err
	}

	var customer *model.Customer
	for res.Next() {
		customer, err = scanCustomer(res)
		if err != nil {
			return nil, err
		}
	}

	return customer, nil
}

func (impl *customerStore) Create(ctx context.Context, customer model.Customer) (*model.Customer, error) {
	now := time.Now()
	nowMysql := now.Format("2006-01-02")
	res, err := impl.db.ExecContext(ctx, `
		INSERT INTO customers (created_at, updated_at, name)
		VALUES (?, ?, ?)
	`, nowMysql, nowMysql, customer.Name)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	customer.ID = int(id)
	customer.CreatedAt = now
	customer.UpdatedAt = now

	return &customer, nil
}

func (impl *customerStore) Update(ctx context.Context, customer model.Customer) (*model.Customer, error) {
	now := time.Now()
	t, err := impl.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	res, err := t.ExecContext(ctx, `
		UPDATE customers
		SET name = ?,
			updated_at = ?
		WHERE id = ?
	`, customer.Name, now.Format("2006-01-02"), customer.ID)

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
		FROM customers
		WHERE id = ?
		LIMIT 1
	`, customer.ID).Scan(&createdAt)

	if err := t.Commit(); err != nil {
		return nil, err
	}

	c, err := time.Parse("2006-01-02", createdAt)
	if err != nil {
		return nil, err
	}

	customer.CreatedAt = c
	customer.UpdatedAt = now
	return &customer, nil
}

func scanCustomer(res *sql.Rows) (*model.Customer, error) {
	var customer = &model.Customer{}
	var createdAt, updatedAt string

	if err := res.Scan(&customer.ID, &createdAt, &updatedAt, &customer.Name); err != nil {
		return nil, err
	}

	t, err := time.Parse("2006-01-02", createdAt)
	if err != nil {
		return nil, err
	}
	customer.CreatedAt = t

	t, err = time.Parse("2006-01-02", updatedAt)
	if err != nil {
		return nil, err
	}
	customer.UpdatedAt = t

	return customer, nil
}
