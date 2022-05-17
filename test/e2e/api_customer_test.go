package component

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/viniosilva/socialassistanceapi/internal/api"
	"github.com/viniosilva/socialassistanceapi/internal/configuration"
	"github.com/viniosilva/socialassistanceapi/internal/service"
	"github.com/viniosilva/socialassistanceapi/internal/store"
)

func TestE2ECustomerApi(t *testing.T) {
	t.Run("E2E Customer API", func(t *testing.T) {
		// given
		mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
		defer mysql.DB.Close()

		customerStore := store.NewCustomerStore(mysql.DB)
		customerService := service.NewCustomerService(customerStore)
		api := api.NewApi("0.0.0.0:8080", nil, customerService)

		// when find all customers then returns empty list
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/customers", nil)
		api.Gin.ServeHTTP(rec, req)

		if rec.Body.String() != `{"data":[]}` {
			t.Errorf("GET /api/v1/customers Body = %v, expected %v", rec.Body.String(), "[]")
		}

		// when create customer then create a customer
		b, _ := json.Marshal(service.CustomerDto{Name: "Test"})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/customers", strings.NewReader(string(b)))
		api.Gin.ServeHTTP(rec, req)

		// when update customer
		b, _ = json.Marshal(service.CustomerDto{Name: "Test updated"})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("PATCH", "/api/v1/customers/1", strings.NewReader(string(b)))
		api.Gin.ServeHTTP(rec, req)

		customer := rec.Body.String()

		// when find a customer by ID then return the customer
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/customers/1", nil)
		api.Gin.ServeHTTP(rec, req)

		if rec.Body.String() != customer {
			t.Errorf("POST /api/v1/customers Body = %v, expected %v", rec.Body.String(), customer)
		}

		// after
		mysql.DB.Exec(`DELETE FROM customers`)
		mysql.DB.Exec(`ALTER TABLE customers AUTO_INCREMENT=1`)
	})

}
