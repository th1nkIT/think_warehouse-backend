package service

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

func (s *WarehouseService) UpdateWarehouse(ctx context.Context, request sqlc.UpdateWarehouseParams) (warehouse sqlc.Warehouse, userBackoffice sqlc.GetUserBackofficeRow, err error) {
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

	userBackoffice, err = q.GetUserBackoffice(ctx, request.UpdatedBy.String)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get user backoffice data")
		err = errors.WithStack(httpservice.ErrUserNotFound)

		return
	}

	warehouse, err = q.UpdateWarehouse(ctx, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed update warehouse")
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
