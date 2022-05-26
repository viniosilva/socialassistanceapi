include .env
export

all:
	go get
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go generate ./...
	swag init

update: mock docs
	go get -u

.PHONY: mock
mock:
	go generate ./...

.PHONY: docs
docs:
	swag init

infra/up:
	docker-compose up --build -d

infra/down:
	docker-compose down --remove-orphans

run: infra/up
	go run main.go

migrate: infra/up
	migrate -source file://db/migrations -database ${MIGRATION_URL} up

migrate/new:
	migrate create -ext sql -dir db/migrations -seq ${name}

.PHONY: test/unit
test/unit:
	go test ./test/unit/...

.PHONY: test/component
test/component: infra/up
	go test ./test/component/...

.PHONY: test/e2e
test/e2e: infra/up
	go test ./test/e2e/...

.PHONY: test
test: infra/up test/unit test/component test/e2e
