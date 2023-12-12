# Data Access / Repository

## Storage Redis / MongoDB (manual repository impl)



## Storage SQL (Postgres)

Because we are using Postgres as an RDBMS, access
the data is enough to use [golang library sqlc](https://docs.sqlc.dev/en/latest/).
It doesn't need to be wrapped in a `repository` object. Just save the `*sql.DB` instance in
in the `service` object, and create a query object in each function that
need.

```golang
// internal/service/service.go
type Service struct {
     db *sql.DB // used for data access to postgres
     ...
}

// internal/service/user_service.go
func (s *Service) GetUser(ctx context.Context, ID int) (*User, error) {
     // create query obj
     query := sqlc.New(s.db)
     u, err := query.FindUserByID(ctx, ID)

     ...
}
```

### Reference

* [sqlc tutorial]()