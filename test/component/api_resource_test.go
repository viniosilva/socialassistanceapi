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

func TestComponentResourceApiFindAll(t *testing.T) {
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
					INSERT INTO resources (id, created_at, updated_at, name, amount, measurement)
					VALUES (1, ?, ?, 'Test', '1', 'Kg')
				`, date, date)
			},
			expectedCode: 200,
			expectedBody: &service.ResourcesResponse{Data: []service.Resource{{
				ID: 1, CreatedAt: DATE, UpdatedAt: DATE,
				Name:        "Test",
				Amount:      1,
				Measurement: "Kg",
			}}},
		},
		"should return empty list when resource not exists": {
			before:       func(bd *sql.DB) {},
			expectedCode: 200,
			expectedBody: &service.ResourcesResponse{Data: []service.Resource{}},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			resourceStore := store.NewResourceStore(mysql)
			resourceService := service.NewResourceService(resourceStore)
			impl := api.NewApi("0.0.0.0:8080", nil, nil, nil, resourceService)

			cs.before(mysql.DB)

			// when
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/resources", nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.ResourcesResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("GET /api/v1/resources StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("GET /api/v1/resources StatusCode = %v, expected %v", rec.Code, cs.expectedBody)
			}

			// after
			mysql.DB.Exec(`DELETE FROM resources`)
			mysql.DB.Exec(`ALTER TABLE resources AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentResourceApiFindOneByID(t *testing.T) {
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
					INSERT INTO resources (id, created_at, updated_at, name, amount, measurement)
					VALUES (1, ?, ?, 'Test', '1', 'Kg')
				`, date, date)
			},
			inputResourceID: "1",
			expectedCode:    200,
			expectedBody: &service.ResourceResponse{Data: &service.Resource{
				ID: 1, CreatedAt: DATE, UpdatedAt: DATE,
				Name:        "Test",
				Amount:      1,
				Measurement: "Kg",
			}},
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

			resourceStore := store.NewResourceStore(mysql)
			resourceService := service.NewResourceService(resourceStore)
			impl := api.NewApi("0.0.0.0:8080", nil, nil, nil, resourceService)

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
			if rec.Code != cs.expectedCode {
				t.Errorf("GET /api/v1/resources/:resourceID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("GET /api/v1/resources/:resourceID Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("GET /api/v1/resources/:resourceID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM resources`)
			mysql.DB.Exec(`ALTER TABLE resources AUTO_INCREMENT=1`)
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
			inputResource: service.ResourceDto{Name: "Test", Amount: 1, Measurement: "Kg"},
			expectedCode:  201,
			expectedBody: &service.ResourceResponse{Data: &service.Resource{
				ID: 1, Name: "Test", Amount: 1, Measurement: "Kg",
			}},
		},
		"should throw bad request error": {
			expectedCode: 400,
			expectedErr: &api.HttpError{Code: 400, Message: strings.Join([]string{
				"Key: 'ResourceDto.Name' Error:Field validation for 'Name' failed on the 'required' tag",
				"Key: 'ResourceDto.Amount' Error:Field validation for 'Amount' failed on the 'required' tag",
				"Key: 'ResourceDto.Measurement' Error:Field validation for 'Measurement' failed on the 'required' tag",
			}, "\n")},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			resourceStore := store.NewResourceStore(mysql)
			resourceService := service.NewResourceService(resourceStore)
			impl := api.NewApi("0.0.0.0:8080", nil, nil, nil, resourceService)

			// when
			b, _ := json.Marshal(cs.inputResource)
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
			if rec.Code != cs.expectedCode {
				t.Errorf("POST /api/v1/resources StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("POST /api/v1/resources Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("POST /api/v1/resources BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM resources`)
			mysql.DB.Exec(`ALTER TABLE resources AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentResourceApiUpdate(t *testing.T) {
	const DATE = "2000-01-01T12:03:00"

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
				date := strings.Replace(DATE, "T", " ", 1)
				db.Exec(`
					INSERT INTO resources (id, created_at, updated_at, name, amount, measurement)
					VALUES (1, ?, ?, 'Test', '1', 'Kg')
				`, date, date)
			},
			inputResourceID: "1",
			inputResource:   service.ResourceDto{Name: "Test update", Measurement: "l"},
			expectedCode:    200,
			expectedBody: &service.ResourceResponse{Data: &service.Resource{
				ID: 1, CreatedAt: DATE,
				Name:        "Test update",
				Amount:      1,
				Measurement: "l",
			}},
		},
		"should throw bad request error when resourceID id not number": {
			before:          func(db *sql.DB) {},
			inputResourceID: "a",
			expectedCode:    400,
			expectedErr:     &api.HttpError{Code: 400, Message: "invalid resourceID"},
		},
		"should throw bad resquest error": {
			before:          func(db *sql.DB) {},
			inputResourceID: "1",
			expectedCode:    400,
			expectedErr:     &api.HttpError{Code: 400, Message: "empty model: resource"},
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

			resourceStore := store.NewResourceStore(mysql)
			resourceService := service.NewResourceService(resourceStore)
			impl := api.NewApi("0.0.0.0:8080", nil, nil, nil, resourceService)

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputResource)
			url := "/api/v1/resources/" + cs.inputResourceID
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
				t.Errorf("PATCH /api/v1/resources/:resourceID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("PATCH /api/v1/resources/:resourceID Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("PATCH /api/v1/resources/:resourceID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM resources`)
			mysql.DB.Exec(`ALTER TABLE resources AUTO_INCREMENT=1`)
		})
	}
}
