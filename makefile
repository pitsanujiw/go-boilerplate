### Variable
RUN=go run
BUILD=go build
DEV=${RUN} main.go
DB=postgres://user:P@ssw0rd1234@localhost:5432/database?sslmode=disable
TZ=Etc/GMT

######################################################################################

### Build

.PHONY: build
build:
	${BUILD} -ldflags="-s -w" -o output/cli.sh ./main.go

######################################################################################

### Dev

.PHONY: api
api: 
	TZ=${TZ} ${DEV} service

.PHONY: worker
worker:
	TZ=${TZ} ${DEV} worker

.PHONY: gateway
gateway:
	TZ=${TZ} ${DEV} gateway

######################################################################################

### Format

.PHONY: format
format: goimport staticcheck

.PHONY: goimport
goimport:
	goimports -w -v ./

.PHONY: staticcheck
staticcheck:
	staticcheck ./...
	go vet ./...

######################################################################################

### Testing

.PHONY: test
test:
	go test -v ./... --cover

.PHONY: test-coverage
test-coverage:
	rm -rf coverage
	mkdir coverage
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	open coverage/coverage.html
	
######################################################################################

### DB

.PHONY: db
db:
	sqlc generate

######################################################################################

### Swagger

.PHONY: swag
swag: swag-api

.PHONY: swag-api
swag-api:
	swag fmt -g ./...
	swag init  --parseDependency -g cmd/service/cmd.go -o cmd/service/docs -t health

######################################################################################

### Utils

.PHONY: setup
setup: mod setup-dependency setup-env cp-env-api cp-env-scheduler cp-env-worker

.PHONY: setup-env
setup-db:
	docker-compose up -d

.PHONY: setup-dependency
setup-dependency:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golang/mock/mockgen@latest
	go install gotest.tools/gotestsum@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: cp-env-api
cp-env-api:
	cp env_default.yaml env.yaml

.PHONY: cp-env-worker
cp-env-worker:
	cp worker_default.yaml worker_env.yaml

######################################################################################

### Migrate
migrate-up:
	migrate -path pkg/database/migration/ -database ${DB} -verbose up

migrate-down:
	migrate -path pkg/database/migration/ -database ${DB} -verbose down ${version}

migrate-fix:
	migrate -path pkg/database/migration/ -database ${dbURI} force ${version}

migrate-new:
	migrate create -ext sql -dir pkg/database/migration/ -seq ${name}

######################################################################################