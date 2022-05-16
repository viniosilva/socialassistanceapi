# Assistencia Social API

Gestão de recursos da assistencia social

## Requiriments

- [Go](https://go.dev/doc/)
- [GNU Make](https://www.gnu.org/software/make/manual/make.html)
- [Docker](https://docs.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Instalation

Run the command:

```shel
make
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

# all
make test
```
