package service

import (
	"database/sql"
	"think_warehouse/toolkit/config"
)

type StockService struct {
	mainDB *sql.DB
	cfg    config.KVStore
}

func NewStockService(
	mainDB *sql.DB,
	cfg config.KVStore,
) *StockService {
	return &StockService{
		mainDB: mainDB,
		cfg:    cfg,
	}
}
