package service

import (
	"context"

	"github.com/viniosilva/socialassistanceapi/internal/repository"
)

//go:generate mockgen -destination ../../mock/health_service_mock.go -package mock . HealthService
type HealthService interface {
	Ping(ctx context.Context) HealthResponse
}

type HealthServiceImpl struct {
	HealthRepository repository.HealthRepository
}

func (impl *HealthServiceImpl) Ping(ctx context.Context) HealthResponse {
	if err := impl.HealthRepository.Ping(ctx); err != nil {
		return HealthResponse{Status: HealthStatusDown}
	}

	return HealthResponse{HealthStatusUp}
}
