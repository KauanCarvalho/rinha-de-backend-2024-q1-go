ifneq (,$(wildcard ./.envrc))
  include .envrc
endif

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

CONNECTION_STRING_PREFIX=postgres://$(DATABASE_USER):$(DATABASE_PASS)@$(DATABASE_HOST):$(DATABASE_PORT)
DATABASE_DSN=$(CONNECTION_STRING_PREFIX)?sslmode=disable
BACKENDFIGHT_DATABASE_DSN=$(CONNECTION_STRING_PREFIX)/$(DATABASE_NAME)?sslmode=disable

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${BACKENDFIGHT_DATABASE_DSN}

# ==================================================================================== #
# DEVELOPMENT AND PRODUCTION
# ==================================================================================== #

## db/check: verify that the database is up and running
.PHONY: db/check
db/check:
	@echo 'Checking connection with database...'
	until nc -z -v -w30 ${DATABASE_HOST} 5432; do \
	  sleep 1; \
	done

## db/drop: drop the database
.PHONY: db/drop
db/drop:
	psql ${DATABASE_DSN} -c "DROP DATABASE IF EXISTS backendfight" || exit 0

## db/create: create the database
.PHONY: db/create
db/create:
	psql $(DATABASE_DSN) -c "CREATE DATABASE backendfight"

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${BACKENDFIGHT_DATABASE_DSN} -verbose up

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...
