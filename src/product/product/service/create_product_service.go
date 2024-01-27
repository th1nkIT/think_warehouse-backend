package service

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"think_warehouse/common/httpservice"
	sqlc "think_warehouse/src/repository/pgbo_sqlc"
	"think_warehouse/toolkit/log"
)

func (s *ProductService) CreateProduct(ctx context.Context, request sqlc.InsertProductParams) (product sqlc.Product, userBackoffice sqlc.GetUserBackofficeRow, err error) {
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
