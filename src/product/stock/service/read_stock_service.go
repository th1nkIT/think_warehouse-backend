package service

import (
	"context"
	"github.com/pkg/errors"
	"think_warehouse/common/httpservice"
	sqlc "think_warehouse/src/repository/pgbo_sqlc"
	"think_warehouse/toolkit/log"
)

func (s *StockService) ListStock(ctx context.Context, request sqlc.ListStockParams) (listStock []sqlc.ListStockRow, totalData int64, err error) {
	q := sqlc.New(s.mainDB)

	// totalData
	totalData, err = s.getCountListStock(ctx, q, request)
	if err != nil {
		return
	}

	// listStock
	listStock, err = q.ListStock(ctx, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get list stock")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}

func (s *StockService) GetStock(ctx context.Context, guid string) (stockData sqlc.GetStockRow, err error) {
	q := sqlc.New(s.mainDB)

	stockData, err = q.GetStock(ctx, guid)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get stock")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}

func (s *StockService) getCountListStock(ctx context.Context, q *sqlc.Queries, request sqlc.ListStockParams) (totalData int64, err error) {
	requestQueryParams := sqlc.GetCountListStockParams{
		SetStockGreater: request.SetStockGreater,
		Stock:           request.Stock,
		SetStockLower:   request.SetStockLower,
		SetProductName:  request.SetProductName,
		ProductName:     request.ProductName,
		SetVariantName:  request.SetVariantName,
		VariantName:     request.VariantName,
	}

	totalData, err = q.GetCountListStock(ctx, requestQueryParams)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get total data list stock")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}
