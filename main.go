package main

import (
	"time"

	"github.com/viniosilva/socialassistanceapi/internal/api"
	"github.com/viniosilva/socialassistanceapi/internal/configuration"
	"github.com/viniosilva/socialassistanceapi/internal/service"
	"github.com/viniosilva/socialassistanceapi/internal/store"
)

func main() {
	mysql := configuration.NewMySQL("socialassistanceapi:c8c59046fca24022@tcp(localhost:3306)/socialassistance", time.Minute*1, 3, 3)
	defer mysql.DB.Close()

	personStore := store.NewPersonStore(mysql)
	healthStore := store.NewHealthStore(mysql)
	resourceStore := store.NewResourceStore(mysql)
	addressStore := store.NewAddressStore(mysql)

	personService := service.NewPersonService(personStore)
	healthService := service.NewHealthService(healthStore)
	resourceService := service.NewResourceService(resourceStore)
	addressService := service.NewAddressService(addressStore)

	api := api.NewApi("localhost:8080", healthService, personService, addressService, resourceService)
	api.Start()
}
