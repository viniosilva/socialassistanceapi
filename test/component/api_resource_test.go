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

const DATE = "2000-01-01T12:03:00"

func TestComponentResourceApiFindAll(t *testing.T) {
	cases := map[string]struct {
		before       func(db *sql.DB)
		expectedCode int
		expectedBody *service.ResourceResponse
	}{
		"should return resource list when resource exists": {
			before: func(db *sql.DB) {
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO resource (id, created_at, update_at, name, amount, measurement)
					VALUES (1, ?, ?, ?, ?, ?, 'Test')
				`, date, date)
			},
			expectedCode: 200,
			expectedBody: &service.ResourceResponse{Data: []service.Resource{{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}}},
		},
		"should return empty list when resource not exists": {
			before:       func(bd *sql.DB) {},
			expectedCode: 200,
			expectedBody: &service.ResourceResponse{Data: []service.Resource{}},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			resourceStore := store.NewResourceStore(mysql.DB)
			resourceService := service.NewResourceService(resourceStore)
			impl := api.NewApi("0.0.0.0:8080", nil, resourceService)

			cs.before(mysql.DB)

			// when
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/resource", nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.ResourceService
			json.Unmarshal(rec.Body.Bytes(), &body)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("GET /api/v1/resource StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("GET /api/v1/resource StatusCode = %v, expected %v", rec.Code, cs.expectedBody)
			}

			// after
			mysql.DB.Exec(`DELETE FROM resource`)
			mysql.DB.Exec(`ALTER TABLE resource AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentResourceApiFindOneByID(t *testing.T) {
	cases := map[string]struct {
		before          func(db *sql.DB)
		inputResourceID string
		expectedCode    int
		expectedBody    *service.ResourceResponse
		expectedErr     *api.HttpError
	}{
		"shouldl return resource when resource exists": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO resource (id, created_at, update_at, name, amount, measurement)
					VALUES (1, ?, ?, ?, ?, ?, 'Test')
				`, DATE, DATE)
			},
			inputResourceID: "1",
			expectedCode:    200,
			expectedBody:    &service.ResourceResponse{Data: &service.Person{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}},
		},
		"should throw bad request error when reosurceID is not a number": {
			before:          func(db *sql.DB) {},
			inputResourceID: "a",
			expectedCode:    400,
			expectedErr:     &api.HttpError{Code: 400, Message: "invalid resourceID"},
		},
		"should throw not found error when resource not exists": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			expectedCode:    404,
			expectedBody:    &service.ResourceResponse{},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			resourceStore := store.NewResourceStore(mysql.DB)
			resourceService := service.NewResourceService(resourceStore)
			impl := api.NewApi("0.0.0.0:8080", nil, resourceService)

			cs.before(mysql.DB)

			// when
			url := "/api/v1/resource/" + cs.inputResourceID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.ResourceResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), httpError)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("GET /api/v1/resource/:resourceID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("GET /api/v1/resource/:resourceID Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("GET /api/v1/resource/:resourceID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM resource`)
			mysql.DB.Exec(`ALTER TABLE resource AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentResourceApiCreate(t *testing.T) {
	cases := map[string]struct {
		inputResource service.ResourceDto
		expectedCode  int
		expectedBody  *service.ResourceResponse
		expectedErr   *api.HttpError
	}{
		"should return create resource": {
			inputResource: service.ResourceDto{Name: "Test"},
			expectedCode:  200,
			expectedBody:  &service.ResourceResponse{Data: &service.Resource{ID: 1, Name: "Test"}},
		},
		"should throw bad request error": {
			expectedCode: 400,
			expectedErr:  &api.HttpError{Code: 400, Message: "Key: 'ResourceDto.Name' Error:Filed validation for 'Name' falied on the 'required' tag"},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			resourceStore := store.NewResourceStore(mysql.DB)
			resourceService := service.NewPersonService(resourceStore)
			impl := api.NewApi("0.0.0.0:8080", nil, resourceService)

			// when
			b, _ := json.Marshal(cs.inputResource)
			url := "/api/v1/resource"
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
				t.Errorf("POST /api/v1/resource StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("POST /api/v1/resource Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("POST /api/v1/resource BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM resource`)
			mysql.DB.Exec(`ALTER TABLE resource AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentResourceApiUpdate(t *testing.T) {
	cases := map[string]struct {
		before          func(db *sql.DB)
		inputResourceID string
		inputResource   service.ResourceDto
		expectedCode    int
		expectedBody    *service.ResourceResponse
		expectedErr     *api.HttpError
	}{
		"should return updated resource": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO resource (id, created_at, updated_at, name, amount, measurement)
					VALUES (1, ?, ?, ?, ?, ? 'Test')
				`, DATE, DATE)
			},
			inputResourceID: "1",
			inputResource:   service.ResourceDto{Name: "Test Update"},
			expectedCode:    200,
			expectedBody:    &service.ResourceResponse{Data: &service.Resource{ID: 1, CreatedAt: DATE, Name: "Test update"}},
		},
		"should throw bad request error when resourceID id not number": {
			before:          func(db *sql.DB) {},
			inputResourceID: "a",
			expectedCode:    400,
			expectedErr:     &api.HttpError{Code: 400, Message: "invalid resourceID"},
		},
		"should throw dab resquest error": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			expectedCode:    400,
			expectedErr:     &api.HttpError{Code: 400, Message: "Key: 'ResourceDto.Name' Error:Field validation for 'Name' failed on the 'required' tag"},
		},
		"should throw not found error when resources not exists": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			inputResource:   service.ResourceDto{Name: "Test update"},
			expectedCode:    404,
			expectedBody:    &service.ResourceResponse{},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			resourceStore := store.NewPersonStore(mysql.DB)
			resourceService := service.NewPersonService(resourceStore)
			impl := api.NewApi("0.0.0.0:8080", nil, resourceService)

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputResource)
			url := "/api/v1/people/" + cs.inputResourceID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var body *service.ResourceResponse
			json.Unmarshal(rec.Body.Bytes(), &body)
			if body.Data != nil {
				body.Data.UpdatedAt = ""
			}

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("PATCH /api/v1/resource/:resourceID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("PATCH /api/v1/resource/:resourceID Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("PATCH /api/v1/resource/:resourceID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM resource`)
			mysql.DB.Exec(`ALTER TABLE resource AUTO_INCREMENT=1`)
		})
	}
}
func TestComponentResourceApiDelete(t *testing.T) {
	cases := map[string]struct {
		before          func(db *sql.DB)
		inputResourceID string
		expectedCode    int
		expectedBody    *service.ResourceResponse
		expectedErr     *api.HttpError
	}{
		"should be successfull": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO resource (id, created_at, updated_at, name, amount, measurement)
					VALUES (1, ?, ?, ?, ?, ? 'Test')
				`, DATE, DATE)
			},
			inputResourceID: "1",
			expectedCode:    200,
			expectedBody:    &service.ResourceResponse{},
		},
		"should throw bad request error when resourceID is not a number": {
			before:          func(db *sql.DB) {},
			inputResourceID: "a",
			expectedCode:    400,
			expectedErr:     &api.HttpError{Code: 400, Message: "invalid resourceID"},
		},
		"should be successfull when people not exists": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			expectedCode:    200,
			expectedBody:    &service.ResourceResponse{},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			resourceStore := store.NewPersonStore(mysql.DB)
			resourceService := service.NewPersonService(resourceStore)
			impl := api.NewApi("0.0.0.0:8080", nil, resourceService)

			cs.before(mysql.DB)

			// when
			url := "/api/v1/people/" + cs.inputResourceID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("PATCH /api/v1/resource/:resourceID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("PATCH /api/v1/resource/:resourceID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM resource`)
			mysql.DB.Exec(`ALTER TABLE resource AUTO_INCREMENT=1`)
		})
	}
}
