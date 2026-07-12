.PHONY: build run test lint clean docker-build migrate-up migrate-down vet

APP_NAME := notification-service
CMD_DIR := ./cmd/server
BIN_DIR := ./bin
MIGRATE_DSN ?= "postgres://postgres:postgres@localhost:5432/notification_service?sslmode=disable"

build:
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

run: build
	$(BIN_DIR)/$(APP_NAME)

test:
	go test -v -race -count=1 ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

clean:
	rm -rf $(BIN_DIR)

docker-build:
	docker build -t $(APP_NAME) -f deploy/Dockerfile .

migrate-up:
	@echo "Run: migrate -path migrations -database $(MIGRATE_DSN) up"

migrate-down:
	@echo "Run: migrate -path migrations -database $(MIGRATE_DSN) down"

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

.PHONY: all
all: vet build test
