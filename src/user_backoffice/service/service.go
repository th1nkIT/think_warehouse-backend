package service

import (
	"database/sql"

	"github.com/wit-id/blueprint-backend-go/toolkit/config"
)

type UserBackofficeService struct {
	mainDB *sql.DB
	cfg    config.KVStore
}

func NewUserBackofficeService(
	mainDB *sql.DB,
	cfg config.KVStore,
) *UserBackofficeService {
	return &UserBackofficeService{
		mainDB: mainDB,
		cfg:    cfg,
	}
}
