package service

import (
	"context"
	"github.com/pkg/errors"
	"think_warehouse/common/httpservice"
	sqlc "think_warehouse/src/repository/pgbo_sqlc"
	"think_warehouse/toolkit/log"
)

func (s *WarehouseService) ListWarehouse(ctx context.Context, request sqlc.ListWarehouseParams) (listWarehouse []sqlc.ListWarehouseRow, totalData int64, err error) {
	q := sqlc.New(s.mainDB)

	// Get Total data
	totalData, err = s.getCountWarehouse(ctx, q, request)
	if err != nil {
		return
	}

	listWarehouse, err = q.ListWarehouse(ctx, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get list warehouse")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}

func (s *WarehouseService) GetWarehouse(ctx context.Context, guid string) (warehouse sqlc.GetWarehouseRow, err error) {
	q := sqlc.New(s.mainDB)

	warehouse, err = q.GetWarehouse(ctx, guid)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get warehouse")
		err = errors.WithStack(httpservice.ErrWarehouseNotFound)

		return
	}

	return
}

func (s *WarehouseService) getCountWarehouse(ctx context.Context, q *sqlc.Queries, request sqlc.ListWarehouseParams) (totalData int64, err error) {
	requestQueryParams := sqlc.GetCountWarehouseParams{
		SetName: request.SetName,
		Name:    request.Name,
	}

	totalData, err = q.GetCountWarehouse(ctx, requestQueryParams)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get total data list product")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}
