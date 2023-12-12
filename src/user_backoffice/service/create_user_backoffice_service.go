package service

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

func (s *UserBackofficeService) CreateUserBackoffice(ctx context.Context, request sqlc.InsertUserBackofficeParams) (userBackoffice sqlc.UserBackoffice, role sqlc.UserBackofficeRole, err error) {
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

	// check user role data
	role, err = q.GetUserBackofficeRole(ctx, int64(request.RoleID))
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get user role data")
		err = errors.WithStack(httpservice.ErrRoleNotFound)

		return
	}

	userBackoffice, err = q.InsertUserBackoffice(ctx, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed insert user backoffice")
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
