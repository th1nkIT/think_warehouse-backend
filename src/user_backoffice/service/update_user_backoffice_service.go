package service

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

func (s *UserBackofficeService) UpdateUserBackoffice(ctx context.Context, request sqlc.UpdateUserBackofficeParams) (userBackoffice sqlc.UserBackoffice, role sqlc.UserBackofficeRole, err error) {
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

	userBackoffice, err = q.UpdateUserBackoffice(ctx, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed update user backoffice")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	// get user role data
	role, err = q.GetUserBackofficeRole(ctx, int64(userBackoffice.RoleID))
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get user role data")
		err = errors.WithStack(httpservice.ErrRoleNotFound)

		return
	}

	if err = tx.Commit(); err != nil {
		log.FromCtx(ctx).Error(err, "error commit")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}

func (s *UserBackofficeService) UpdateUserBackofficePassword(ctx context.Context, request sqlc.UpdateUserBackofficePasswordParams) (err error) {
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

	err = q.UpdateUserBackofficePassword(ctx, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed update user backoffice")
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

func (s *UserBackofficeService) UpdateUserBackofficeActive(ctx context.Context, guid string, userData sqlc.GetUserBackofficeRow) (err error) {
	userBackofficeData, err := s.GetUserBackoffice(ctx, guid)
	if err != nil {
		return
	}

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

	isActive := false
	if !userBackofficeData.IsActive.Bool {
		isActive = true
	}

	if err = q.UpdateUserBackofficeIsActive(ctx, sqlc.UpdateUserBackofficeIsActiveParams{
		IsActive: sql.NullBool{
			Bool:  isActive,
			Valid: true,
		},
		UpdatedBy: sql.NullString{
			String: userData.Guid,
			Valid:  true,
		},
		Guid: guid,
	}); err != nil {
		log.FromCtx(ctx).Error(err, "failed update is active status")
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
