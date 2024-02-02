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

func (s *StockService) UpdateStock(ctx context.Context, request []payload.UpdateStockPayload, users sqlc.GetUserBackofficeRow) (stocks []sqlc.Stock, err error) {
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

	for i := range request {
		stockExisting, errGetStock := q.GetStock(ctx, request[i].GUID)
		if errGetStock != nil {
			log.FromCtx(ctx).Error(errGetStock, "failed get stock")
			errGetStock = errors.WithStack(httpservice.ErrUnknownSource)

			return
		}

		stockUpdate, stockLog := request[i].ToEntityUpdateStock(sqlc.GetUserBackofficeRow{Guid: users.Guid}, stockExisting)

		stock, errStockUpdate := q.UpdateStock(ctx, stockUpdate)
		if errStockUpdate != nil {
			log.FromCtx(ctx).Error(errStockUpdate, "failed update stock")
			errStockUpdate = errors.WithStack(httpservice.ErrUnknownSource)

			return
		}

		_, err = q.InsertStockLog(ctx, stockLog)
		if err != nil {
			log.FromCtx(ctx).Error(err, "failed insert stock log")
			err = errors.WithStack(httpservice.ErrUnknownSource)

			return
		}

		stocks = append(stocks, stock)
	}

	if err = tx.Commit(); err != nil {
		log.FromCtx(ctx).Error(err, "failed commit tx")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}
