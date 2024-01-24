package service

import (
	"database/sql"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
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
