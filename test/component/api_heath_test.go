package component

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/viniosilva/socialassistanceapi/internal/api"
	"github.com/viniosilva/socialassistanceapi/internal/configuration"
	"github.com/viniosilva/socialassistanceapi/internal/service"
	"github.com/viniosilva/socialassistanceapi/internal/store"
	"github.com/viniosilva/socialassistanceapi/mock"
)

func TestComponentHealthApiHealth(t *testing.T) {
	cases := map[string]struct {
		expectedCode int
		expectedBody service.Health
		prepareMock  func(mock *mock.MockHealthStore)
	}{
		"should return health status up": {
			expectedCode: 200,
			expectedBody: service.Health{Status: service.HealthStatusUp},
			prepareMock: func(mock *mock.MockHealthStore) {
				mock.EXPECT().Health(gomock.Any()).Return(true)
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			healthStore := store.NewHealthStore(mysql.DB)
			healthService := service.NewHealthService(healthStore)
			api := api.NewApi("0.0.0.0:8080", healthService, nil)

			// when
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/health", nil)
			api.Gin.ServeHTTP(rec, req)

			expectedBody, _ := json.Marshal(cs.expectedBody)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("GET /api/health StatusCode= %v, expected %v", rec.Code, cs.expectedCode)
			}
			if rec.Body.String() != string(expectedBody) {
				t.Errorf("GET /api/health Body= %v, expected %v", rec.Body.String(), string(expectedBody))
			}
		})
	}
}
