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
	Router        *gin.RouterGroup
	HealthService service.HealthService
}

func (impl *HealthApiImpl) Configure() {
	impl.Router.GET("", impl.HealthCheck)
}

func (impl *HealthApiImpl) HealthCheck(c *gin.Context) {
	res := impl.HealthService.Ping(c)
	c.JSON(http.StatusOK, res)
}
