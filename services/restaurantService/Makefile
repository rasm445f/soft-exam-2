# read .env variables for use with commands. Maybe not necessary
include .env
.EXPORT_ALL_VARIABLES:

# Declare all targets as phony to avoid conflicts with file names
.PHONY: default run build docs sqlc migrate-new migrate-up migrate-down test test-verbose test-cover test-cover-html test-mutation

# Default target to show available commands
default:
	@echo "Available targets:"
	@echo "  run			# Run with air"
	@echo "  build			# Build the project"
	@echo "  docs 			# Generate docs with swag"
	@echo "  sqlc			# Generate sqlc boilerplate"
	@echo "  migrate-new		# Create new sql migration"
	@echo "  migrate-up		# Migrate the db"
	@echo "  migrate-down		# Roll back the db"
	@echo "  test			# Run tests"
	@echo "  test-verbose		# Run test with verbose flag"
	@echo "  test-cover		# Show test coverage"
	@echo "  test-cover-html	# Show test coverage in the browser"
	@echo "  test-mutation		# Run mutation testing with gremlins"


run-go:
	DBSTRING=$(DBSTRING) go run main.go
# Unix-like systems (Linux/Mac)
run-unix:
	DBSTRING=$(DBSTRING) air --build.cmd "go build -o bin/api" --build.bin "./bin/api"
	
# Windows
run-windows:
	DBSTRING="postgres://postgres:test@127.0.0.1:5434/restaurant?sslmode=disable" \
		air --build.cmd "go build -o bin\\api.exe" --build.bin "bin\\api.exe"

# Main run target
run:
ifeq ($(OS),Windows_NT)
	$(MAKE) run-windows
else
	$(MAKE) run-unix
endif

build:
	go build -o bin/api

# Used for pgtype.Numeric fields
docs:
	swag init --parseDependency --parseInternal

sqlc:
	sqlc generate

migrate-new:
	@read -p "Enter migration name: " NAME; \
	goose create -dir db/migrations $$NAME sql

migrate-up: # execute migrations
	@goose -dir db/migrations postgres "$(DBSTRING)" up

migrate-down:
	@goose -dir db/migrations postgres "$(DBSTRING)" down

test:
	go test ./...

test-verbose:
	go test ./... -v

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-cover-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

test-mutation:
	gremlins unleash --output=mutation.json
