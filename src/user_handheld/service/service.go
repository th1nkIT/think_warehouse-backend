package service

import (
	"database/sql"

	"github.com/wit-id/blueprint-backend-go/toolkit/config"
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
