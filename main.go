package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/viniosilva/socialassistanceapi/internal/api"
	"github.com/viniosilva/socialassistanceapi/internal/configuration"
	"github.com/viniosilva/socialassistanceapi/internal/repository"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	cfg, err := configuration.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	mysql := configuration.MySQLConfigure(cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database, cfg.MySQL.Username,
		cfg.MySQL.Password, time.Duration(cfg.MySQL.ConnMaxLifetimeMs), cfg.MySQL.MaxOpenConns, cfg.MySQL.MaxIdleConns)
	defer mysql.DB.Close()

	healthRepository := &repository.HealthRepositoryImpl{DB: mysql}
	personRepository := &repository.PersonRepositoryImpl{DB: mysql}
	resourceRepository := &repository.ResourceRepositoryImpl{DB: mysql}
	addressRepository := &repository.AddressRepositoryImpl{DB: mysql}
	donateResourceRepository := &repository.DonateResourceRepositoryImpl{DB: mysql}

	healthService := &service.HealthServiceImpl{HealthRepository: healthRepository}
	personService := &service.PersonServiceImpl{PersonRepository: personRepository}
	resourceService := &service.ResourceServiceImpl{ResourceRepository: resourceRepository}
	addressService := &service.AddressServiceImpl{AddressRepository: addressRepository}
	donateResourceService := &service.DonateResourceServiceImpl{DonateResourceRepository: donateResourceRepository}

	api := &api.ApiImpl{
		Addr:                  fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port),
		HealthService:         healthService,
		PersonService:         personService,
		AddressService:        addressService,
		ResourceService:       resourceService,
		DonateResourceService: donateResourceService,
	}

	api.Configure()
	api.Start()
}
