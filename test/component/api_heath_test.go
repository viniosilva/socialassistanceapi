package component

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/socialassistanceapi/internal/api"
	"github.com/viniosilva/socialassistanceapi/internal/configuration"
	"github.com/viniosilva/socialassistanceapi/internal/repository"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

func TestComponentHealthApiPing(t *testing.T) {
	cases := map[string]struct {
		expectedCode int
		expectedBody *service.HealthResponse
	}{
		"should return health status up": {
			expectedCode: http.StatusOK,
			expectedBody: &service.HealthResponse{Status: service.HealthStatusUp},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			cfg, err := configuration.LoadConfig("../..")
			if err != nil {
				log.Fatal("cannot load config: ", err)
			}

			mysql := configuration.MySQLConfigure(cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database, cfg.MySQL.Username,
				cfg.MySQL.Password, time.Duration(cfg.MySQL.ConnMaxLifetimeMs), cfg.MySQL.MaxOpenConns, cfg.MySQL.MaxIdleConns)
			defer mysql.DB.Close()

			healthRepository := &repository.HealthRepositoryImpl{DB: mysql}
			healthService := &service.HealthServiceImpl{HealthRepository: healthRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", HealthService: healthService}
			impl.Configure()

			// when
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/health", nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.HealthResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedBody, body)
		})
	}
}
