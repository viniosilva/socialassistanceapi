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

		personStore := store.NewPersonStore(mysql.DB)
		personService := service.NewPersonService(personStore)
		api := api.NewApi("0.0.0.0:8080", nil, personService)

		// when find all people then returns empty list
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/people", nil)
		api.Gin.ServeHTTP(rec, req)

		if rec.Body.String() != `{"data":[]}` {
			t.Errorf("GET /api/v1/people Body = %v, expected %v", rec.Body.String(), "[]")
		}

		// when create person then create a person
		b, _ := json.Marshal(service.PersonDto{Name: "Test"})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/people", strings.NewReader(string(b)))
		api.Gin.ServeHTTP(rec, req)

		// when update person
		b, _ = json.Marshal(service.PersonDto{Name: "Test updated"})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("PATCH", "/api/v1/people/1", strings.NewReader(string(b)))
		api.Gin.ServeHTTP(rec, req)

		person := rec.Body.String()

		// when find a person by ID then return the person
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/people/1", nil)
		api.Gin.ServeHTTP(rec, req)

		if rec.Body.String() != person {
			t.Errorf("GET /api/v1/people/:personID Body = %v, expected %v", rec.Body.String(), person)
		}

		// when delete a person by ID then return ok
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/api/v1/people/1", nil)
		api.Gin.ServeHTTP(rec, req)

		if rec.Code != 200 {
			t.Errorf("DELETE /api/v1/people Code = %v, expected %v", rec.Body, 200)
		}

		// after
		mysql.DB.Exec(`DELETE FROM people`)
		mysql.DB.Exec(`ALTER TABLE people AUTO_INCREMENT=1`)
	})

}
