package service

import (
	"database/sql"

	"think_warehouse/toolkit/config"
)

type AuthorizationBackofficeService struct {
	mainDB *sql.DB
	cfg    config.KVStore
}

func NewAuthorizationBackofficeService(
	mainDB *sql.DB,
	cfg config.KVStore,
) *AuthorizationBackofficeService {
	return &AuthorizationBackofficeService{
		mainDB: mainDB,
		cfg:    cfg,
	}
}
