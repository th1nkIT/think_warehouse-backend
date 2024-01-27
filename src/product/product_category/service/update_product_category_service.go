package service

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

func (s *ProductCategoryService) UpdateProductCategory(ctx context.Context, payload sqlc.UpdateProductCategoryParams) (productCategory sqlc.ProductCategory, userBackoffice sqlc.GetUserBackofficeRow, err error) {
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

	// Check if product category already deleted
	productCategoryData, err := q.GetProductCategory(ctx, payload.Guid)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get product category data")
		err = errors.WithStack(httpservice.ErrProductCategoryNotFound)

		return
	}

	if productCategoryData.DeletedAt.Valid {
		log.FromCtx(ctx).Error(err, "product category already deleted")
		err = errors.WithStack(httpservice.ErrProductCategoryIsInactive)

		return
	}

	userBackoffice, err = q.GetUserBackoffice(ctx, payload.UpdatedBy.String)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get user backoffice data")
		err = errors.WithStack(httpservice.ErrUserNotFound)

		return
	}

	productCategory, err = q.UpdateProductCategory(ctx, payload)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed update product category")
		err = errors.WithStack(httpservice.ErrProductCategoryFailed)

		return
	}

	if err = tx.Commit(); err != nil {
		log.FromCtx(ctx).Error(err, "error commit")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}
