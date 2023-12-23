NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
GOLANGCI_CMD := $(shell command -v golangci-lint 2> /dev/null)
PKGS := $(shell go list ./... | grep -v /vendor/ | grep -v /internal/mock)
ALL_PACKAGES := $(shell go list ./... | grep -v /vendor/ | grep -v /internal/mock)

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