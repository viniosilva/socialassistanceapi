package store

import (
	"context"

	"github.com/viniosilva/socialassistanceapi/internal/configuration"
)

//go:generate mockgen -destination ../../mock/health_store_mock.go -package mock . HealthStore
type HealthStore interface {
	Health(ctx context.Context) bool
}

type healthStore struct {
	db configuration.MySQL
}

func NewHealthStore(db configuration.MySQL) HealthStore {
	return &healthStore{
		db: db,
	}
}

func (impl *healthStore) Health(ctx context.Context) bool {
	if err := impl.db.DB.PingContext(ctx); err != nil {
		return false
	}

	return true
}
