package repository

import (
	"context"

	"github.com/viniosilva/socialassistanceapi/internal/configuration"
)

//go:generate mockgen -destination ../../mock/health_repository_mock.go -package mock . HealthRepository
type HealthRepository interface {
	Ping(ctx context.Context) error
}

type HealthRepositoryImpl struct {
	DB configuration.MySQL
}

func (impl *HealthRepositoryImpl) Ping(ctx context.Context) error {
	return impl.DB.DB.PingContext(ctx)
}
