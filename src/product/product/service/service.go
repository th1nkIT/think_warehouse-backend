package service

import (
	"database/sql"
	"think_warehouse/toolkit/config"
)

type ProductService struct {
	mainDB *sql.DB
	cfg    config.KVStore
}

func NewProductService(
	mainDB *sql.DB,
	cfg config.KVStore,
) *ProductService {
	return &ProductService{
		mainDB: mainDB,
		cfg:    cfg,
	}
}
