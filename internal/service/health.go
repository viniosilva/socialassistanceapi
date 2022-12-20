package service

import (
	"context"

	"github.com/sirupsen/logrus"
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
	log := logrus.WithFields(logrus.Fields{"span_id": ctx.Value("span_id"), "path": "internal.service.health.ping"})

	if err := impl.HealthRepository.Ping(ctx); err != nil {
		log.Error(err.Error())
		return HealthResponse{Status: HealthStatusDown}
	}

	return HealthResponse{HealthStatusUp}
}
