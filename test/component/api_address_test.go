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
	"github.com/viniosilva/socialassistanceapi/internal/repository"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

func TestComponentAddressApiFindAll(t *testing.T) {
	DATE := "2000-01-01T12:03:00"

	cases := map[string]struct {
		before       func(db *sql.DB)
		expectedCode int
		expectedBody *service.AddressesResponse
	}{
		"should return address list when addresses exists": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO addresses (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)
			},
			expectedCode: http.StatusOK,
			expectedBody: &service.AddressesResponse{Data: []service.Address{{
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
		"should return empty list when addresses not exists": {
			before:       func(db *sql.DB) {},
			expectedCode: http.StatusOK,
			expectedBody: &service.AddressesResponse{Data: []service.Address{}},
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

			addressRepository := &repository.AddressRepositoryImpl{DB: mysql}
			addressService := &service.AddressServiceImpl{AddressRepository: addressRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", AddressService: addressService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/addresses", nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.AddressesResponse
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
			mysql.DB.Exec(`DELETE FROM addresses`)
			mysql.DB.Exec(`ALTER TABLE addresses AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentAddressApiFindOneByID(t *testing.T) {
	DATE := "2000-01-01T12:03:00"

	cases := map[string]struct {
		before         func(db *sql.DB)
		inputAddressID string
		expectedCode   int
		expectedBody   *service.AddressResponse
		expectedErr    *api.HttpError
	}{
		"should return address when addresses exists": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO addresses (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, DATE, DATE)
			},
			inputAddressID: "1",
			expectedCode:   http.StatusOK,
			expectedBody: &service.AddressResponse{Data: &service.Address{
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
		"should throw bad request error when addressID is not a number": {
			before:         func(db *sql.DB) {},
			inputAddressID: "a",
			expectedCode:   http.StatusBadRequest,
			expectedBody:   &service.AddressResponse{},
			expectedErr:    &api.HttpError{Code: 400, Message: "invalid addressID"},
		},
		"should throw not found error when addresses not exists": {
			before:         func(db *sql.DB) {},
			inputAddressID: "1",
			expectedCode:   http.StatusNotFound,
			expectedBody:   &service.AddressResponse{},
			expectedErr:    &api.HttpError{Code: http.StatusNotFound, Message: "address 1 not found"},
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

			addressRepository := &repository.AddressRepositoryImpl{DB: mysql}
			addressService := &service.AddressServiceImpl{AddressRepository: addressRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", AddressService: addressService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			url := "/api/v1/addresses/" + cs.inputAddressID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.AddressResponse
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
			mysql.DB.Exec(`DELETE FROM addresses`)
			mysql.DB.Exec(`ALTER TABLE addresses AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentAddressApiCreate(t *testing.T) {
	cases := map[string]struct {
		inputDto     service.AddressCreateDto
		expectedCode int
		expectedBody *service.AddressResponse
		expectedErr  *api.HttpError
	}{
		"should return created address": {
			inputDto: service.AddressCreateDto{
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
			expectedBody: &service.AddressResponse{Data: &service.Address{
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
			expectedBody: &service.AddressResponse{},
			expectedErr: &api.HttpError{
				Code: http.StatusBadRequest,
				Message: strings.Join([]string{
					"Key: 'AddressCreateDto.Country' Error:Field validation for 'Country' failed on the 'required' tag",
					"Key: 'AddressCreateDto.State' Error:Field validation for 'State' failed on the 'required' tag",
					"Key: 'AddressCreateDto.City' Error:Field validation for 'City' failed on the 'required' tag",
					"Key: 'AddressCreateDto.Neighborhood' Error:Field validation for 'Neighborhood' failed on the 'required' tag",
					"Key: 'AddressCreateDto.Street' Error:Field validation for 'Street' failed on the 'required' tag",
					"Key: 'AddressCreateDto.Number' Error:Field validation for 'Number' failed on the 'required' tag",
					"Key: 'AddressCreateDto.Complement' Error:Field validation for 'Complement' failed on the 'required' tag",
					"Key: 'AddressCreateDto.Zipcode' Error:Field validation for 'Zipcode' failed on the 'required' tag",
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

			mysql := configuration.MySQLConfigure(cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database, cfg.MySQL.Username,
				cfg.MySQL.Password, time.Duration(cfg.MySQL.ConnMaxLifetimeMs), cfg.MySQL.MaxOpenConns, cfg.MySQL.MaxIdleConns)
			defer mysql.DB.Close()

			addressRepository := &repository.AddressRepositoryImpl{DB: mysql}
			addressService := &service.AddressServiceImpl{AddressRepository: addressRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", AddressService: addressService}
			impl.Configure()

			// when
			b, _ := json.Marshal(cs.inputDto)
			url := "/api/v1/addresses"
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var body *service.AddressResponse
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
			mysql.DB.Exec(`DELETE FROM addresses`)
			mysql.DB.Exec(`ALTER TABLE addresses AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentAddressApiUpdate(t *testing.T) {
	DATE := "2000-01-01T12:03:00"

	cases := map[string]struct {
		before         func(db *sql.DB)
		inputAddressID string
		inputDto       service.AddressCreateDto
		expectedCode   int
		expectedErr    *api.HttpError
	}{
		"should update address": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO addresses (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, DATE, DATE)
			},
			inputAddressID: "1",
			inputDto: service.AddressCreateDto{
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
		"should update address when is a partial update": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO addresses (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, DATE, DATE)
			},
			inputAddressID: "1",
			inputDto:       service.AddressCreateDto{Number: "2"},
			expectedCode:   http.StatusNoContent,
		},
		"should throw bad request error when addressID is not a number": {
			before:         func(db *sql.DB) {},
			inputAddressID: "a",
			expectedCode:   http.StatusBadRequest,
			expectedErr:    &api.HttpError{Code: http.StatusBadRequest, Message: "invalid addressID"},
		},
		"should throw bad request error": {
			before:         func(db *sql.DB) {},
			inputAddressID: "1",
			expectedCode:   http.StatusBadRequest,
			expectedErr:    &api.HttpError{Code: http.StatusBadRequest, Message: "empty address model"},
		},
		"should throw not found error when addresses not exists": {
			before:         func(db *sql.DB) {},
			inputAddressID: "1",
			inputDto: service.AddressCreateDto{
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
			expectedErr:  &api.HttpError{Code: http.StatusNotFound, Message: "address 1 not found"},
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

			addressRepository := &repository.AddressRepositoryImpl{DB: mysql}
			addressService := &service.AddressServiceImpl{AddressRepository: addressRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", AddressService: addressService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputDto)
			url := "/api/v1/addresses/" + cs.inputAddressID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM addresses`)
			mysql.DB.Exec(`ALTER TABLE addresses AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentAddressApiDelete(t *testing.T) {
	DATE := "2000-01-01T12:03:00"

	cases := map[string]struct {
		before         func(db *sql.DB)
		inputAddressID string
		expectedCode   int
		expectedBody   *service.AddressResponse
		expectedErr    *api.HttpError
	}{
		"should be successfull": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO addresses (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, DATE, DATE)
			},
			inputAddressID: "1",
			expectedCode:   http.StatusNoContent,
			expectedBody:   &service.AddressResponse{},
		},
		"should throw bad request error when addressID is not a number": {
			before:         func(db *sql.DB) {},
			inputAddressID: "a",
			expectedCode:   http.StatusBadRequest,
			expectedErr:    &api.HttpError{Code: http.StatusBadRequest, Message: "invalid addressID"},
		},
		"should be successfull when addresses not exists": {
			before:         func(db *sql.DB) {},
			inputAddressID: "1",
			expectedCode:   http.StatusNoContent,
			expectedBody:   &service.AddressResponse{},
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

			addressRepository := &repository.AddressRepositoryImpl{DB: mysql}
			addressService := &service.AddressServiceImpl{AddressRepository: addressRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", AddressService: addressService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			url := "/api/v1/addresses/" + cs.inputAddressID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM addresses`)
			mysql.DB.Exec(`ALTER TABLE addresses AUTO_INCREMENT=1`)
		})
	}
}
