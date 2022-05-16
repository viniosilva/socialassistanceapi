package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

type HealthApi struct {
	service *service.HealthService
}

func NewHealthApi(router *gin.RouterGroup, service *service.HealthService) *HealthApi {
	impl := &HealthApi{service}

	router.GET("", impl.Health)

	return impl
}

func (impl *HealthApi) Health(c *gin.Context) {
	res := impl.service.Health(c)
	c.JSON(http.StatusOK, res)
}
