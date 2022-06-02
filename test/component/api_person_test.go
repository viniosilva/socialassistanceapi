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
		expectedBody *service.PeopleResponse
	}{
		"should return person list when people exists": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO people (id, created_at, updated_at, name)
					VALUES (1, ?, ?, 'Test')
				`, date, date)
			},
			expectedCode: 200,
			expectedBody: &service.PeopleResponse{Data: []service.Person{{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}}},
		},
		"should return empty list when people not exists": {
			before:       func(db *sql.DB) {},
			expectedCode: 200,
			expectedBody: &service.PeopleResponse{Data: []service.Person{}},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			personStore := store.NewPersonStore(mysql.DB)
			personService := service.NewPersonService(personStore)
			impl := api.NewApi("0.0.0.0:8080", nil, personService)

			cs.before(mysql.DB)

			// when
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/people", nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.PeopleResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("GET /api/v1/people StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("GET /api/v1/people Body = %v, expected %v", body, cs.expectedBody)
			}

			// after
			mysql.DB.Exec(`DELETE FROM people`)
			mysql.DB.Exec(`ALTER TABLE people AUTO_INCREMENT=1`)
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
		"should return person when people exists": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO people (id, created_at, updated_at, name)
					VALUES (1, ?, ?, 'Test')
				`, DATE, DATE)
			},
			inputPersonID: "1",
			expectedCode:  200,
			expectedBody:  &service.PersonResponse{Data: &service.Person{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}},
		},
		"should throw bad request error when personID is not a number": {
			before:        func(db *sql.DB) {},
			inputPersonID: "a",
			expectedCode:  400,
			expectedErr:   &api.HttpError{Code: 400, Message: "invalid personID"},
		},
		"should throw not found error when people not exists": {
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

			personStore := store.NewPersonStore(mysql.DB)
			personService := service.NewPersonService(personStore)
			impl := api.NewApi("0.0.0.0:8080", nil, personService)

			cs.before(mysql.DB)

			// when
			url := "/api/v1/people/" + cs.inputPersonID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.PersonResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("GET /api/v1/people/:personID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("GET /api/v1/people/:personID Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("GET /api/v1/people/:personID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM people`)
			mysql.DB.Exec(`ALTER TABLE people AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentPersonApiCreate(t *testing.T) {
	cases := map[string]struct {
		inputPerson  service.PersonDto
		expectedCode int
		expectedBody *service.PersonResponse
		expectedErr  *api.HttpError
	}{
		"should return created person": {
			inputPerson:  service.PersonDto{Name: "Test"},
			expectedCode: 201,
			expectedBody: &service.PersonResponse{Data: &service.Person{ID: 1, Name: "Test"}},
		},
		"should throw bad request error": {
			expectedCode: 400,
			expectedErr:  &api.HttpError{Code: 400, Message: "Key: 'PersonDto.Name' Error:Field validation for 'Name' failed on the 'required' tag"},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			personStore := store.NewPersonStore(mysql.DB)
			personService := service.NewPersonService(personStore)
			impl := api.NewApi("0.0.0.0:8080", nil, personService)

			// when
			b, _ := json.Marshal(cs.inputPerson)
			url := "/api/v1/people"
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
				t.Errorf("POST /api/v1/people StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("POST /api/v1/people Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("POST /api/v1/people BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM people`)
			mysql.DB.Exec(`ALTER TABLE people AUTO_INCREMENT=1`)
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
				db.Exec(`
					INSERT INTO people (id, created_at, updated_at, name)
					VALUES (1, ?, ?, 'Test')
				`, DATE, DATE)
			},
			inputPersonID: "1",
			inputPerson:   service.PersonDto{Name: "Test update"},
			expectedCode:  200,
			expectedBody:  &service.PersonResponse{Data: &service.Person{ID: 1, CreatedAt: DATE, Name: "Test update"}},
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
			expectedErr:   &api.HttpError{Code: 400, Message: "Key: 'PersonDto.Name' Error:Field validation for 'Name' failed on the 'required' tag"},
		},
		"should throw not found error when people not exists": {
			before:        func(db *sql.DB) {},
			inputPersonID: "1",
			inputPerson:   service.PersonDto{Name: "Test update"},
			expectedCode:  404,
			expectedBody:  &service.PersonResponse{},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			personStore := store.NewPersonStore(mysql.DB)
			personService := service.NewPersonService(personStore)
			impl := api.NewApi("0.0.0.0:8080", nil, personService)

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputPerson)
			url := "/api/v1/people/" + cs.inputPersonID
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
				t.Errorf("PATCH /api/v1/people/:personID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("PATCH /api/v1/people/:personID Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("PATCH /api/v1/people/:personID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM people`)
			mysql.DB.Exec(`ALTER TABLE people AUTO_INCREMENT=1`)
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
				db.Exec(`
					INSERT INTO people (id, created_at, updated_at, name)
					VALUES (1, ?, ?, 'Test')
				`, DATE, DATE)
			},
			inputPersonID: "1",
			expectedCode:  200,
			expectedBody:  &service.PersonResponse{},
		},
		"should throw bad request error when personID is not a number": {
			before:        func(db *sql.DB) {},
			inputPersonID: "a",
			expectedCode:  400,
			expectedErr:   &api.HttpError{Code: 400, Message: "invalid personID"},
		},
		"should be successfull when people not exists": {
			before:        func(db *sql.DB) {},
			inputPersonID: "1",
			expectedCode:  200,
			expectedBody:  &service.PersonResponse{},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			personStore := store.NewPersonStore(mysql.DB)
			personService := service.NewPersonService(personStore)
			impl := api.NewApi("0.0.0.0:8080", nil, personService)

			cs.before(mysql.DB)

			// when
			url := "/api/v1/people/" + cs.inputPersonID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("PATCH /api/v1/people/:personID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("PATCH /api/v1/people/:personID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM people`)
			mysql.DB.Exec(`ALTER TABLE people AUTO_INCREMENT=1`)
		})
	}
}
