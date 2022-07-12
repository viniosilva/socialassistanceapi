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

func TestE2EApi(t *testing.T) {
	t.Run("E2E API", func(t *testing.T) {
		// given
		mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
		defer mysql.DB.Close()

		personStore := store.NewPersonStore(mysql)
		personService := service.NewPersonService(personStore)
		addressStore := store.NewAddressStore(mysql)
		addressService := service.NewAddressService(addressStore)
		resourceStore := store.NewResourceStore(mysql)
		resourceService := service.NewResourceService(resourceStore)

		api := api.NewApi("0.0.0.0:8080", nil, personService, addressService, resourceService)

		// when find all persons then returns empty list
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/persons", nil)
		api.Gin.ServeHTTP(rec, req)

		if rec.Body.String() != `{"data":[]}` {
			t.Errorf("GET /api/v1/persons Body = %v, expected %v", rec.Body.String(), "[]")
		}

		// when find all addresses then returns empty list
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/addresses", nil)
		api.Gin.ServeHTTP(rec, req)

		if rec.Body.String() != `{"data":[]}` {
			t.Errorf("GET /api/v1/addresses Body = %v, expected %v", rec.Body.String(), "[]")
		}

		// when find all resources then returns empty list
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/resources", nil)
		api.Gin.ServeHTTP(rec, req)

		if rec.Body.String() != `{"data":[]}` {
			t.Errorf("GET /api/v1/resources Body = %v, expected %v", rec.Body.String(), "[]")
		}

		// when create person then create a person
		b, _ := json.Marshal(service.PersonDto{Name: "Test"})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/persons", strings.NewReader(string(b)))
		api.Gin.ServeHTTP(rec, req)

		if rec.Code != 201 {
			t.Errorf("POST /api/v1/persons Code = %v, expected %v", rec.Body, 201)
		}

		// when create address then create an address
		b, _ = json.Marshal(service.AddressDto{
			Country:      "BR",
			State:        "SP",
			City:         "São Paulo",
			Neighborhood: "Pq. Novo Mundo",
			Street:       "R. Sd. Teodoro Francisco Ribeiro",
			Number:       "1",
			Complement:   "1",
			Zipcode:      "02180110",
		})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/addresses", strings.NewReader(string(b)))
		api.Gin.ServeHTTP(rec, req)

		if rec.Code != 201 {
			t.Errorf("POST /api/v1/addresses Code = %v, expected %v", rec.Body, 201)
		}

		// when create resource then create a resource
		b, _ = json.Marshal(service.ResourceDto{Name: "Test", Amount: 10, Measurement: "l"})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/resources", strings.NewReader(string(b)))
		api.Gin.ServeHTTP(rec, req)

		if rec.Code != 201 {
			t.Errorf("POST /api/v1/resources Code = %v, expected %v", rec.Body, 201)
		}

		// when update person
		b, _ = json.Marshal(service.PersonDto{Name: "Test Update"})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("PATCH", "/api/v1/persons/1", strings.NewReader(string(b)))
		api.Gin.ServeHTTP(rec, req)

		if rec.Code != 200 {
			t.Errorf("PATH /api/v1/addresses Code = %v, expected %v", rec.Body, 200)
		}
		person := rec.Body.String()

		// when update address
		b, _ = json.Marshal(service.AddressDto{
			State:        "RS",
			City:         "Porto Alegre",
			Neighborhood: "Hípica",
			Street:       "R. J",
			Number:       "1",
			Zipcode:      "91755450",
		})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("PATCH", "/api/v1/addresses/1", strings.NewReader(string(b)))
		api.Gin.ServeHTTP(rec, req)

		if rec.Code != 200 {
			t.Errorf("PATH /api/v1/addresses Code = %v, expected %v", rec.Body, 200)
		}
		address := rec.Body.String()

		// when update resource
		b, _ = json.Marshal(service.ResourceUpdateDto{Measurement: "Kg"})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("PATCH", "/api/v1/resources/1", strings.NewReader(string(b)))
		api.Gin.ServeHTTP(rec, req)

		if rec.Code != 200 {
			t.Errorf("PATH /api/v1/resources Code = %v, expected %v", rec.Body, 200)
		}
		resource := rec.Body.String()

		// when find a person by ID then return the person
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/persons/1", nil)
		api.Gin.ServeHTTP(rec, req)

		if rec.Body.String() != person {
			t.Errorf("GET /api/v1/persons/:personID Body = %v, expected %v", rec.Body.String(), person)
		}

		// when find an address by ID then return the address
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/addresses/1", nil)
		api.Gin.ServeHTTP(rec, req)

		if rec.Body.String() != address {
			t.Errorf("GET /api/v1/addresses/:addressID Body = %v, expected %v", rec.Body.String(), address)
		}

		// when find a resource by ID then return the resource
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/resources/1", nil)
		api.Gin.ServeHTTP(rec, req)

		if rec.Body.String() != resource {
			t.Errorf("GET /api/v1/resources/:resourceID Body = %v, expected %v", rec.Body.String(), resource)
		}

		// when transfer amount to resource by ID then return the resource
		b, _ = json.Marshal(service.ResourceDto{Name: "Test", Amount: 10, Measurement: "l"})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/resources/1/amount/transfer", strings.NewReader(string(b)))
		api.Gin.ServeHTTP(rec, req)

		if rec.Code != 200 {
			t.Errorf("POST /api/v1/resources/:resourceID/amount/transfer Code = %v, expected %v", rec.Body, 200)
		}

		// when delete a person by ID then return ok
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/api/v1/persons/1", nil)
		api.Gin.ServeHTTP(rec, req)

		if rec.Code != 204 {
			t.Errorf("DELETE /api/v1/persons/:addressID Code = %v, expected %v", rec.Body, 204)
		}

		// when delete an address by ID then return ok
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/api/v1/addresses/1", nil)
		api.Gin.ServeHTTP(rec, req)

		if rec.Code != 204 {
			t.Errorf("DELETE /api/v1/addresses/:addressID Code = %v, expected %v", rec.Body, 204)
		}

		// after
		mysql.DB.Exec(`DELETE FROM persons`)
		mysql.DB.Exec(`DELETE FROM addresses`)
		mysql.DB.Exec(`DELETE FROM resources`)
		mysql.DB.Exec(`ALTER TABLE persons AUTO_INCREMENT=1`)
		mysql.DB.Exec(`ALTER TABLE addresses AUTO_INCREMENT=1`)
		mysql.DB.Exec(`ALTER TABLE resources AUTO_INCREMENT=1`)
	})
}
