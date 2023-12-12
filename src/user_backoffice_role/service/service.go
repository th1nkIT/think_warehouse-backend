package service

import (
	"database/sql"

	"github.com/wit-id/blueprint-backend-go/toolkit/config"
)

type UserBackofficeRoleService struct {
	mainDB *sql.DB
	cfg    config.KVStore
}

func NewUserBackofficeRoleService(
	mainDB *sql.DB,
	cfg config.KVStore,
) *UserBackofficeRoleService {
	return &UserBackofficeRoleService{
		mainDB: mainDB,
		cfg:    cfg,
	}
}
