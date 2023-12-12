package service

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/common/jwt"
	"github.com/wit-id/blueprint-backend-go/common/utility"
	"github.com/wit-id/blueprint-backend-go/src/repository/payload"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"

	userBackofficeService "github.com/wit-id/blueprint-backend-go/src/user_backoffice/service"
)

func (s *AuthorizationBackofficeService) Login(ctx context.Context, request payload.AuthorizationBackofficePayload, jwtRequest jwt.RequestJWTToken) (userBackoffice sqlc.GetUserBackofficeByEmailRow, err error) {
	userBackofficeSvc := userBackofficeService.NewUserBackofficeService(s.mainDB, s.cfg)

	// Check user backoffice by mail
	userBackoffice, err = userBackofficeSvc.GetUserBackofficeByEmail(ctx, request.Email)
	if err != nil {
		return
	}

	// Check user password valid
	password := utility.GeneratePassword(userBackoffice.Salt, request.Password)
	if password != userBackoffice.Password {
		err = errors.WithStack(httpservice.ErrPasswordNotMatch)
		return
	}

	// check active user {
	if !userBackoffice.IsActive.Bool {
		err = errors.WithStack(httpservice.ErrInActiveUser)
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

	// Update Last login user backoffice
	if err = q.RecordUserBackofficeLastLogin(ctx, userBackoffice.Guid); err != nil {
		log.FromCtx(ctx).Error(err, "failed record last login")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	// Update token auth record
	if err = q.RecordAuthTokenUserLogin(ctx, sqlc.RecordAuthTokenUserLoginParams{
		UserLogin: sql.NullString{
			String: userBackoffice.Guid,
			Valid:  true,
		},
		Name:       jwtRequest.AppName,
		DeviceID:   jwtRequest.DeviceID,
		DeviceType: jwtRequest.DeviceType,
	}); err != nil {
		log.FromCtx(ctx).Error(err, "failed update token auth login user")
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

func (s *AuthorizationBackofficeService) Logout(ctx context.Context, request jwt.RequestJWTToken) (err error) {
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

	if err = q.ClearAuthTokenUserLogin(ctx, sqlc.ClearAuthTokenUserLoginParams{
		Name:       request.AppName,
		DeviceID:   request.DeviceID,
		DeviceType: request.DeviceType,
	}); err != nil {
		log.FromCtx(ctx).Error(err, "failed clear auth user login")
		err = errors.WithStack(err)

		return
	}

	if err = tx.Commit(); err != nil {
		log.FromCtx(ctx).Error(err, "error commit")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}
