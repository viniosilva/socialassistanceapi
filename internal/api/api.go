package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/viniosilva/socialassistanceapi/docs"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

type Api struct {
	Gin  *gin.Engine
	addr string
}

// @title Ipanema Box API
// @version 1.0
// @description Customer, budget and service management
// @BasePath /api/v1
func NewApi(addr string, healthService *service.HealthService, customerService *service.CustomerService) *Api {
	api := gin.Default()
	api.Use(cors.Default())

	docs.SwaggerInfo.Host = addr
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	NewHealthApi(api.Group("/api/health"), healthService)
	NewCustomerApi(api.Group("/api/v1/customers"), customerService)

	return &Api{api, addr}
}

func (impl *Api) Start() {
	impl.Gin.Run(impl.addr)
}
