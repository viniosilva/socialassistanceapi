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

const DATE = "2000-01-01"

func TestComponentCustomerApiFindAll(t *testing.T) {
	cases := map[string]struct {
		before       func(db *sql.DB)
		expectedCode int
		expectedBody *service.CustomersResponse
	}{
		"should return customer list when customers exists": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO customers (id, created_at, updated_at, name)
					VALUES (1, ?, ?, 'Test')
				`, DATE, DATE)
			},
			expectedCode: 200,
			expectedBody: &service.CustomersResponse{Data: []service.Customer{{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}}},
		},
		"should return empty list when customers not exists": {
			before:       func(db *sql.DB) {},
			expectedCode: 200,
			expectedBody: &service.CustomersResponse{Data: []service.Customer{}},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			customerStore := store.NewCustomerStore(mysql.DB)
			customerService := service.NewCustomerService(customerStore)
			impl := api.NewApi("0.0.0.0:8080", nil, customerService)

			cs.before(mysql.DB)

			// when
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/customers", nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.CustomersResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("GET /api/v1/customers StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("GET /api/v1/customers Body = %v, expected %v", body, cs.expectedBody)
			}

			// after
			mysql.DB.Exec(`DELETE FROM customers`)
			mysql.DB.Exec(`ALTER TABLE customers AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentCustomerApiFindOneByID(t *testing.T) {
	cases := map[string]struct {
		before          func(db *sql.DB)
		inputCustomerID string
		expectedCode    int
		expectedBody    *service.CustomerResponse
		expectedErr     *api.HttpError
	}{
		"should return customer when customers exists": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO customers (id, created_at, updated_at, name)
					VALUES (1, ?, ?, 'Test')
				`, DATE, DATE)
			},
			inputCustomerID: "1",
			expectedCode:    200,
			expectedBody:    &service.CustomerResponse{Data: &service.Customer{ID: 1, CreatedAt: DATE, UpdatedAt: DATE, Name: "Test"}},
		},
		"should throw bad request error when customerID is not a number": {
			before:          func(db *sql.DB) {},
			inputCustomerID: "a",
			expectedCode:    400,
			expectedErr:     &api.HttpError{Code: 400, Message: "invalid customerID"},
		},
		"should throw not found error when customers not exists": {
			before:          func(db *sql.DB) {},
			inputCustomerID: "1",
			expectedCode:    404,
			expectedBody:    &service.CustomerResponse{},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			customerStore := store.NewCustomerStore(mysql.DB)
			customerService := service.NewCustomerService(customerStore)
			impl := api.NewApi("0.0.0.0:8080", nil, customerService)

			cs.before(mysql.DB)

			// when
			url := "/api/v1/customers/" + cs.inputCustomerID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.CustomerResponse
			json.Unmarshal(rec.Body.Bytes(), &body)

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("GET /api/v1/customers/:customerID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("GET /api/v1/customers/:customerID Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("GET /api/v1/customers/:customerID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM customers`)
			mysql.DB.Exec(`ALTER TABLE customers AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentCustomerApiCreate(t *testing.T) {
	cases := map[string]struct {
		inputCustomer service.CustomerDto
		expectedCode  int
		expectedBody  *service.CustomerResponse
		expectedErr   *api.HttpError
	}{
		"should return created customer": {
			inputCustomer: service.CustomerDto{Name: "Test"},
			expectedCode:  201,
			expectedBody:  &service.CustomerResponse{Data: &service.Customer{ID: 1, Name: "Test"}},
		},
		"should throw bad request error": {
			expectedCode: 400,
			expectedErr:  &api.HttpError{Code: 400, Message: "Key: 'CustomerDto.Name' Error:Field validation for 'Name' failed on the 'required' tag"},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			customerStore := store.NewCustomerStore(mysql.DB)
			customerService := service.NewCustomerService(customerStore)
			impl := api.NewApi("0.0.0.0:8080", nil, customerService)

			// when
			b, _ := json.Marshal(cs.inputCustomer)
			url := "/api/v1/customers"
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var body *service.CustomerResponse
			json.Unmarshal(rec.Body.Bytes(), &body)
			if body.Data != nil {
				body.Data.CreatedAt = ""
				body.Data.UpdatedAt = ""
			}

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("POST /api/v1/customers StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("POST /api/v1/customers Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("POST /api/v1/customers BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM customers`)
			mysql.DB.Exec(`ALTER TABLE customers AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentCustomerApiUpdate(t *testing.T) {
	cases := map[string]struct {
		before          func(db *sql.DB)
		inputCustomerID string
		inputCustomer   service.CustomerDto
		expectedCode    int
		expectedBody    *service.CustomerResponse
		expectedErr     *api.HttpError
	}{
		"should return updated customer": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO customers (id, created_at, updated_at, name)
					VALUES (1, ?, ?, 'Test')
				`, DATE, DATE)
			},
			inputCustomerID: "1",
			inputCustomer:   service.CustomerDto{Name: "Test update"},
			expectedCode:    200,
			expectedBody:    &service.CustomerResponse{Data: &service.Customer{ID: 1, CreatedAt: DATE, Name: "Test update"}},
		},
		"should throw bad request error when customerID is not a number": {
			before:          func(db *sql.DB) {},
			inputCustomerID: "a",
			expectedCode:    400,
			expectedErr:     &api.HttpError{Code: 400, Message: "invalid customerID"},
		},
		"should throw bad request error": {
			before:          func(db *sql.DB) {},
			inputCustomerID: "1",
			expectedCode:    400,
			expectedErr:     &api.HttpError{Code: 400, Message: "Key: 'CustomerDto.Name' Error:Field validation for 'Name' failed on the 'required' tag"},
		},
		"should throw not found error when customers not exists": {
			before:          func(db *sql.DB) {},
			inputCustomerID: "1",
			inputCustomer:   service.CustomerDto{Name: "Test update"},
			expectedCode:    404,
			expectedBody:    &service.CustomerResponse{},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			customerStore := store.NewCustomerStore(mysql.DB)
			customerService := service.NewCustomerService(customerStore)
			impl := api.NewApi("0.0.0.0:8080", nil, customerService)

			cs.before(mysql.DB)

			// when
			b, _ := json.Marshal(cs.inputCustomer)
			url := "/api/v1/customers/" + cs.inputCustomerID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(b)))
			impl.Gin.ServeHTTP(rec, req)

			var body *service.CustomerResponse
			json.Unmarshal(rec.Body.Bytes(), &body)
			if body.Data != nil {
				body.Data.UpdatedAt = ""
			}

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("PATCH /api/v1/customers/:customerID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("PATCH /api/v1/customers/:customerID Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("PATCH /api/v1/customers/:customerID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM customers`)
			mysql.DB.Exec(`ALTER TABLE customers AUTO_INCREMENT=1`)
		})
	}
}

