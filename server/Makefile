APP_NAME := drawing-api

include .env
export $(shell sed 's/=.*//' .env) 

.PHONY: build test run clean postgresinit postgres createdb dropdb migrateup migratedown


all: build
dev:
	go run cmd/drawing-app/main.go

build:
	@echo "Building..."
	go build -o $(APP_NAME) .

test:
	@echo "Running tests..."
	go test -v ./...

run: build
	@echo "Running..."
	./$(APP_NAME)

clean:
	@echo "Cleaning..."
	go clean
	rm -f $(APP_NAME)
	
postgresinit:
	docker run --name postgres15 -p $(POSTGRES_PORT):5432 --env-file .env -d postgres:15-alpine

postgres:
	docker exec -it postgres15 psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

createdb:
	docker exec -it postgres15 createdb password=$(POSTGRES_PASSWORD) --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(POSTGRES_DB)

dropdb:
	docker exec -it postgres15 dropdb 

migrateup:
	migrate -path internal/db/migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/db/migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose down

migrate-force:
	@read -p "Enter version to force: " v; \
	migrate -path internal/db/migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" force $$v
