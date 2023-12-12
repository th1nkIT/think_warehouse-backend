
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
GOLANGCI_CMD := $(shell command -v golangci-lint 2> /dev/null)
PKGS := $(shell go list ./... | grep -v /vendor/ | grep -v /internal/mock)
ALL_PACKAGES := $(shell go list ./... | grep -v /vendor/ | grep -v /internal/mock)
PG_MIGRATIONS_FOLDER=./scripts/pgbo/migrations
PG_DB_URL=postgresql://postgres:postgres@127.0.01:5432/think_laundry?sslmode=disable

CMD_SQLC := $(shell command -v sqlc 2> /dev/null)
CMD_MIGRATE := $(shell command -v migrate 2> /dev/null)

check-migrations-cmd:
ifndef CMD_MIGRATE
	$(error "migrate is not installed, see: https://github.com/golang-migrate/migrate")
endif
ifndef CMD_SQLC
	$(error "sqlc is not installed, see: github.com/kyleconroy/sqlc)")
endif

check-golangci:
ifndef GOLANGCI_CMD
	$(error "Please install golangci linters from https://golangci-lint.run/usage/install/")
endif

lint: check-golangci fmt
	@echo -e "$(OK_COLOR)==> linting projects$(NO_COLOR)..."
	@golangci-lint run --fix

fmt:
	@go fmt $(ALL_PACKAGES)

test:
	@go test -coverprofile=test_coverage.out $(PKGS)
	go tool cover -html=test_coverage.out -o test_coverage.html
	rm test_coverage.out
	@echo -e "$(OK_COLOR)==> Open test_coverage.html file on your web browser for detailed coverage$(NO_COLOR)..."

deps:
	go mod tidy && go mod vendor

gen.pg.models: check-migrations-cmd
	cd ./scripts/pgbo && sqlc generate

gen.pg.migration: check-migrations-cmd
	migrate create -ext sql -dir $(PG_MIGRATIONS_FOLDER) --seq $(name)

pg.migrate.up: check-migrations-cmd
	migrate -path $(PG_MIGRATIONS_FOLDER) -database $(PG_DB_URL) --verbose up

pg.migrate.down: check-migrations-cmd
	migrate -path $(PG_MIGRATIONS_FOLDER) -database $(PG_DB_URL) --verbose down

GRPCUI_CMD := $(shell command -v grpcui 2> /dev/null)
PROTOC_CMD := $(shell command -v protoc 3> /dev/null)
GRPC_PORT = 10237

grpc.gen.proto:
ifndef PROTOC_CMD
	$(error "protoc-gen-go is not installed. Run command 'go get -u google.golang.org/protobuf/proto && go install github.com/golang/protobuf/protoc-gen-go'")
endif
	@echo -e "$(OK_COLOR)==> Generate proto objects to pkg/grpc/v1$(NO_COLOR)..."
	@protoc --proto_path=pkg/grpc --go_out=plugins=grpc:./ \
		pkg/grpc/*.proto
	@echo -e "$(OK_COLOR)==> Done$(NO_COLOR)..."

grpc-ui:
ifndef GRPCUI_CMD
	$(error "Please install grpcui first! See: https://github.com/fullstorydev/grpcui")
endif
	grpcui -plaintext localhost:$(GRPC_PORT)

run-service-local:
	go run -mod=vendor cmd/api/application.go
