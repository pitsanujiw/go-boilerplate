# Golang Boilerplate

## Build & Run Service

for local machine please install go version 1.8+ [download](https://go.dev/doc/install) and follow command below to start this service
and please install docker desktop for support docker compose [download](https://www.docker.com/products/docker-desktop)

### Setup before develop service

```bash
$ make setup
```

## Build service

```bash
$ make build
```

### How to run service

for running a service for a part of cli is service
```bash
$ make api
```

for running a service for a part of worker
```bash
$ make worker
```

## Generate Swagger

```bash
$ make swag
```


## Format and lint checking
```bash
$ make format
```

## Testing

```bash
$ make test
```
or testing coverage

```bash
$ make test-coverage
```

## Generate Go DB Code From SQL

Generates Go code from SQL queries defined in the `.sql` files in the `pkg/database/queries/` directory.

```shell
sqlc generate
```

Requires [Go SQLC 1.18](https://docs.sqlc.dev/en/latest/overview/install.html)

The settings are in `sqlc.yaml`

### Setup and Installation

```bash
 	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

#### for Linux

```bush
    curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey| apt-key add -
    echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
    apt-get update
    apt-get install -y migrate
```


### How to create migration file
```bash
make postgres-migrate-new name=${name}                                   
```

### How to migrate up

```bash
make migrate-up
```

### How to migrate down

```bash
make migrate-down version=${version}
```

### How to fix version

```bash
make migrate-fix version=${version}
```