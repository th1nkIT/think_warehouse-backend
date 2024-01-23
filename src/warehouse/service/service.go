package service

import (
	"database/sql"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
)

type WarehouseService struct {
	mainDB *sql.DB
	cfg    config.KVStore
}

func NewWarehouseService(
	mainDB *sql.DB,
	cfg config.KVStore,
) *WarehouseService {
	return &WarehouseService{
		mainDB: mainDB,
		cfg:    cfg,
	}
}
