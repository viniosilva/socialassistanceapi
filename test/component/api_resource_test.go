package component

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/socialassistanceapi/internal/api"
	"github.com/viniosilva/socialassistanceapi/internal/configuration"
	"github.com/viniosilva/socialassistanceapi/internal/infra"
	"github.com/viniosilva/socialassistanceapi/internal/repository"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

func Test_ResourceApi_FindAll(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"

	cases := map[string]struct {
		before       func(db *sql.DB)
		expectedCode int
		expectedBody *service.ResourcesResponse
	}{
		"should return resource list when resource exists": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO resources (id, created_at, updated_at, name, amount, measurement, quantity)
					VALUES (1, ?, ?, 'Test', '1', 'Kg', 1)
				`, date, date)
			},
			expectedCode: http.StatusOK,
			expectedBody: &service.ResourcesResponse{Data: []service.Resource{{
				ID: 1, CreatedAt: DATE, UpdatedAt: DATE,
				Name:        "Test",
				Amount:      1,
				Measurement: "Kg",
				Quantity:    1,
			}}},
		},
		"should return empty list when resource not exists": {
			before:       func(bd *sql.DB) {},
			expectedCode: http.StatusOK,
			expectedBody: &service.ResourcesResponse{Data: []service.Resource{}},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			cfg, err := configuration.LoadConfig("../..")
			if err != nil {
				log.Fatal("cannot load config: ", err)
			}

			mysql := infra.MySQLConfigure(cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database, cfg.MySQL.Username,
				cfg.MySQL.Password, time.Duration(cfg.MySQL.ConnMaxLifetimeMs), cfg.MySQL.MaxOpenConns, cfg.MySQL.MaxIdleConns)
			defer mysql.DB.Close()

			resourceRepository := &repository.ResourceRepositoryImpl{DB: mysql}
			resourceService := &service.ResourceServiceImpl{ResourceRepository: resourceRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", ResourceService: resourceService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/resources", nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.ResourcesResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedBody, body)

			// after
			mysql.DB.Exec(`DELETE FROM resources`)
			mysql.DB.Exec(`ALTER TABLE resources AUTO_INCREMENT=1`)
		})
	}
}

func Test_ResourceApi_FindOneByID(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"

	cases := map[string]struct {
		before          func(db *sql.DB)
		inputResourceID string
		expectedCode    int
		expectedBody    *service.ResourceResponse
		expectedErr     *api.HttpError
	}{
		"shouldl return resource when resource exists": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO resources (id, created_at, updated_at, name, amount, measurement, quantity)
					VALUES (1, ?, ?, 'Test', '1', 'Kg', 1)
				`, date, date)
			},
			inputResourceID: "1",
			expectedCode:    http.StatusOK,
			expectedBody: &service.ResourceResponse{Data: &service.Resource{
				ID: 1, CreatedAt: DATE, UpdatedAt: DATE,
				Name:        "Test",
				Amount:      1,
				Measurement: "Kg",
				Quantity:    1,
			}},
			expectedErr: &api.HttpError{},
		},
		"should throw bad request error when reosurceID is not a number": {
			before:          func(db *sql.DB) {},
			inputResourceID: "a",
			expectedCode:    http.StatusBadRequest,
			expectedBody:    &service.ResourceResponse{},
			expectedErr:     &api.HttpError{Code: http.StatusBadRequest, Message: "invalid resourceID"},
		},
		"should throw not found error when resource not exists": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			expectedCode:    http.StatusNotFound,
			expectedBody:    &service.ResourceResponse{},
			expectedErr:     &api.HttpError{Code: http.StatusNotFound, Message: "resource 1 not found"},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			cfg, err := configuration.LoadConfig("../..")
			if err != nil {
				log.Fatal("cannot load config: ", err)
			}

			mysql := infra.MySQLConfigure(cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database, cfg.MySQL.Username,
				cfg.MySQL.Password, time.Duration(cfg.MySQL.ConnMaxLifetimeMs), cfg.MySQL.MaxOpenConns, cfg.MySQL.MaxIdleConns)
			defer mysql.DB.Close()

			resourceRepository := &repository.ResourceRepositoryImpl{DB: mysql}
			resourceService := &service.ResourceServiceImpl{ResourceRepository: resourceRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", ResourceService: resourceService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			url := "/api/v1/resources/" + cs.inputResourceID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.ResourceResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedBody, body)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM resources`)
			mysql.DB.Exec(`ALTER TABLE resources AUTO_INCREMENT=1`)
		})
	}
}

