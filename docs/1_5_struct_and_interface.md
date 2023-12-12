# Using Structures and Interfaces

## Struct

* **Naming**, avoid redundant / stuttering *exported struct* names
  with the package name.

    ```go
    package auth

    // ! BAD - don't do this, when used from another package it will be
    // `auth.AuthService{}`
    type AuthService struct { }

    type AuthConfig struct { } // ==> auth.AuthConfig

    // ! Good, when used from another package it becomes `auth.Service{}`
    type Service struct{}

    type Config struct { } // ==> auth.Config
    ```

* **Constructor**, with the format
    * to be exported (called from outside package) `New<Struct Name>(<dependency list>)`
    * for internal (private constructor) `new<struct name>(<dependency list>)`

    ```go
    // NewService returns new service for given postgres database
    func NewService(db *sql.DB) *Service {
        return &Service{
            db: db,
        }
    }

    func newUserProductResponse(p postgres.TbProductRow) UserProductResponse {
        var response UserProductResponse
        response.ID = p.ID
        ...

        return response
    }
    ```

* **Dependency**, for structs that have methods, external
  its dependencies (both interface, config, db, redis, httpClient) are injected via
  constructor.

    ```go
    // BAD
    type Service struct{}

    const key := "my_static_key"

    // * methods are more difficult to unit test, because of their dependencies
    // bound to external package
    // * code is more difficult to read, because you have to look at implementation details
    // to see dependencies2 used
    func (s *Service) CreateResource(ctx context.Context, r *Resource) error {
        // get global db variable
        db := configdb.GetDB()
        db.Query...

        // get global value by config
        configVal := configStore.GetVal(key)
        if configVal == ... {
            ....
        }

        // get redis variables
        rClient := configRedis.GetRedis()
        rCient.Set(..)

        ...

        return nil
    }

    // GOOD
    type Service struct {
        db *sql.DB
        rClient *redis.Client
        configStore config.KVStore
    }

    func NewService(db *sql.DB, r *redis.Client, kv config.KVStore) *Service {
        return &Service{
            db: db,
            rClient: r,
            configStore: kv,
        }
    }

    func (s *Service) CreateResource(ctx context.Context, r *Resource) error {
        s.db.Query...

        // get global value by config
        configVal := s.configStore.GetVal(key)
        if configVal == ... {
            ....
        }

        // get redis variables
        s.rClient.Set(..)

        ...

        return nil
    }
    ```

## Interfaces

* **When to use the interface?**

    * When we need a method, regardless of implementation and
      there is a possibility that **implementation is more than 1**.
    * When you need *mocking* functionality during unit tests.
    * If it turns out that in our program there is only an interface declaration with 1
      implementation, it's better to use struct/pointer-struct with
      usual methods.

* **Naming**, use a name that reflects the behavior of the interface.
    ``` go
    package io

    type Reader interface {
        Read(...)
    }

    // =======================
    promotional packages

    type Processor interface {
        Process(...)
    }
    ```

* Define interface in the package consumer side (packages that use interfaces
  tb)

    ```go
    // BAD, interface and implementation defined in the same package
    package auth

    type Validator interface {
        Validate(ctx context.Context, data Data) error
    }

    type ValidatorImpl struct {
        ...
    }

    func (v *ValidatorImpl) Validate(ctx context.Context, data Data) error {
        ....
    }

    // GOOD,
    package controller

    type validator interface {
        Validate(ctx context.Context, data Data) error
    }

    func createUser(v validator) echo.HandlerFunc {
        return func(ctx echo.Context) error {
            data := ctx.Param(...)

            // use the interface method
            if err := v.Validate(data); err != nil {
                ...
            }

            ...
        }
    }

    //================
    package auth

    type Validator struct { }

    func (v *Validator) Validate(ctx context.Context, data Data) error {
    ```

* If consumed by multiple packages, then define in the general package

    ```go
    package service

    type Validator interface {
        Validate(ctx context.Context, data Data) error
    }

    // =======================
    package grpc

    type Server struct {
        v service.Validator
    }

    func (s *Server) CreateUser(ctx context.Context, req Request) (Response, error) {
        if err := s.v.Validate(...); err != nil {
            ...
        }

        ...
    }

    // ======================
    package controller

    func createUser(v service.Vali