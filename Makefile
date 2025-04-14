include .env

LOCAL_BIN := $(CURDIR)/bin

install-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

install-oapi:
	GOBIN=$(LOCAL_BIN) go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

install-mockgen:
	go install github.com/gojuno/minimock/v3/cmd/minimock@latest

install-pgx:
	go get github.com/jackc/pgx/v5/pgxpool@latest

install-squirrel:
	go get github.com/Masterminds/squirrel@latest

install-deps: install-goose install-oapi install-mockgen install-pgx install-squirrel

LOCAL_MIGRATION_DIR = $(MIGRATION_DIR)
LOCAL_MIGRATION_DSN = "host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

local-migration-status:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) down -v

OPENAPI_FILE = swagger.yaml
OAPI_PACKAGE = api
OAPI_OUT_DIR = internal/api
OAPI_OUTPUT = $(OAPI_OUT_DIR)/api.gen.go

generate-api:
	@echo "⚙️  Generating API from $(OPENAPI_FILE)..."
	@mkdir -p $(OAPI_OUT_DIR)
	$(LOCAL_BIN)/oapi-codegen -generate types,chi-server -package $(OAPI_PACKAGE) -o $(OAPI_OUTPUT) $(OPENAPI_FILE)
