package service

import (
	"database/sql"

	"think_warehouse/toolkit/config"
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
