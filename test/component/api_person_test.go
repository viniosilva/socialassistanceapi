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
					INSERT INTO addresses (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)

				db.Exec(`
					INSERT INTO persons (id, created_at, updated_at, address_id, name)
					VALUES (1, ?, ?, 1, 'Test')
				`, date, date)
			},
			expectedCode: 200,
			expectedBody: &service.PersonsResponse{
				Data: []service.Person{{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, AddressID: 1, Name: "Test"}},
			},
		},
		"should return empty list when persons not exists": {
			before:       func(db *sql.DB) {},
			expectedCode: 200,
			expectedBody: &service.PersonsResponse{Data: []service.Person{}},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			personStore := store.NewPersonStore(mysql)
			personService := service.NewPersonService(personStore)
			impl := api.NewApi("0.0.0.0:8080", nil, personService, nil, nil)

			cs.before(mysql.DB)

			// when
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/persons", nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.PersonsResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("GET /api/v1/persons StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("GET /api/v1/persons Body = %v, expected %v", body, cs.expectedBody)
			}

			// after
			mysql.DB.Exec(`DELETE FROM persons`)
			mysql.DB.Exec(`DELETE FROM addresses`)
			mysql.DB.Exec(`ALTER TABLE persons AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE addresses AUTO_INCREMENT=1`)
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
					INSERT INTO addresses (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)

				db.Exec(`
					INSERT INTO persons (id, created_at, updated_at, address_id, name)
					VALUES (1, ?, ?, 1, 'Test')
				`, date, date)
			},
			inputPersonID: "1",
			expectedCode:  200,
			expectedBody: &service.PersonResponse{
				Data: &service.Person{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, AddressID: 1, Name: "Test"},
			},
		},
		"should throw bad request error when personID is not a number": {
			before:        func(db *sql.DB) {},
			inputPersonID: "a",
			expectedCode:  400,
			expectedErr:   &api.HttpError{Code: 400, Message: "invalid personID"},
		},
		"should throw not found error when persons not exists": {
			before:        func(db *sql.DB) {},
			inputPersonID: "1",
			expectedCode:  404,
			expectedBody:  &service.PersonResponse{},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			personStore := store.NewPersonStore(mysql)
			personService := service.NewPersonService(personStore)
			impl := api.NewApi("0.0.0.0:8080", nil, personService, nil, nil)

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
			if rec.Code != cs.expectedCode {
				t.Errorf("GET /api/v1/persons/:personID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("GET /api/v1/persons/:personID Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("GET /api/v1/persons/:personID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM persons`)
			mysql.DB.Exec(`DELETE FROM addresses`)
			mysql.DB.Exec(`ALTER TABLE persons AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE addresses AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentPersonApiCreate(t *testing.T) {
	cases := map[string]struct {
		before       func(db *sql.DB)
		inputPerson  service.PersonDto
		expectedCode int
		expectedBody *service.PersonResponse
		expectedErr  *api.HttpError
	}{
		"should return created person": {
			before: func(db *sql.DB) {
				date := "2000-01-01 12:03:00"
				db.Exec(`
					INSERT INTO addresses (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)
			},
			inputPerson:  service.PersonDto{AddressID: 1, Name: "Test"},
			expectedCode: 201,
			expectedBody: &service.PersonResponse{Data: &service.Person{ID: 1, AddressID: 1, Name: "Test"}},
		},
		"should throw bad request error": {
			before:       func(db *sql.DB) {},
			expectedCode: 400,
			expectedErr: &api.HttpError{
				Code: 400,
				Message: strings.Join([]string{
					"Key: 'PersonDto.AddressID' Error:Field validation for 'AddressID' failed on the 'required' tag",
					"Key: 'PersonDto.Name' Error:Field validation for 'Name' failed on the 'required' tag",
				}, "\n"),
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			personStore := store.NewPersonStore(mysql)
			personService := service.NewPersonService(personStore)
			impl := api.NewApi("0.0.0.0:8080", nil, personService, nil, nil)

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputPerson)
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
			if rec.Code != cs.expectedCode {
				t.Errorf("POST /api/v1/persons StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("POST /api/v1/persons Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("POST /api/v1/persons BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM persons`)
			mysql.DB.Exec(`DELETE FROM addresses`)
			mysql.DB.Exec(`ALTER TABLE persons AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE addresses AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentPersonApiUpdate(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"

	cases := map[string]struct {
		before        func(db *sql.DB)
		inputPersonID string
		inputPerson   service.PersonDto
		expectedCode  int
		expectedBody  *service.PersonResponse
		expectedErr   *api.HttpError
	}{
		"should return updated person": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO addresses (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)

				db.Exec(`
					INSERT INTO persons (id, created_at, updated_at, address_id, name)
					VALUES (1, ?, ?, 1, 'Test')
				`, date, date)
			},
			inputPersonID: "1",
			inputPerson:   service.PersonDto{Name: "Test update"},
			expectedCode:  200,
			expectedBody: &service.PersonResponse{
				Data: &service.Person{ID: 1, CreatedAt: DATE, AddressID: 1, Name: "Test update"},
			},
		},
		"should throw bad request error when personID is not a number": {
			before:        func(db *sql.DB) {},
			inputPersonID: "a",
			expectedCode:  400,
			expectedErr:   &api.HttpError{Code: 400, Message: "invalid personID"},
		},
		"should throw bad request error": {
			before:        func(db *sql.DB) {},
			inputPersonID: "1",
			expectedCode:  400,
			expectedErr:   &api.HttpError{Code: 400, Message: "empty model: person"},
		},
		"should throw not found error when persons not exists": {
			before:        func(db *sql.DB) {},
			inputPersonID: "1",
			inputPerson:   service.PersonDto{Name: "Test update"},
			expectedCode:  404,
			expectedErr:   &api.HttpError{Code: 404, Message: "not found: person"},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			personStore := store.NewPersonStore(mysql)
			personService := service.NewPersonService(personStore)
			impl := api.NewApi("0.0.0.0:8080", nil, personService, nil, nil)

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputPerson)
			url := "/api/v1/persons/" + cs.inputPersonID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var body *service.PersonResponse
			json.Unmarshal(rec.Body.Bytes(), &body)
			if body.Data != nil {
				body.Data.UpdatedAt = ""
			}

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("PATCH /api/v1/persons/:personID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("PATCH /api/v1/persons/:personID Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("PATCH /api/v1/persons/:personID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM persons`)
			mysql.DB.Exec(`DELETE FROM addresses`)
			mysql.DB.Exec(`ALTER TABLE persons AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE addresses AUTO_INCREMENT=1`)
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
					INSERT INTO addresses (id, created_at, updated_at, country,
						state, city, neighborhood, street, number, complement, zipcode)
					VALUES (1, ?, ?, 'BR', 'SP', 'São Paulo', 'Pq. Novo Mundo', 'R. Sd. Teodoro Francisco Ribeiro', '1', '1', '02180110')
				`, date, date)

				db.Exec(`
					INSERT INTO persons (id, created_at, updated_at, address_id, name)
					VALUES (1, ?, ?, 1, 'Test')
				`, date, date)
			},
			inputPersonID: "1",
			expectedCode:  204,
			expectedBody:  &service.PersonResponse{},
		},
		"should throw bad request error when personID is not a number": {
			before:        func(db *sql.DB) {},
			inputPersonID: "a",
			expectedCode:  400,
			expectedErr:   &api.HttpError{Code: 400, Message: "invalid personID"},
		},
		"should be successfull when persons not exists": {
			before:        func(db *sql.DB) {},
			inputPersonID: "1",
			expectedCode:  204,
			expectedBody:  &service.PersonResponse{},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			personStore := store.NewPersonStore(mysql)
			personService := service.NewPersonService(personStore)
			impl := api.NewApi("0.0.0.0:8080", nil, personService, nil, nil)

			cs.before(mysql.DB)

			// when
			url := "/api/v1/persons/" + cs.inputPersonID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("PATCH /api/v1/persons/:personID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("PATCH /api/v1/persons/:personID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM persons`)
			mysql.DB.Exec(`DELETE FROM addresses`)
			mysql.DB.Exec(`ALTER TABLE persons AUTO_INCREMENT=1`)
			mysql.DB.Exec(`ALTER TABLE addresses AUTO_INCREMENT=1`)
		})
	}
}
