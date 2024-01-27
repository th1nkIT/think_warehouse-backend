package service

import (
	"database/sql"

	"think_warehouse/toolkit/config"
)

type UserHandheldService struct {
	mainDB *sql.DB
	cfg    config.KVStore
}

func NewUserHandheldService(
	mainDB *sql.DB,
	cfg config.KVStore,
) *UserHandheldService {
	return &UserHandheldService{
		mainDB: mainDB,
		cfg:    cfg,
	}
}
