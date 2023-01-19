package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	Gin                   *gin.Engine
	Addr                  string
	HealthService         service.HealthService
	PersonService         service.PersonService
	FamilyService         service.FamilyService
	ResourceService       service.ResourceService
	DonateResourceService service.DonateResourceService
}

// @title Ipanema Box API
// @version 1.0
// @description person, budget and service management
// @BasePath /api/v1
func (impl *ApiImpl) Configure() {
	api := gin.New()
	api.Use(cors.Default())
	api.Use(gin.Recovery())
	api.Use(impl.JSONLogMiddleware())

	docs.SwaggerInfo.Host = impl.Addr
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	healthApi := &HealthApiImpl{
		Router:          api.Group("/api/health"),
		HealthService:   impl.HealthService,
		TraceMiddleware: impl.TraceMiddleware,
	}
	personApi := &PersonApiImpl{
		Router:          api.Group("/api/v1/persons"),
		PersonService:   impl.PersonService,
		TraceMiddleware: impl.TraceMiddleware,
	}
	familyApi := &FamilyApiImpl{
		Router:          api.Group("/api/v1/families"),
		FamilyService:   impl.FamilyService,
		TraceMiddleware: impl.TraceMiddleware,
		Addr:            fmt.Sprintf("%s/api/v1/families", impl.Addr),
	}
	resourceApi := &ResourceApiImpl{
		Router:          api.Group("/api/v1/resources"),
		ResourceService: impl.ResourceService,
		TraceMiddleware: impl.TraceMiddleware,
	}
	donateResourceApi := &DonateResourceApiImpl{
		Router:                api.Group("/api/v1/resources"),
		DonateResourceService: impl.DonateResourceService,
		TraceMiddleware:       impl.TraceMiddleware,
	}

	healthApi.Configure()
	personApi.Configure()
	familyApi.Configure()
	resourceApi.Configure()
	donateResourceApi.Configure()

	impl.Gin = api
}

func (impl *ApiImpl) Start() {
	impl.Gin.Run(impl.Addr)
}

func (impl *ApiImpl) JSONLogMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := map[string]interface{}{
				"status":      params.StatusCode,
				"path":        params.Path,
				"method":      params.Method,
				"start":       params.TimeStamp.Format("2006-01-02T15:04:05Z07:00"),
				"remote_addr": params.ClientIP,
				"duration_ms": params.Latency.Milliseconds(),
				"trace_id":    params.Request.Header.Get("Trace-Id"),
				"span_id":     params.Request.Header.Get("Request-Id"),
			}

			if params.Request.Header.Get("Span-Id") != "" {
				log["parent_span_id"] = params.Request.Header.Get("Span-Id")
			}

			jsonLog, _ := json.Marshal(log)

			return string(jsonLog) + "\n"
		},
	)
}

func (impl *ApiImpl) TraceMiddleware(c *gin.Context) {
	traceID := c.Request.Header.Get("Trace-Id")
	if traceID == "" {
		traceID = strings.Replace(uuid.New().String(), "-", "", -1)[:32]
		c.Request.Header.Set("Trace-Id", traceID)
	}

	spanID := strings.Replace(uuid.New().String(), "-", "", -1)[:16]
	c.Request.Header.Set("Request-Id", spanID)

	c.Set("trace_id", c.Request.Header.Get("Trace-Id"))
	c.Set("span_id", c.Request.Header.Get("Request-Id"))
	c.Set("parent_span_id", c.Request.Header.Get("Span-Id"))

	c.Next()
}
