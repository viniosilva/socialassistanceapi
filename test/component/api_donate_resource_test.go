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

func Test_DonateResourceApi_Donate(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"

	cases := map[string]struct {
		before          func(db *sql.DB)
		inputResourceID string
		inputDto        service.DonateResourceDonateDto
		expectedCode    int
		expectedErr     *api.HttpError
	}{
		"should donate resource": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO resources (id, created_at, updated_at, name, amount, measurement, quantity)
					VALUES (1, ?, ?, 'Test', '1', 'Kg', 1)
				`, date, date)
				db.Exec(`
					INSERT INTO families (id, created_at, updated_at, name, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'Sauro', 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)
			},
			inputResourceID: "1",
			inputDto:        service.DonateResourceDonateDto{FamilyID: 1, Quantity: 1},
			expectedCode:    http.StatusNoContent,
		},
		"should throw not found error when resource is not found": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			inputDto:        service.DonateResourceDonateDto{FamilyID: 1, Quantity: 1},
			expectedCode:    http.StatusNotFound,
			expectedErr:     &api.HttpError{Code: http.StatusNotFound, Message: "resource 1 not found"},
		},
		"should throw not found error when family is not found": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO resources (id, created_at, updated_at, name, amount, measurement, quantity)
					VALUES (1, ?, ?, 'Test', '1', 'Kg', 1)
				`, date, date)
			},
			inputResourceID: "1",
			inputDto:        service.DonateResourceDonateDto{FamilyID: 1, Quantity: 1},
			expectedCode:    http.StatusNotFound,
			expectedErr:     &api.HttpError{Code: http.StatusNotFound, Message: "family 1 not found"},
		},
		"should throw bad request error when quantity is negative": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO resources (id, created_at, updated_at, name, amount, measurement, quantity)
					VALUES (1, ?, ?, 'Test', '1', 'Kg', 1)
				`, date, date)
			},
			inputResourceID: "1",
			inputDto:        service.DonateResourceDonateDto{FamilyID: 1, Quantity: 1.5},
			expectedCode:    http.StatusBadRequest,
			expectedErr:     &api.HttpError{Code: http.StatusBadRequest, Message: "resource 1 quantity is 1.0"},
		},
		"should throw bad request error when resourceID is not a number": {
			before:          func(db *sql.DB) {},
			inputResourceID: "a",
			expectedCode:    http.StatusBadRequest,
			expectedErr:     &api.HttpError{Code: http.StatusBadRequest, Message: "invalid resourceID"},
		},
		"should throw bad request error": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			expectedCode:    http.StatusBadRequest,
			expectedErr: &api.HttpError{
				Code: http.StatusBadRequest,
				Message: strings.Join([]string{
					"Key: 'DonateResourceDonateDto.FamilyID' Error:Field validation for 'FamilyID' failed on the 'required' tag",
					"Key: 'DonateResourceDonateDto.Quantity' Error:Field validation for 'Quantity' failed on the 'required' tag",
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

			donateResourceRepository := &repository.DonateResourceRepositoryImpl{DB: mysql}
			donateResourceService := &service.DonateResourceServiceImpl{DonateResourceRepository: donateResourceRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", DonateResourceService: donateResourceService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputDto)
			url := fmt.Sprintf("/api/v1/resources/%s/donate", cs.inputResourceID)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM resources_to_families`)
			mysql.DB.Exec(`DELETE FROM resources`)
			mysql.DB.Exec(`DELETE FROM families`)
			mysql.DB.Exec(`ALTER TABLE resources_to_families AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE resources AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE families AUTO_INCREMENT=1`)
		})
	}
}

func Test_DonateResourceApi_Return(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"

	cases := map[string]struct {
		before          func(db *sql.DB)
		inputResourceID string
		expectedCode    int
		expectedErr     *api.HttpError
	}{
		"should return resource": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO resources (id, created_at, updated_at, name, amount, measurement, quantity)
					VALUES (1, ?, ?, 'Test', '1', 'Kg', 1)
				`, date, date)
				db.Exec(`
					INSERT INTO families (id, created_at, updated_at, name, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'Sauro', 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)
				db.Exec(`
					INSERT INTO resources_to_families (id, created_at, resource_id, family_id, quantity)
					VALUES (1, ?, 1, 1, 1)
				`, date)
			},
			inputResourceID: "1",
			expectedCode:    http.StatusNoContent,
		},
		"should throw not found error when resource is not found": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			expectedCode:    http.StatusNotFound,
			expectedErr:     &api.HttpError{Code: http.StatusNotFound, Message: "resource 1 not found"},
		},
		"should throw bad request error when resourceID is not a number": {
			before:          func(db *sql.DB) {},
			inputResourceID: "a",
			expectedCode:    http.StatusBadRequest,
			expectedErr:     &api.HttpError{Code: http.StatusBadRequest, Message: "invalid resourceID"},
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

			donateResourceRepository := &repository.DonateResourceRepositoryImpl{DB: mysql}
			donateResourceService := &service.DonateResourceServiceImpl{DonateResourceRepository: donateResourceRepository}
			impl := &api.ApiImpl{Addr: "0.0.0.0:8080", DonateResourceService: donateResourceService}
			impl.Configure()

			cs.before(mysql.DB)

			// when
			url := fmt.Sprintf("/api/v1/resources/%s/return", cs.inputResourceID)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			assert.Equal(t, cs.expectedCode, rec.Code)
			assert.Equal(t, cs.expectedErr, httpError)

			// after
			mysql.DB.Exec(`DELETE FROM resources_to_families`)
			mysql.DB.Exec(`DELETE FROM resources`)
			mysql.DB.Exec(`DELETE FROM families`)
			mysql.DB.Exec(`ALTER TABLE resources_to_families AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE resources AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE families AUTO_INCREMENT=1`)
		})
	}
}
