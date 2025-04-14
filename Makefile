include .env

LOCAL_BIN := $(CURDIR)/bin

install-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0


install-grpc-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-grpc-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


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



generate-pvz-api:
	mkdir -p pkg/pvz_v1
	PATH=$(LOCAL_BIN):$$PATH \
	protoc --proto_path=proto \
	--go_out=pkg/pvz_v1 --go_opt=paths=source_relative \
	--go-grpc_out=pkg/pvz_v1 --go-grpc_opt=paths=source_relative \
	proto/pvz/pvz.proto

