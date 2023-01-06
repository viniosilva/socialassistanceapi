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

func TestComponentPersonApiFindAll(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"

	cases := map[string]struct {
		before       func(db *sql.DB)
		expectedCode int
		expectedBody *service.PersonsResponse
	}{
		"should return person list when persons exists": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO families (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)

				db.Exec(`
					INSERT INTO persons (id, created_at, updated_at, family_id, name)
					VALUES (1, ?, ?, 1, 'Test')
				`, date, date)
			},
			expectedCode: http.StatusOK,
			expectedBody: &service.PersonsResponse{
				Data: []service.Person{{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, FamilyID: 1, Name: "Test"}},
			},
		},
		"should return empty list when persons not exists": {
			before:       func(db *sql.DB) {},
			expectedCode: http.StatusOK,
			expectedBody: &service.PersonsResponse{Data: []service.Person{}},
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

			personRepository := &repository.PersonRepositoryImpl{DB: mysql}
			personService := &service.PersonServiceImpl{PersonRepository: personRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", PersonService: personService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/persons", nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.PersonsResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedBody, body)

			// after
			mysql.DB.Exec(`DELETE FROM persons`)
			mysql.DB.Exec(`DELETE FROM families`)
			mysql.DB.Exec(`ALTER TABLE persons AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE families AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentPersonApiFindOneByID(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"

	cases := map[string]struct {
		before        func(db *sql.DB)
		inputPersonID string
		expectedCode  int
		expectedBody  *service.PersonResponse
		expectedErr   *api.HttpError
	}{
		"should return person when persons exists": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO families (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)

				db.Exec(`
					INSERT INTO persons (id, created_at, updated_at, family_id, name)
					VALUES (1, ?, ?, 1, 'Test')
				`, date, date)
			},
			inputPersonID: "1",
			expectedCode:  http.StatusOK,
			expectedBody: &service.PersonResponse{
				Data: &service.Person{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, FamilyID: 1, Name: "Test"},
			},
			expectedErr: &api.HttpError{},
		},
		"should throw bad request error when personID is not a number": {
			before:        func(db *sql.DB) {},
			inputPersonID: "a",
			expectedCode:  http.StatusBadRequest,
			expectedBody:  &service.PersonResponse{},
			expectedErr:   &api.HttpError{Code: http.StatusBadRequest, Message: "invalid personID"},
		},
		"should throw not found error when persons not exists": {
			before:        func(db *sql.DB) {},
			inputPersonID: "1",
			expectedCode:  http.StatusNotFound,
			expectedBody:  &service.PersonResponse{},
			expectedErr:   &api.HttpError{Code: http.StatusNotFound, Message: "person 1 not found"},
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

			personRepository := &repository.PersonRepositoryImpl{DB: mysql}
			personService := &service.PersonServiceImpl{PersonRepository: personRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", PersonService: personService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			url := "/api/v1/persons/" + cs.inputPersonID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.PersonResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedBody, body)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM persons`)
			mysql.DB.Exec(`DELETE FROM families`)
			mysql.DB.Exec(`ALTER TABLE persons AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE families AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentPersonApiCreate(t *testing.T) {
	cases := map[string]struct {
		before       func(db *sql.DB)
		inputDto     service.PersonCreateDto
		expectedCode int
		expectedBody *service.PersonResponse
		expectedErr  *api.HttpError
	}{
		"should return created person": {
			before: func(db *sql.DB) {
				date := "2000-01-01 12:03:00"
				db.Exec(`
					INSERT INTO families (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)
			},
			inputDto:     service.PersonCreateDto{FamilyID: 1, Name: "Test"},
			expectedCode: http.StatusCreated,
			expectedBody: &service.PersonResponse{Data: &service.Person{ID: 1, FamilyID: 1, Name: "Test"}},
			expectedErr:  &api.HttpError{},
		},
		"should throw bad request error": {
			before:       func(db *sql.DB) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: &service.PersonResponse{},
			expectedErr: &api.HttpError{
				Code: http.StatusBadRequest,
				Message: strings.Join([]string{
					"Key: 'PersonCreateDto.FamilyID' Error:Field validation for 'FamilyID' failed on the 'required' tag",
					"Key: 'PersonCreateDto.Name' Error:Field validation for 'Name' failed on the 'required' tag",
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

			personRepository := &repository.PersonRepositoryImpl{DB: mysql}
			personService := &service.PersonServiceImpl{PersonRepository: personRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", PersonService: personService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputDto)
			url := "/api/v1/persons"
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var body *service.PersonResponse
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
			mysql.DB.Exec(`DELETE FROM persons`)
			mysql.DB.Exec(`DELETE FROM families`)
			mysql.DB.Exec(`ALTER TABLE persons AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE families AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentPersonApiUpdate(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"

	cases := map[string]struct {
		before        func(db *sql.DB)
		inputPersonID string
		inputDto      service.PersonCreateDto
		expectedCode  int
		expectedErr   *api.HttpError
	}{
		"should update person": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO families (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)

				db.Exec(`
					INSERT INTO persons (id, created_at, updated_at, family_id, name)
					VALUES (1, ?, ?, 1, 'Test')
				`, date, date)
			},
			inputPersonID: "1",
			inputDto:      service.PersonCreateDto{Name: "Test update"},
			expectedCode:  http.StatusNoContent,
		},
		"should throw bad request error when personID is not a number": {
			before:        func(db *sql.DB) {},
			inputPersonID: "a",
			expectedCode:  http.StatusBadRequest,
			expectedErr:   &api.HttpError{Code: http.StatusBadRequest, Message: "invalid personID"},
		},
		"should throw bad request error": {
			before:        func(db *sql.DB) {},
			inputPersonID: "1",
			expectedCode:  http.StatusBadRequest,
			expectedErr:   &api.HttpError{Code: http.StatusBadRequest, Message: "empty person model"},
		},
		"should throw not found error when persons not exists": {
			before:        func(db *sql.DB) {},
			inputPersonID: "1",
			inputDto:      service.PersonCreateDto{Name: "Test update"},
			expectedCode:  http.StatusNotFound,
			expectedErr:   &api.HttpError{Code: http.StatusNotFound, Message: "person 1 not found"},
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

			personRepository := &repository.PersonRepositoryImpl{DB: mysql}
			personService := &service.PersonServiceImpl{PersonRepository: personRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", PersonService: personService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputDto)
			url := "/api/v1/persons/" + cs.inputPersonID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM persons`)
			mysql.DB.Exec(`DELETE FROM families`)
			mysql.DB.Exec(`ALTER TABLE persons AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE families AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentPersonApiDelete(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"

	cases := map[string]struct {
		before        func(db *sql.DB)
		inputPersonID string
		expectedCode  int
		expectedBody  *service.PersonResponse
		expectedErr   *api.HttpError
	}{
		"should be successfull": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO families (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)

				db.Exec(`
					INSERT INTO persons (id, created_at, updated_at, family_id, name)
					VALUES (1, ?, ?, 1, 'Test')
				`, date, date)
			},
			inputPersonID: "1",
			expectedCode:  http.StatusNoContent,
			expectedBody:  &service.PersonResponse{},
		},
		"should throw bad request error when personID is not a number": {
			before:        func(db *sql.DB) {},
			inputPersonID: "a",
			expectedCode:  http.StatusBadRequest,
			expectedErr:   &api.HttpError{Code: 400, Message: "invalid personID"},
		},
		"should be successfull when persons not exists": {
			before:        func(db *sql.DB) {},
			inputPersonID: "1",
			expectedCode:  http.StatusNoContent,
			expectedBody:  &service.PersonResponse{},
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

			personRepository := &repository.PersonRepositoryImpl{DB: mysql}
			personService := &service.PersonServiceImpl{PersonRepository: personRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", PersonService: personService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			url := "/api/v1/persons/" + cs.inputPersonID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM persons`)
			mysql.DB.Exec(`DELETE FROM families`)
			mysql.DB.Exec(`ALTER TABLE persons AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE families AUTO_INCREMENT=1`)
		})
	}
}
