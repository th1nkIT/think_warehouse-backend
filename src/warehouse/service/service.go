package service

import (
	"database/sql"
	"think_warehouse/toolkit/config"
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
