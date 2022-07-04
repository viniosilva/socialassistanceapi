package component

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/viniosilva/socialassistanceapi/internal/api"
	"github.com/viniosilva/socialassistanceapi/internal/configuration"
	"github.com/viniosilva/socialassistanceapi/internal/service"
	"github.com/viniosilva/socialassistanceapi/internal/store"
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
			expectedCode: 200,
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
			expectedCode: 200,
			expectedBody: &service.AddressesResponse{Data: []service.Address{}},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			addressStore := store.NewAddressStore(mysql.DB)
			addressService := service.NewAddressService(addressStore)
			impl := api.NewApi("0.0.0.0:8080", nil, nil, addressService, nil)

			cs.before(mysql.DB)

			// when
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/addresses", nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.AddressesResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("GET /api/v1/addresses StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("GET /api/v1/addresses Body = %v, expected %v", body, cs.expectedBody)
			}

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
			expectedCode:   200,
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
		},
		"should throw bad request error when addressID is not a number": {
			before:         func(db *sql.DB) {},
			inputAddressID: "a",
			expectedCode:   400,
			expectedErr:    &api.HttpError{Code: 400, Message: "invalid addressID"},
		},
		"should throw not found error when addresses not exists": {
			before:         func(db *sql.DB) {},
			inputAddressID: "1",
			expectedCode:   404,
			expectedBody:   &service.AddressResponse{},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			addressStore := store.NewAddressStore(mysql.DB)
			addressService := service.NewAddressService(addressStore)
			impl := api.NewApi("0.0.0.0:8080", nil, nil, addressService, nil)

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

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("GET /api/v1/addresses/:addressID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("GET /api/v1/addresses/:addressID Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("GET /api/v1/addresses/:addressID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM addresses`)
			mysql.DB.Exec(`ALTER TABLE addresses AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentAddressApiCreate(t *testing.T) {
	cases := map[string]struct {
		inputAddress service.AddressDto
		expectedCode int
		expectedBody *service.AddressResponse
		expectedErr  *api.HttpError
	}{
		"should return created address": {
			inputAddress: service.AddressDto{
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			},
			expectedCode: 201,
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
		},
		"should throw bad request error": {
			expectedCode: 400,
			expectedErr: &api.HttpError{
				Code: 400,
				Message: strings.Join([]string{
					"Key: 'AddressDto.Country' Error:Field validation for 'Country' failed on the 'required' tag",
					"Key: 'AddressDto.State' Error:Field validation for 'State' failed on the 'required' tag",
					"Key: 'AddressDto.City' Error:Field validation for 'City' failed on the 'required' tag",
					"Key: 'AddressDto.Neighborhood' Error:Field validation for 'Neighborhood' failed on the 'required' tag",
					"Key: 'AddressDto.Street' Error:Field validation for 'Street' failed on the 'required' tag",
					"Key: 'AddressDto.Number' Error:Field validation for 'Number' failed on the 'required' tag",
					"Key: 'AddressDto.Complement' Error:Field validation for 'Complement' failed on the 'required' tag",
					"Key: 'AddressDto.Zipcode' Error:Field validation for 'Zipcode' failed on the 'required' tag",
				}, "\n"),
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			addressStore := store.NewAddressStore(mysql.DB)
			addressService := service.NewAddressService(addressStore)
			impl := api.NewApi("0.0.0.0:8080", nil, nil, addressService, nil)

			// when
			b, _ := json.Marshal(cs.inputAddress)
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
			if rec.Code != cs.expectedCode {
				t.Errorf("POST /api/v1/addresses StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("POST /api/v1/addresses Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("POST /api/v1/addresses BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

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
		inputAddress   service.AddressDto
		expectedCode   int
		expectedBody   *service.AddressResponse
		expectedErr    *api.HttpError
	}{
		"should return updated address": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO addresses (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, DATE, DATE)
			},
			inputAddressID: "1",
			inputAddress: service.AddressDto{
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			},
			expectedCode: 200,
			expectedBody: &service.AddressResponse{Data: &service.Address{
				ID:           1,
				CreatedAt:    DATE,
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			}},
		},
		"should return updated address when is a partial update": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO addresses (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, DATE, DATE)
			},
			inputAddressID: "1",
			inputAddress:   service.AddressDto{Number: "2"},
			expectedCode:   200,
			expectedBody: &service.AddressResponse{Data: &service.Address{
				ID:           1,
				CreatedAt:    DATE,
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "2",
				Complement:   "1",
				Zipcode:      "02180110",
			}},
		},
		"should throw bad request error when addressID is not a number": {
			before:         func(db *sql.DB) {},
			inputAddressID: "a",
			expectedCode:   400,
			expectedErr:    &api.HttpError{Code: 400, Message: "invalid addressID"},
		},
		"should throw bad request error": {
			before:         func(db *sql.DB) {},
			inputAddressID: "1",
			expectedCode:   400,
			expectedErr:    &api.HttpError{Code: 400, Message: "empty model: address"},
		},
		"should throw not found error when addresses not exists": {
			before:         func(db *sql.DB) {},
			inputAddressID: "1",
			inputAddress: service.AddressDto{
				Country:      "BR",
				State:        "SP",
				City:         "São Paulo",
				Neighborhood: "Pq. Novo Mundo",
				Street:       "R. Sd. Teodoro Francisco Ribeiro",
				Number:       "1",
				Complement:   "1",
				Zipcode:      "02180110",
			},
			expectedCode: 404,
			expectedBody: &service.AddressResponse{},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			addressStore := store.NewAddressStore(mysql.DB)
			addressService := service.NewAddressService(addressStore)
			impl := api.NewApi("0.0.0.0:8080", nil, nil, addressService, nil)

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputAddress)
			url := "/api/v1/addresses/" + cs.inputAddressID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var body *service.AddressResponse
			json.Unmarshal(rec.Body.Bytes(), &body)
			if body.Data != nil {
				body.Data.UpdatedAt = ""
			}

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("PATCH /api/v1/addresses/:addressID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("PATCH /api/v1/addresses/:addressID Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("PATCH /api/v1/addresses/:addressID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

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
			expectedCode:   204,
			expectedBody:   &service.AddressResponse{},
		},
		"should throw bad request error when addressID is not a number": {
			before:         func(db *sql.DB) {},
			inputAddressID: "a",
			expectedCode:   400,
			expectedErr:    &api.HttpError{Code: 400, Message: "invalid addressID"},
		},
		"should be successfull when addresses not exists": {
			before:         func(db *sql.DB) {},
			inputAddressID: "1",
			expectedCode:   204,
			expectedBody:   &service.AddressResponse{},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			addressStore := store.NewAddressStore(mysql.DB)
			addressService := service.NewAddressService(addressStore)
			impl := api.NewApi("0.0.0.0:8080", nil, nil, addressService, nil)

			cs.before(mysql.DB)

			// when
			url := "/api/v1/addresses/" + cs.inputAddressID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("PATCH /api/v1/addresses/:addressID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("PATCH /api/v1/addresses/:addressID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM addresses`)
			mysql.DB.Exec(`ALTER TABLE addresses AUTO_INCREMENT=1`)
		})
	}
}