func TestComponentCustomerApiDelete(t *testing.T) {
	cases := map[string]struct {
		before          func(db *sql.DB)
		inputCustomerID string
		expectedCode    int
		expectedBody    *service.CustomerResponse
		expectedErr     *api.HttpError
	}{
		"should be successfull": {
			before: func(db *sql.DB) {
				db.Exec(`
					INSERT INTO customers (id, created_at, updated_at, name)
					VALUES (1, ?, ?, 'Test')
				`, DATE, DATE)
			},
			inputCustomerID: "1",
			expectedCode:    200,
			expectedBody:    &service.CustomerResponse{},
		},
		"should throw bad request error when customerID is not a number": {
			before:          func(db *sql.DB) {},
			inputCustomerID: "a",
			expectedCode:    400,
			expectedErr:     &api.HttpError{Code: 400, Message: "invalid customerID"},
		},
		"should be successfull when customers not exists": {
			before:          func(db *sql.DB) {},
			inputCustomerID: "1",
			expectedCode:    200,
			expectedBody:    &service.CustomerResponse{},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
			defer mysql.DB.Close()

			customerStore := store.NewCustomerStore(mysql.DB)
			customerService := service.NewCustomerService(customerStore)
			impl := api.NewApi("0.0.0.0:8080", nil, customerService)

			cs.before(mysql.DB)

			// when
			url := "/api/v1/customers/" + cs.inputCustomerID
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", url, nil)
			impl.Gin.ServeHTTP(rec, req)

			var body *service.CustomerResponse
			json.Unmarshal(rec.Body.Bytes(), &body)
			if body.Data != nil {
				body.Data.UpdatedAt = ""
			}

			var httpError *api.HttpError
			json.Unmarshal(rec.Body.Bytes(), &httpError)

			// then
			if rec.Code != cs.expectedCode {
				t.Errorf("PATCH /api/v1/customers/:customerID StatusCode = %v, expected %v", rec.Code, cs.expectedCode)
			}
			if cs.expectedBody != nil && !reflect.DeepEqual(body, cs.expectedBody) {
				t.Errorf("PATCH /api/v1/customers/:customerID Body = %v, expected %v", body, cs.expectedBody)
			}
			if cs.expectedErr != nil && !reflect.DeepEqual(httpError, cs.expectedErr) {
				t.Errorf("PATCH /api/v1/customers/:customerID BodyErr = %v, expected %v", httpError, cs.expectedErr)
			}

			// after
			mysql.DB.Exec(`DELETE FROM customers`)
			mysql.DB.Exec(`ALTER TABLE customers AUTO_INCREMENT=1`)
		})
	}
}
