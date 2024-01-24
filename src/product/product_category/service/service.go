package service

import (
	"database/sql"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
)

type ProductCategoryService struct {
	mainDB *sql.DB
	cfg    config.KVStore
}

func NewProductCategoryService(
	mainDB *sql.DB,
	cfg config.KVStore,
) *ProductCategoryService {
	return &ProductCategoryService{
		mainDB: mainDB,
		cfg:    cfg,
	}
}
