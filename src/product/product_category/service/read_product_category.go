package service

import (
	"context"
	"github.com/pkg/errors"
	"think_warehouse/common/httpservice"
	sqlc "think_warehouse/src/repository/pgbo_sqlc"
	"think_warehouse/toolkit/log"
)

func (s *ProductCategoryService) ListProductCategory(ctx context.Context, request sqlc.ListProductCategoryParams) (listProductCategory []sqlc.ListProductCategoryRow, totalData int64, err error) {
	q := sqlc.New(s.mainDB)

	// Get Total Data
	totalData, err = s.getProductCategoryCount(ctx, q, request)
	if err != nil {
		return
	}

	listProductCategory, err = q.ListProductCategory(ctx, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get list product category")
		err = errors.WithStack(httpservice.ErrUnknownSource)
	}

	return
}

func (s *ProductCategoryService) GetProductCategory(ctx context.Context, guid string) (productCategory sqlc.GetProductCategoryRow, err error) {
	q := sqlc.New(s.mainDB)

	productCategory, err = q.GetProductCategory(ctx, guid)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get product category")
		err = errors.WithStack(httpservice.ErrProductCategoryNotFound)
	}

	return
}

func (s *ProductCategoryService) getProductCategoryCount(ctx context.Context, q *sqlc.Queries, request sqlc.ListProductCategoryParams) (totalData int64, err error) {
	requestQueryParams := sqlc.GetCountListProductCategoryParams{
		SetName:   request.SetName,
		Name:      request.Name,
		SetActive: request.SetActive,
		Active:    request.Active,
	}

	totalData, err = q.GetCountListProductCategory(ctx, requestQueryParams)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get total data list product category")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}
