# 1.3 Go Package Layout

Broadly speaking, it is divided into 3 main package folders
* `cmd` contains main program
* `internal` app business logic
* `pkg` for packages that can be imported by other go projects / repos

## Service Package Design

```console
user_service_repo
    cmd/ <- boleh import package /internal & /pkg
        api/ <- bisa run echo http & grpc di port yg sama, run /internal/{grpc,http}
            main.go
        pubsubworker/ <- bisa run multiple pubsub subscription worker, run /internal/worker
            main.go
    internal/ <-- boleh import /pkg, tidak boleh diimport secara go-module dari repo luar
        service/ <- boleh import repository, api_client, tidak boleh import /internal/http & internal/grpc
            user_password_service.go
            ...
            service.go
        http/ <- expose /internal/service functions. Request/response body struct define di /pkg/http
            handler/
                user_password_handler.go
                ...
                handler.go
            my_custom_middleware.go
            my_custom2_middleware.go
            ...
            server.go
        grpc/ <- expose /internal/service functions
            user_password_handler.go
            server.go
        worker/ <- import repository, api_client
            pubsub_worker.go
            user_updated_event_worker.go
            new_oba_created_event_worker.go
            ...
        repository/ <- group by datasource / host url instance
            userpostgres/ <- should be generated from sqlc
                user_password_query.sql.go
                user_profile_query.sql.go
                ...
            marketplacepostgres/ <- also should be generated from sqlc
                user_email.sql.go
                user_address.sql.go
                ...
            registrationmongo/
                user_registration_repository.go
                ...
            membermongo/
                officer_attribute_repository.go
                ...
            redis/
                redis_repository.go
                user_cart_repository.go
                ...
            sessionredis/
                redis_repository.go
                user_session_repository.go
                ...
            ...
        config/
            config_keys.go <- define constant config keys here
            config.go <- contains httpClient, db, redis or pubsub connection
            constructor
        ...
    pkg/ <- tidak boleh import /internal, boleh diimport secara gomodule dari repo luar
        pb/
            user.proto
            user.pb.go <-- generated
            user_validation.go
            ...
        http/
            http_api_request.go
            http_api_response.go
            http_api_client.go
        ...
```

## Package Descriptions

1. `/cmd`

* `/cmd/api/main.go` -> run gRPC & echo HTTP on the same runtime, different ports. Or
  choose one as needed. But don't make 2 main programs
  for gRPC & http unless absolutely necessary for scaling / traffic problems
  tall.
* `/cmd/pubsubworker/main.go` -> run pubsub subscription worker. Can subscribe
  multiple subscriptions at the same time (with different goroutines).
  No need to create main programs for different workers unless necessary
  for *scaling / high throughput* problems.

The point is, we avoid *premature-optimization* & make the program more
maintainable in terms of development, deployment & operation.

2. `/internal`

`internal` packages for program logic needed internally by services, and
not intended to be imported by other go programs by `gomodule`.

* `/internal/service` core logic of application gRPC api / JSON webservice
* `/internal/http` implementation of `echo.Echo` http server, router & handler (if
  required). Implement refer to [echo HTTP API
  design](./2_echo_api_design.md)
* `/internal/grpc` gRPC server implementation (if required). Implementation
  refer to [gRPC design](./3_grpc_api_design.md)
* `/internal/worker` logic for GCP Pubsub / other messaging. Implementation
* refer to [pubsub document](./8_pubsub.md)
* `/internal/repository` data access layer logic. Implementation refer to
* [repository design](./4_data_access_repository.md)
* `/internal/apiclient` logic regarding data access via API to third party systems
  third. Can pass http JSON API, SOAP, gRPC etc. Implementation refer to
  [api client design](./5_data_access_api_client.md)
* `/internal/config` contains constant config keys, constructor of dependency
  services.

3. `/pkg`

* `/pkg/pb` contains definition of *.proto file, generated *.pb.go, & validation
  protobuff struct
* `/pkg/http` contains the definition of the request/response body (go/json struct). Can
  also include api-client if needed (remember, communication between internal services
  use gRPC)

### Reference

* [Go Module](https://blog.golang.org/using-go-modules)
* [Go Protobuff](https://developers.google.com/protocol-buffers/docs/gotutorial)
