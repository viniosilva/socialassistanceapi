# AssistenciaSocial API

Management of social assistance resources

## Requiriments

- [Go](https://go.dev/doc/)
- [Gomock](https://github.com/golang/mock)
- [GNU Make](https://www.gnu.org/software/make/manual/make.html)
- [Docker](https://docs.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Instalation

Run the command:

```shel
make
```

## Configuration

To run locally, create the `.env` file as in the example below:

```
GIN_MODE=debug
MIGRATION_URL=mysql://socialassistanceapi:c8c59046fca24022@tcp\(localhost:3306\)/socialassistance
MYSQL_PASSWORD=c8c59046fca24022
```

## Migrations

Run the command:

```shel
make migrate
```

## Running

Run the command:

```shel
make run
```

To see the documentation, just enter here: `http://localhost:8080/swagger/index.html`

## Tests

Run the commands:

```shel
make test/unit
make test/component
make test/e2e

# all
make test
```
