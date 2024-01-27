package service

import (
	"context"
	"github.com/pkg/errors"
	"think_warehouse/common/httpservice"
	sqlc "think_warehouse/src/repository/pgbo_sqlc"
	"think_warehouse/toolkit/log"
)

func (s *ProductService) ListProduct(ctx context.Context, request sqlc.ListProductParams) (listProduct []sqlc.ListProductRow, totalData int64, err error) {
	q := sqlc.New(s.mainDB)

	// Get Total data
	totalData, err = s.getCountProduct(ctx, q, request)
	if err != nil {
		return
	}

	listProduct, err = q.ListProduct(ctx, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get list product")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}

func (s *ProductService) GetProduct(ctx context.Context, guid string) (product sqlc.GetProductRow, err error) {
	q := sqlc.New(s.mainDB)

	product, err = q.GetProduct(ctx, guid)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get product")
		err = errors.WithStack(httpservice.ErrProductNotFound)

		return
	}

	return
}

func (s *ProductService) getCountProduct(ctx context.Context, q *sqlc.Queries, request sqlc.ListProductParams) (totalData int64, err error) {
	requestQueryParams := sqlc.GetCountProductListParams{
		SetName: request.SetName,
		Name:    request.Name,
	}

	totalData, err = q.GetCountProductList(ctx, requestQueryParams)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get total data list product")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}
