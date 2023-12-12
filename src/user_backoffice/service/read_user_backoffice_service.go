package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

func (s *UserBackofficeService) ListUserBackoffice(ctx context.Context, request sqlc.ListUserBackofficeParams) (listUserBackoffice []sqlc.ListUserBackofficeRow, totalData int64, err error) {
	q := sqlc.New(s.mainDB)

	// Get Total data
	totalData, err = s.getCountUserBackoffice(ctx, q, request)
	if err != nil {
		return
	}

	listUserBackoffice, err = q.ListUserBackoffice(ctx, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get list user backoffice")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}

func (s *UserBackofficeService) GetUserBackoffice(ctx context.Context, guid string) (userBackoffice sqlc.GetUserBackofficeRow, err error) {
	q := sqlc.New(s.mainDB)

	userBackoffice, err = q.GetUserBackoffice(ctx, guid)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get user backoffice")
		err = errors.WithStack(httpservice.ErrUserNotFound)

		return
	}

	return
}

func (s *UserBackofficeService) GetUserBackofficeByEmail(ctx context.Context, email string) (userBackoffice sqlc.GetUserBackofficeByEmailRow, err error) {
	q := sqlc.New(s.mainDB)

	userBackoffice, err = q.GetUserBackofficeByEmail(ctx, email)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get user backoffice by email")
		err = errors.WithStack(httpservice.ErrUserNotFound)

		return
	}

	return
}

func (s *UserBackofficeService) getCountUserBackoffice(ctx context.Context, q *sqlc.Queries, request sqlc.ListUserBackofficeParams) (totalData int64, err error) {
	requestQueryParam := sqlc.GetCountListUserBackofficeParams{
		SetName:     request.SetName,
		Name:        request.Name,
		SetPhone:    request.SetPhone,
		Phone:       request.Phone,
		SetEmail:    request.SetEmail,
		Email:       request.Email,
		SetRoleID:   request.SetRoleID,
		RoleID:      request.RoleID,
		SetIsActive: request.SetIsActive,
		IsActive:    request.IsActive,
	}

	totalData, err = q.GetCountListUserBackoffice(ctx, requestQueryParam)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get total data list user backoffice")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}