func Test_ResourceApi_Create(t *testing.T) {
	cases := map[string]struct {
		inputDto     service.CreateResourceDto
		expectedCode int
		expectedBody *service.ResourceResponse
		expectedErr  *api.HttpError
	}{
		"should return create resource": {
			inputDto:     service.CreateResourceDto{Name: "Test", Amount: 1, Measurement: "Kg", Quantity: 1},
			expectedCode: http.StatusCreated,
			expectedBody: &service.ResourceResponse{Data: &service.Resource{
				ID: 1, Name: "Test", Amount: 1, Measurement: "Kg", Quantity: 1,
			}},
			expectedErr: &api.HttpError{},
		},
		"should throw bad request error": {
			expectedCode: http.StatusBadRequest,
			expectedBody: &service.ResourceResponse{},
			expectedErr: &api.HttpError{Code: http.StatusBadRequest, Message: strings.Join([]string{
				"Key: 'CreateResourceDto.Name' Error:Field validation for 'Name' failed on the 'required' tag",
				"Key: 'CreateResourceDto.Amount' Error:Field validation for 'Amount' failed on the 'required' tag",
				"Key: 'CreateResourceDto.Measurement' Error:Field validation for 'Measurement' failed on the 'required' tag",
				"Key: 'CreateResourceDto.Quantity' Error:Field validation for 'Quantity' failed on the 'required' tag",
			}, "\n")},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			cfg, err := configuration.LoadConfig("../..")
			if err != nil {
				log.Fatal("cannot load config: ", err)
			}

			mysql := infra.MySQLConfigure(cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database, cfg.MySQL.Username,
				cfg.MySQL.Password, time.Duration(cfg.MySQL.ConnMaxLifetimeMs), cfg.MySQL.MaxOpenConns, cfg.MySQL.MaxIdleConns)
			defer mysql.DB.Close()

			resourceRepository := &repository.ResourceRepositoryImpl{DB: mysql}
			resourceService := &service.ResourceServiceImpl{ResourceRepository: resourceRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", ResourceService: resourceService}
			impl.Configure()

			// when
			b, _ := json.Marshal(cs.inputDto)
			url := "/api/v1/resources"
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var body *service.ResourceResponse
			json.Unmarshal(rec.Body.Bytes(), &body)
			if body.Data != nil {
				body.Data.CreatedAt = ""
				body.Data.UpdatedAt = ""
			}

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedBody, body)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM resources`)
			mysql.DB.Exec(`ALTER TABLE resources AUTO_INCREMENT=1`)
		})
	}
}

func Test_ResourceApi_Update(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"

	cases := map[string]struct {
		before          func(db *sql.DB)
		inputResourceID string
		inputDto        service.UpdateResourceDto
		expectedCode    int
		expectedErr     *api.HttpError
	}{
		"should update resource": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO resources (id, created_at, updated_at, name, amount, measurement, quantity)
					VALUES (1, ?, ?, 'Test', '1', 'Kg', 1)
				`, date, date)
			},
			inputResourceID: "1",
			inputDto:        service.UpdateResourceDto{Name: "Test update", Measurement: "l"},
			expectedCode:    http.StatusNoContent,
		},
		"should throw bad request error when resourceID id not number": {
			before:          func(db *sql.DB) {},
			inputResourceID: "a",
			expectedCode:    http.StatusBadRequest,
			expectedErr:     &api.HttpError{Code: http.StatusBadRequest, Message: "invalid resourceID"},
		},
		"should throw bad request error": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			expectedCode:    http.StatusBadRequest,
			expectedErr:     &api.HttpError{Code: http.StatusBadRequest, Message: "empty resource model"},
		},
		"should throw not found error when resources not exists": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			inputDto:        service.UpdateResourceDto{Name: "Test update"},
			expectedCode:    http.StatusNotFound,
			expectedErr:     &api.HttpError{Code: http.StatusNotFound, Message: "resource 1 not found"},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			cfg, err := configuration.LoadConfig("../..")
			if err != nil {
				log.Fatal("cannot load config: ", err)
			}

			mysql := infra.MySQLConfigure(cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database, cfg.MySQL.Username,
				cfg.MySQL.Password, time.Duration(cfg.MySQL.ConnMaxLifetimeMs), cfg.MySQL.MaxOpenConns, cfg.MySQL.MaxIdleConns)
			defer mysql.DB.Close()

			resourceRepository := &repository.ResourceRepositoryImpl{DB: mysql}
			resourceService := &service.ResourceServiceImpl{ResourceRepository: resourceRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", ResourceService: resourceService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputDto)
			url := "/api/v1/resources/" + cs.inputResourceID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM resources`)
			mysql.DB.Exec(`ALTER TABLE resources AUTO_INCREMENT=1`)
		})
	}
}

func Test_ResourceApi_UpdateQuantity(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"

	cases := map[string]struct {
		before          func(db *sql.DB)
		inputResourceID string
		inputDto        service.UpdateResourceQuantityDto
		expectedCode    int
		expectedErr     *api.HttpError
	}{
		"should update resource": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO resources (id, created_at, updated_at, name, amount, measurement, quantity)
					VALUES (1, ?, ?, 'Test', '1', 'Kg', 1)
				`, date, date)
			},
			inputResourceID: "1",
			inputDto:        service.UpdateResourceQuantityDto{Quantity: 0.5},
			expectedCode:    http.StatusNoContent,
		},
		"should throw bad request error when resourceID id not number": {
			before:          func(db *sql.DB) {},
			inputResourceID: "a",
			expectedCode:    http.StatusBadRequest,
			expectedErr:     &api.HttpError{Code: http.StatusBadRequest, Message: "invalid resourceID"},
		},
		"should throw bad request error when quantity is less than zero": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			inputDto:        service.UpdateResourceQuantityDto{Quantity: -1.5},
			expectedCode:    http.StatusBadRequest,
			expectedErr: &api.HttpError{
				Code:    http.StatusBadRequest,
				Message: "Key: 'UpdateResourceQuantityDto.Quantity' Error:Field validation for 'Quantity' failed on the 'gte' tag"},
		},
		"should throw bad request error": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			expectedCode:    http.StatusBadRequest,
			expectedErr: &api.HttpError{
				Code:    http.StatusBadRequest,
				Message: "Key: 'UpdateResourceQuantityDto.Quantity' Error:Field validation for 'Quantity' failed on the 'required' tag",
			},
		},
		"should throw not found error when resources not exists": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			inputDto:        service.UpdateResourceQuantityDto{Quantity: 2},
			expectedCode:    http.StatusNotFound,
			expectedErr:     &api.HttpError{Code: http.StatusNotFound, Message: "resource 1 not found"},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			cfg, err := configuration.LoadConfig("../..")
			if err != nil {
				log.Fatal("cannot load config: ", err)
			}

			mysql := infra.MySQLConfigure(cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database, cfg.MySQL.Username,
				cfg.MySQL.Password, time.Duration(cfg.MySQL.ConnMaxLifetimeMs), cfg.MySQL.MaxOpenConns, cfg.MySQL.MaxIdleConns)
			defer mysql.DB.Close()

			resourceRepository := &repository.ResourceRepositoryImpl{DB: mysql}
			resourceService := &service.ResourceServiceImpl{ResourceRepository: resourceRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", ResourceService: resourceService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputDto)
			url := fmt.Sprintf("/api/v1/resources/%s/quantity", cs.inputResourceID)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM resources`)
			mysql.DB.Exec(`ALTER TABLE resources AUTO_INCREMENT=1`)
		})
	}
}
