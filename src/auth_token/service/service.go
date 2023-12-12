package service

import (
	"database/sql"

	"github.com/wit-id/blueprint-backend-go/toolkit/config"
)

type AuthTokenService struct {
	mainDB *sql.DB
	cfg    config.KVStore
}

func NewAuthTokenService(
	mainDB *sql.DB,
	cfg config.KVStore,
) *AuthTokenService {
	return &AuthTokenService{
		mainDB: mainDB,
		cfg:    cfg,
	}
}
