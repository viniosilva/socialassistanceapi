package store

import (
	"context"
	"database/sql"
)

//go:generate mockgen -destination ../../mock/health_store_mock.go -package mock . HealthStore
type HealthStore interface {
	Health(ctx context.Context) bool
}

type healthStore struct {
	db *sql.DB
}

func NewHealthStore(db *sql.DB) HealthStore {
	return &healthStore{db}
}

func (impl *healthStore) Health(ctx context.Context) bool {
	if err := impl.db.PingContext(ctx); err != nil {
		return false
	}

	return true
}
