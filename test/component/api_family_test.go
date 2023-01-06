package component

import (
	"database/sql"
	"encoding/json"
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

func TestComponentFamilyApiFindAll(t *testing.T) {
	DATE := "2000-01-01T12:03:00"

	cases := map[string]struct {
		before       func(db *sql.DB)
		expectedCode int
		expectedBody *service.FamiliesResponse
	}{
		"should return family list when families exists": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO families (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)
			},
			expectedCode: http.StatusOK,
			expectedBody: &service.FamiliesResponse{Data: []service.Family{{
				ID:           1,
				CreatedAt:    DATE,
				UpdatedAt:    DATE,
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			}}},
		},
		"should return empty list when families not exists": {
			before:       func(db *sql.DB) {},
			expectedCode: http.StatusOK,
			expectedBody: &service.FamiliesResponse{Data: []service.Family{}},
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

			familyRepository := &repository.FamilyRepositoryImpl{DB: mysql}
			familyService := &service.FamilyServiceImpl{FamilyRepository: familyRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", FamilyService: familyService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/families", nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.FamiliesResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			// clean
			for _, a := range body.Data {
				a.CreatedAt = DATE
				a.UpdatedAt = DATE
			}

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedBody, body)

			// after
			mysql.DB.Exec(`DELETE FROM families`)
			mysql.DB.Exec(`ALTER TABLE families AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentFamilyApiFindOneByID(t *testing.T) {
	DATE := "2000-01-01T12:03:00"

	cases := map[string]struct {
		before        func(db *sql.DB)
		inputFamilyID string
		expectedCode  int
		expectedBody  *service.FamilyResponse
		expectedErr   *api.HttpError
	}{
		"should return family when families exists": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO families (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, DATE, DATE)
			},
			inputFamilyID: "1",
			expectedCode:  http.StatusOK,
			expectedBody: &service.FamilyResponse{Data: &service.Family{
				ID:           1,
				CreatedAt:    DATE,
				UpdatedAt:    DATE,
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			}},
			expectedErr: &api.HttpError{},
		},
		"should throw bad request error when familyID is not a number": {
			before:        func(db *sql.DB) {},
			inputFamilyID: "a",
			expectedCode:  http.StatusBadRequest,
			expectedBody:  &service.FamilyResponse{},
			expectedErr:   &api.HttpError{Code: 400, Message: "invalid familyID"},
		},
		"should throw not found error when families not exists": {
			before:        func(db *sql.DB) {},
			inputFamilyID: "1",
			expectedCode:  http.StatusNotFound,
			expectedBody:  &service.FamilyResponse{},
			expectedErr:   &api.HttpError{Code: http.StatusNotFound, Message: "family 1 not found"},
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

			familyRepository := &repository.FamilyRepositoryImpl{DB: mysql}
			familyService := &service.FamilyServiceImpl{FamilyRepository: familyRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", FamilyService: familyService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			url := "/api/v1/families/" + cs.inputFamilyID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.FamilyResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// clean
			if body.Data != nil {
				body.Data.CreatedAt = DATE
				body.Data.UpdatedAt = DATE
			}

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedBody, body)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM families`)
			mysql.DB.Exec(`ALTER TABLE families AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentFamilyApiCreate(t *testing.T) {
	cases := map[string]struct {
		inputDto     service.FamilyCreateDto
		expectedCode int
		expectedBody *service.FamilyResponse
		expectedErr  *api.HttpError
	}{
		"should return created family": {
			inputDto: service.FamilyCreateDto{
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			},
			expectedCode: http.StatusCreated,
			expectedBody: &service.FamilyResponse{Data: &service.Family{
				ID:           1,
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			}},
			expectedErr: &api.HttpError{},
		},
		"should throw bad request error": {
			expectedCode: http.StatusBadRequest,
			expectedBody: &service.FamilyResponse{},
			expectedErr: &api.HttpError{
				Code: http.StatusBadRequest,
				Message: strings.Join([]string{
					"Key: 'FamilyCreateDto.Country' Error:Field validation for 'Country' failed on the 'required' tag",
					"Key: 'FamilyCreateDto.State' Error:Field validation for 'State' failed on the 'required' tag",
					"Key: 'FamilyCreateDto.City' Error:Field validation for 'City' failed on the 'required' tag",
					"Key: 'FamilyCreateDto.Neighborhood' Error:Field validation for 'Neighborhood' failed on the 'required' tag",
					"Key: 'FamilyCreateDto.Street' Error:Field validation for 'Street' failed on the 'required' tag",
					"Key: 'FamilyCreateDto.Number' Error:Field validation for 'Number' failed on the 'required' tag",
					"Key: 'FamilyCreateDto.Complement' Error:Field validation for 'Complement' failed on the 'required' tag",
					"Key: 'FamilyCreateDto.Zipcode' Error:Field validation for 'Zipcode' failed on the 'required' tag",
				}, "\n"),
			},
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

			familyRepository := &repository.FamilyRepositoryImpl{DB: mysql}
			familyService := &service.FamilyServiceImpl{FamilyRepository: familyRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", FamilyService: familyService}
			impl.Configure()

			// when
			b, _ := json.Marshal(cs.inputDto)
			url := "/api/v1/families"
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var body *service.FamilyResponse
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
			mysql.DB.Exec(`DELETE FROM families`)
			mysql.DB.Exec(`ALTER TABLE families AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentFamilyApiUpdate(t *testing.T) {
	DATE := "2000-01-01T12:03:00"

	cases := map[string]struct {
		before        func(db *sql.DB)
		inputFamilyID string
		inputDto      service.FamilyCreateDto
		expectedCode  int
		expectedErr   *api.HttpError
	}{
		"should update family": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO families (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, DATE, DATE)
			},
			inputFamilyID: "1",
			inputDto: service.FamilyCreateDto{
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			},
			expectedCode: http.StatusNoContent,
		},
		"should update family when is a partial update": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO families (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, DATE, DATE)
			},
			inputFamilyID: "1",
			inputDto:      service.FamilyCreateDto{Number: "2"},
			expectedCode:  http.StatusNoContent,
		},
		"should throw bad request error when familyID is not a number": {
			before:        func(db *sql.DB) {},
			inputFamilyID: "a",
			expectedCode:  http.StatusBadRequest,
			expectedErr:   &api.HttpError{Code: http.StatusBadRequest, Message: "invalid familyID"},
		},
		"should throw bad request error": {
			before:        func(db *sql.DB) {},
			inputFamilyID: "1",
			expectedCode:  http.StatusBadRequest,
			expectedErr:   &api.HttpError{Code: http.StatusBadRequest, Message: "empty family model"},
		},
		"should throw not found error when families not exists": {
			before:        func(db *sql.DB) {},
			inputFamilyID: "1",
			inputDto: service.FamilyCreateDto{
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			},
			expectedCode: http.StatusNotFound,
			expectedErr:  &api.HttpError{Code: http.StatusNotFound, Message: "family 1 not found"},
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

			familyRepository := &repository.FamilyRepositoryImpl{DB: mysql}
			familyService := &service.FamilyServiceImpl{FamilyRepository: familyRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", FamilyService: familyService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputDto)
			url := "/api/v1/families/" + cs.inputFamilyID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM families`)
			mysql.DB.Exec(`ALTER TABLE families AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentFamilyApiDelete(t *testing.T) {
	DATE := "2000-01-01T12:03:00"

	cases := map[string]struct {
		before        func(db *sql.DB)
		inputFamilyID string
		expectedCode  int
		expectedBody  *service.FamilyResponse
		expectedErr   *api.HttpError
	}{
		"should be successfull": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO families (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, DATE, DATE)
			},
			inputFamilyID: "1",
			expectedCode:  http.StatusNoContent,
			expectedBody:  &service.FamilyResponse{},
		},
		"should throw bad request error when familyID is not a number": {
			before:        func(db *sql.DB) {},
			inputFamilyID: "a",
			expectedCode:  http.StatusBadRequest,
			expectedErr:   &api.HttpError{Code: http.StatusBadRequest, Message: "invalid familyID"},
		},
		"should be successfull when families not exists": {
			before:        func(db *sql.DB) {},
			inputFamilyID: "1",
			expectedCode:  http.StatusNoContent,
			expectedBody:  &service.FamilyResponse{},
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

			familyRepository := &repository.FamilyRepositoryImpl{DB: mysql}
			familyService := &service.FamilyServiceImpl{FamilyRepository: familyRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", FamilyService: familyService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			url := "/api/v1/families/" + cs.inputFamilyID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM families`)
			mysql.DB.Exec(`ALTER TABLE families AUTO_INCREMENT=1`)
		})
	}
}
