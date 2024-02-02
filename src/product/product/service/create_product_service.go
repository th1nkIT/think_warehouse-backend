package service

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"think_warehouse/common/httpservice"
	sqlc "think_warehouse/src/repository/pgbo_sqlc"
	"think_warehouse/toolkit/log"
)

func (s *ProductService) CreateProduct(ctx context.Context, request sqlc.InsertProductParams) (product sqlc.Product, userBackoffice sqlc.GetUserBackofficeRow, categoryData sqlc.GetProductCategoryRow, err error) {
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

	userBackoffice, err = q.GetUserBackoffice(ctx, request.CreatedBy)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get user backoffice by guid")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	categoryData, err = q.GetProductCategory(ctx, request.CategoryID)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get product category by guid")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	product, err = q.InsertProduct(ctx, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed insert product")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	if err = tx.Commit(); err != nil {
		log.FromCtx(ctx).Error(err, "error commit")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}

func (s *ProductService) CreateProductWithoutVariant(
	ctx context.Context,
	productParams sqlc.InsertProductParams,
	productPriceParams sqlc.InsertProductPriceParams,
	stockLogParams sqlc.InsertStockLogParams,
	stockParams sqlc.InsertStockParams,
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

	product, err = q.InsertProduct(ctx, productParams)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed insert product")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	productCategory, err = q.GetProductCategory(ctx, productParams.CategoryID)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get product category by guid")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	productPrice, err = q.InsertProductPrice(ctx, productPriceParams)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed insert product price")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	stockLog, err = q.InsertStockLog(ctx, stockLogParams)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed insert product stock log")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	stock, err = q.InsertStock(ctx, stockParams)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed insert product stock")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	if err = tx.Commit(); err != nil {
		log.FromCtx(ctx).Error(err, "error commit")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}

func (s *ProductService) CreateProductWithVariant(
	ctx context.Context,
	productParams sqlc.InsertProductParams,
	productVariantParams []sqlc.InsertProductVariantParams,
	productPriceParams []sqlc.InsertProductPriceParams,
	stockLogParams []sqlc.InsertStockLogParams,
	stockParams []sqlc.InsertStockParams,
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

	product, err = q.InsertProduct(ctx, productParams)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed insert product")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	productCategory, err = q.GetProductCategory(ctx, productParams.CategoryID)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get product category by guid")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	for i := range productVariantParams {
		var (
			productVariant sqlc.ProductVariant
			productPrice   sqlc.ProductPrice
			stockLog       sqlc.StockLog
			stock          sqlc.Stock
		)

		// Check Duplicate Product Variant
		if checkProductVariant, errCheckProductVariant := q.GetProductVariantByNameAndProductID(ctx, sqlc.GetProductVariantByNameAndProductIDParams{
			Name:      productVariantParams[i].Name,
			ProductID: productVariantParams[i].ProductID,
		}); errCheckProductVariant == nil {
			err = errors.Wrapf(err, "duplicate product variant: %s", checkProductVariant.Name)
			log.FromCtx(ctx).Error(err, "failed product variant duplicate:")

			err = errors.WithStack(httpservice.ErrProductVariantDuplicate)

			return
		}

		productVariant, err = q.InsertProductVariant(ctx, productVariantParams[i])
		if err != nil {
			log.FromCtx(ctx).Error(err, "failed insert product variant")
			err = errors.WithStack(httpservice.ErrUnknownSource)

			return
		}

		productPrice, err = q.InsertProductPrice(ctx, productPriceParams[i])
		if err != nil {
			log.FromCtx(ctx).Error(err, "failed insert product price")
			err = errors.WithStack(httpservice.ErrUnknownSource)

			return
		}

		stockLog, err = q.InsertStockLog(ctx, stockLogParams[i])
		if err != nil {
			log.FromCtx(ctx).Error(err, "failed insert product stock log")
			err = errors.WithStack(httpservice.ErrUnknownSource)

			return
		}

		stock, err = q.InsertStock(ctx, stockParams[i])
		if err != nil {
			log.FromCtx(ctx).Error(err, "failed insert product stock")
			err = errors.WithStack(httpservice.ErrUnknownSource)

			return
		}

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
