package service

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"think_warehouse/common/httpservice"
	"think_warehouse/src/repository/payload"
	sqlc "think_warehouse/src/repository/pgbo_sqlc"
	"think_warehouse/toolkit/log"
)

func (s *ProductService) UpdateProductWithoutVariant(
	ctx context.Context,
	productParams sqlc.UpdateProductParams,
	productPriceParams sqlc.UpdateProductPriceParams,
	stockLogParams sqlc.InsertStockLogParams,
	stockParams sqlc.UpdateStockParams,
) (product sqlc.Product, productCategory sqlc.GetProductCategoryRow, productPrice sqlc.ProductPrice, stockLog sqlc.StockLog, stock sqlc.Stock, err error) {
	tx, err := s.mainDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed begin tx")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	q := sqlc.New(s.mainDB).WithTx(tx)

	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				log.FromCtx(ctx).Error(err, "error rollback", rollBackErr)
				err = errors.WithStack(httpservice.ErrUnknownSource)

				return
			}
		}
	}()

	product, err = q.UpdateProduct(ctx, productParams)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed update product")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	productCategory, err = q.GetProductCategory(ctx, product.CategoryID)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get product category")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	productPrice, err = q.UpdateProductPrice(ctx, productPriceParams)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed update is active product price")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	getParams := payload.GetStockByProductAndVariantParams{
		ProductID:           product.Guid,
		SetProductVariantID: false,
	}

	stockExist, err := q.GetStockByProductAndVariant(ctx, getParams.ToEntity())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.FromCtx(ctx).Error(err, "failed get existing stock")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	if stockExist.Guid != "" {
		stockParams.Guid = stockExist.Guid
		if stockExist.Stock > stockParams.Stock {
			stockLogParams.StockLog = int32(stockExist.Stock) - int32(stockParams.Stock)
			stockLogParams.StockType = sqlc.NullStockTypeEnum{
				Valid:         true,
				StockTypeEnum: sqlc.StockTypeEnumOUT,
			}
		} else {
			stockLogParams.StockLog = int32(stockParams.Stock) - int32(stockExist.Stock)
			stockLogParams.StockType = sqlc.NullStockTypeEnum{
				Valid:         true,
				StockTypeEnum: sqlc.StockTypeEnumIN,
			}
		}

		stock, err = q.UpdateStock(ctx, stockParams)
		if err != nil {
			log.FromCtx(ctx).Error(err, "failed update stock")
			err = errors.WithStack(httpservice.ErrUnknownSource)

			return
		}

		stockLog, err = q.InsertStockLog(ctx, stockLogParams)
		if err != nil {
			log.FromCtx(ctx).Error(err, "failed insert stock log")
			err = errors.WithStack(httpservice.ErrUnknownSource)

			return
		}
	}
	// assumptions didn't need create new stock for update product

	if err = tx.Commit(); err != nil {
		log.FromCtx(ctx).Error(err, "error commit")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}

func (s *ProductService) UpdateProductWithVariant(
	ctx context.Context,
	productParams sqlc.UpdateProductParams,
	productVariantParams []sqlc.UpdateProductVariantParams,
	productPriceParams []sqlc.UpdateProductPriceParams,
	stockLogParams []sqlc.InsertStockLogParams,
	stockParams []sqlc.UpdateStockParams,
) (product sqlc.Product, productCategory sqlc.GetProductCategoryRow, productVariants []sqlc.ProductVariant, productPrices []sqlc.ProductPrice, stockLogs []sqlc.StockLog, stocks []sqlc.Stock, err error) {
	tx, err := s.mainDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed begin tx")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	q := sqlc.New(s.mainDB).WithTx(tx)

	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				log.FromCtx(ctx).Error(err, "error rollback", rollBackErr)
				err = errors.WithStack(httpservice.ErrUnknownSource)

				return
			}
		}
	}()

	product, err = q.UpdateProduct(ctx, productParams)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed insert product")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	productCategory, err = q.GetProductCategory(ctx, product.CategoryID)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get product category")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	for i := range productVariantParams {
		var (
			productVariant sqlc.ProductVariant
			productPrice   sqlc.ProductPrice
			stockLog       sqlc.StockLog
			stock          sqlc.Stock
			stockExist     sqlc.Stock
		)

		// Check duplicate product variant
		if checkProductVariant, errCheckProductVariant := q.GetProductVariantByNameAndProductID(ctx, sqlc.GetProductVariantByNameAndProductIDParams{
			Name:      productVariantParams[i].Name,
			ProductID: productVariantParams[i].ProductID,
		}); errCheckProductVariant == nil {
			err = errors.Wrapf(err, "duplicate product variant: %s", checkProductVariant.Name)
			log.FromCtx(ctx).Error(err, "failed product variant duplicate:")

			err = errors.WithStack(httpservice.ErrProductVariantDuplicate)

			return
		}

		productVariant, err = q.UpdateProductVariant(ctx, productVariantParams[i])
		if err != nil {
			log.FromCtx(ctx).Error(err, "failed update product variant")
			err = errors.WithStack(httpservice.ErrUnknownSource)

			return
		}

		productPrice, err = q.UpdateProductPrice(ctx, productPriceParams[i])
		if err != nil {
			log.FromCtx(ctx).Error(err, "failed update is active product price")
			err = errors.WithStack(httpservice.ErrUnknownSource)

			return
		}

		getParams := payload.GetStockByProductAndVariantParams{
			ProductID:           product.Guid,
			SetProductVariantID: true,
			ProductVariantID:    productVariant.Guid,
		}

		stockExist, err = q.GetStockByProductAndVariant(ctx, getParams.ToEntity())
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.FromCtx(ctx).Error(err, "failed get existing stock")
			err = errors.WithStack(httpservice.ErrUnknownSource)

			return
		}

		if stockExist.Guid != "" {
			stockParams[i].Guid = stockExist.Guid
			if stockExist.Stock > stockParams[i].Stock {
				stockLogParams[i].StockLog = int32(stockExist.Stock) - int32(stockParams[i].Stock)
				stockLogParams[i].StockType = sqlc.NullStockTypeEnum{
					Valid:         true,
					StockTypeEnum: sqlc.StockTypeEnumOUT,
				}
			} else {
				stockLogParams[i].StockLog = int32(stockParams[i].Stock) - int32(stockExist.Stock)
				stockLogParams[i].StockType = sqlc.NullStockTypeEnum{
					Valid:         true,
					StockTypeEnum: sqlc.StockTypeEnumIN,
				}
			}

			stockLogParams[i].ProductID = sql.NullString{String: product.Guid, Valid: true}
			stockLogParams[i].ProductVariantID = sql.NullString{String: productVariant.Guid, Valid: true}

			stock, err = q.UpdateStock(ctx, stockParams[i])
			if err != nil {
				log.FromCtx(ctx).Error(err, "failed update stock")
				err = errors.WithStack(httpservice.ErrUnknownSource)

				return
			}

			stockLog, err = q.InsertStockLog(ctx, stockLogParams[i])
			if err != nil {
				log.FromCtx(ctx).Error(err, "failed insert product")
				err = errors.WithStack(httpservice.ErrUnknownSource)

				return
			}
		}
		// assumptions didn't need create new stock for update product

		productVariants = append(productVariants, productVariant)
		productPrices = append(productPrices, productPrice)
		stockLogs = append(stockLogs, stockLog)
		stocks = append(stocks, stock)
	}

	if err = tx.Commit(); err != nil {
		log.FromCtx(ctx).Error(err, "error commit")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}
