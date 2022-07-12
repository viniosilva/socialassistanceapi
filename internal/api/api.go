package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/viniosilva/socialassistanceapi/docs"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

//go:generate mockgen -destination ../../mock/api_mock.go -package mock . Api
type Api interface {
	Configure()
	Start()
}

type ApiImpl struct {
	Gin             *gin.Engine
	Addr            string
	HealthService   service.HealthService
	PersonService   service.PersonService
	AddressService  service.AddressService
	ResourceService service.ResourceService
}

// @title Ipanema Box API
// @version 1.0
// @description person, budget and service management
// @BasePath /api/v1
func (impl *ApiImpl) Configure() {
	api := gin.Default()
	api.Use(cors.Default())

	docs.SwaggerInfo.Host = impl.Addr
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	healthApi := &HealthApiImpl{Router: api.Group("/api/health"), HealthService: impl.HealthService}
	personApi := &PersonApiImpl{Router: api.Group("/api/v1/persons"), PersonService: impl.PersonService}
	addressApi := &AddressApiImpl{Router: api.Group("/api/v1/addresses"), AddressService: impl.AddressService}
	resourceApi := &ResourceApiImpl{Router: api.Group("/api/v1/resources"), ResourceService: impl.ResourceService}

	healthApi.Configure()
	personApi.Configure()
	addressApi.Configure()
	resourceApi.Configure()

	impl.Gin = api
}

func (impl *ApiImpl) Start() {
	impl.Gin.Run(impl.Addr)
}
