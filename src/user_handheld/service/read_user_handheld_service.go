package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

func (s *UserHandheldService) ListUserHandheld(ctx context.Context, request sqlc.ListUserHandheldParams) (listUserHandheld []sqlc.UserHandheld, totalData int64, err error) {
	q := sqlc.New(s.mainDB)

	// Get Total Data
	totalData, err = s.getCountUserHandheld(ctx, q, request)
	if err != nil {
		return
	}

	listUserHandheld, err = q.ListUserHandheld(ctx, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get list user handheld")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}

func (s *UserHandheldService) GetUserHandheld(ctx context.Context, guid string) (userHandheld sqlc.UserHandheld, err error) {
	q := sqlc.New(s.mainDB)

	userHandheld, err = q.GetUserHandheld(ctx, guid)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get user handheld")
		err = errors.WithStack(httpservice.ErrUserNotFound)

		return
	}

	return
}

func (s *UserHandheldService) GetUserhandheldByEmail(ctx context.Context, email string) (userHandheld sqlc.UserHandheld, err error) {
	q := sqlc.New(s.mainDB)

	userHandheld, err = q.GetUserHandheldByEmail(ctx, email)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get user handheld by email")
		err = errors.WithStack(httpservice.ErrUserNotFound)

		return
	}

	return
}

func (s *UserHandheldService) getCountUserHandheld(ctx context.Context, q *sqlc.Queries, request sqlc.ListUserHandheldParams) (totalData int64, err error) {
	requestQueryParam := sqlc.GetCountUserHandheldParams{
		SetName:     request.SetName,
		Name:        request.Name,
		SetPhone:    request.SetPhone,
		Phone:       request.Phone,
		SetEmail:    request.SetEmail,
		Email:       request.Email,
		SetGender:   request.SetGender,
		Gender:      request.Gender,
		SetAddress:  request.SetAddress,
		Address:     request.Address,
		SetIsActive: request.SetIsActive,
		IsActive:    request.IsActive,
	}

	totalData, err = q.GetCountUserHandheld(ctx, requestQueryParam)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get total data user handheld")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}
