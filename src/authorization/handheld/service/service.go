package service

import (
	"database/sql"

	"github.com/wit-id/blueprint-backend-go/toolkit/config"
)

type AuthorizationHandheldService struct {
	mainDB *sql.DB
	cfg    config.KVStore
}

func NewAuthorizationHandheldService(
	mainDB *sql.DB,
	cfg config.KVStore,
) *AuthorizationHandheldService {
	return &AuthorizationHandheldService{
		mainDB: mainDB,
		cfg:    cfg,
	}
}
