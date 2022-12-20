package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

//go:generate mockgen -destination ../../mock/health_api_mock.go -package mock . HealthApi
type HealthApi interface {
	Configure()
}

type HealthApiImpl struct {
	Router          *gin.RouterGroup
	HealthService   service.HealthService
	TraceMiddleware func(c *gin.Context)
}

func (impl *HealthApiImpl) Configure() {
	impl.Router.GET("", impl.TraceMiddleware, impl.HealthCheck)
}

func (impl *HealthApiImpl) HealthCheck(c *gin.Context) {
	res := impl.HealthService.Ping(c)
	c.JSON(http.StatusOK, res)
}
