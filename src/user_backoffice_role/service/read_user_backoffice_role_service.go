package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

func (s *UserBackofficeRoleService) ListUserBackofficeRole(ctx context.Context, request sqlc.ListUserBackofficeRoleParams) (listUserBackofficeRole []sqlc.UserBackofficeRole, totalData int64, err error) {
	q := sqlc.New(s.mainDB)

	// Get Total data
	totalData, err = s.getCountUserBackofficeRole(ctx, q, request)
	if err != nil {
		return
	}

	listUserBackofficeRole, err = q.ListUserBackofficeRole(ctx, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get list user backoffice role")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}

func (s *UserBackofficeRoleService) GetUserBackofficeRole(ctx context.Context, id int64) (userBackofficeRole sqlc.UserBackofficeRole, err error) {
	q := sqlc.New(s.mainDB)

	userBackofficeRole, err = q.GetUserBackofficeRole(ctx, id)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get user backoffice role")
		err = errors.WithStack(httpservice.ErrRoleNotFound)

		return
	}

	return
}

func (s *UserBackofficeRoleService) getCountUserBackofficeRole(ctx context.Context, q *sqlc.Queries, request sqlc.ListUserBackofficeRoleParams) (totalData int64, err error) {
	requestQueryParam := sqlc.GetCountListUserBackofficeRoleParams{
		SetName: request.SetName,
		Name:    request.Name,
	}

	totalData, err = q.GetCountListUserBackofficeRole(ctx, requestQueryParam)

	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get total data list user backoffice role")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}
