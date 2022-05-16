package service

import (
	"context"

	"github.com/viniosilva/socialassistanceapi/internal/store"
)

type HealthService struct {
	store store.HealthStore
}

func NewHealthService(store store.HealthStore) *HealthService {
	return &HealthService{store}
}

func (impl *HealthService) Health(ctx context.Context) Health {
	health := impl.store.Health(ctx)
	if !health {
		return Health{Status: HealthStatusDown}
	}

	return Health{HealthStatusUp}
}
