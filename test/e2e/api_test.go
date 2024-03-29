package component

import (
	"encoding/json"
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

func Test_Api(t *testing.T) {
	t.Run("E2E API", func(t *testing.T) {
		// given
		cfg, err := configuration.LoadConfig("../..")
		if err != nil {
			log.Fatal("cannot load config: ", err)
		}

		mysql := infra.MySQLConfigure(cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database, cfg.MySQL.Username,
			cfg.MySQL.Password, time.Duration(cfg.MySQL.ConnMaxLifetimeMs), cfg.MySQL.MaxOpenConns, cfg.MySQL.MaxIdleConns)
		defer mysql.DB.Close()

		personRepository := &repository.PersonRepositoryImpl{DB: mysql}
		resourceRepository := &repository.ResourceRepositoryImpl{DB: mysql}
		familyRepository := &repository.FamilyRepositoryImpl{DB: mysql}
		donateResourceRepository := &repository.DonateResourceRepositoryImpl{DB: mysql}

		personService := &service.PersonServiceImpl{PersonRepository: personRepository}
		resourceService := &service.ResourceServiceImpl{ResourceRepository: resourceRepository}
		familyService := &service.FamilyServiceImpl{FamilyRepository: familyRepository}
		donateResourceService := &service.DonateResourceServiceImpl{DonateResourceRepository: donateResourceRepository}

		impl := &api.ApiImpl{
			Addr:                  "0.0.0.0:8080",
			PersonService:         personService,
			FamilyService:         familyService,
			ResourceService:       resourceService,
			DonateResourceService: donateResourceService,
		}
		impl.Configure()

		// when find all persons then return OK
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/persons", nil)

		impl.Gin.ServeHTTP(rec, req)
		assert.Equal(t, rec.Code, http.StatusOK)

		// when find all families then return OK
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/families", nil)

		impl.Gin.ServeHTTP(rec, req)
		assert.Equal(t, rec.Code, http.StatusOK)

		// when find all resources then return OK
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/resources", nil)
		impl.Gin.ServeHTTP(rec, req)
		assert.Equal(t, rec.Code, http.StatusOK)

		// when create family then return Created
		b, _ := json.Marshal(service.FamilyCreateDto{
			Name:         "Sauro",
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
		req, _ = http.NewRequest("POST", "/api/v1/families", strings.NewReader(string(b)))

		impl.Gin.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// when create person then return Created
		b, _ = json.Marshal(service.PersonCreateDto{FamilyID: 1, Name: "Test"})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/persons", strings.NewReader(string(b)))
		impl.Gin.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		// when create resource then return Created
		b, _ = json.Marshal(service.CreateResourceDto{Name: "Test", Amount: 1, Measurement: "l", Quantity: 10})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/resources", strings.NewReader(string(b)))
		impl.Gin.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		// when update person then return NoContent
		b, _ = json.Marshal(service.PersonCreateDto{Name: "Test Update"})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("PATCH", "/api/v1/persons/1", strings.NewReader(string(b)))
		impl.Gin.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		// when update family then return NoContent
		b, _ = json.Marshal(service.FamilyCreateDto{
			State:        "RS",
			City:         "Porto Alegre",
			Neighborhood: "Hípica",
			Street:       "R. J",
			Number:       "1",
			Zipcode:      "91755450",
		})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("PATCH", "/api/v1/families/1", strings.NewReader(string(b)))
		impl.Gin.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		// when update resource then return NoContent
		b, _ = json.Marshal(service.UpdateResourceDto{Measurement: "Kg"})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("PATCH", "/api/v1/resources/1", strings.NewReader(string(b)))
		impl.Gin.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		// when update resource quantity then return NoContent
		b, _ = json.Marshal(service.UpdateResourceQuantityDto{Quantity: 1})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("PATCH", "/api/v1/resources/1/quantity", strings.NewReader(string(b)))
		impl.Gin.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		// when donate resource then return NoContent
		b, _ = json.Marshal(service.DonateResourceDonateDto{FamilyID: 1, Quantity: 1})
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/resources/1/donate", strings.NewReader(string(b)))
		impl.Gin.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		// when find a person by ID then return status OK
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/persons/1", nil)
		impl.Gin.ServeHTTP(rec, req)
		assert.Equal(t, rec.Code, http.StatusOK)

		// when find an family by ID then return status OK
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/families/1", nil)
		impl.Gin.ServeHTTP(rec, req)
		assert.Equal(t, rec.Code, http.StatusOK)

		// when find a resource by ID then return status OK
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/resources/1", nil)
		impl.Gin.ServeHTTP(rec, req)
		assert.Equal(t, rec.Code, http.StatusOK)

		// when delete a person by ID then return NoContent
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/api/v1/persons/1", nil)
		impl.Gin.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// when delete an family by ID then return NoContent
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/api/v1/families/1", nil)
		impl.Gin.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// after
		mysql.DB.Exec(`DELETE FROM persons`)
		mysql.DB.Exec(`DELETE FROM families`)
		mysql.DB.Exec(`DELETE FROM resources`)
		mysql.DB.Exec(`DELETE FROM resources_to_families`)
		mysql.DB.Exec(`ALTER TABLE persons AUTO_INCREMENT=1`)
		mysql.DB.Exec(`ALTER TABLE families AUTO_INCREMENT=1`)
		mysql.DB.Exec(`ALTER TABLE resources AUTO_INCREMENT=1`)
		mysql.DB.Exec(`ALTER TABLE resources_to_families AUTO_INCREMENT=1`)
	})
}
