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

	userHandheldService "github.com/wit-id/blueprint-backend-go/src/user_handheld/service"
)

func (s *AuthorizationHandheldService) Login(ctx context.Context, request payload.AuthorizationHandheldPayload, jwtRequest jwt.RequestJWTToken) (userHandheld sqlc.UserHandheld, err error) {
	userhandheldSvc := userHandheldService.NewUserHandheldService(s.mainDB, s.cfg)

	// Check user backoffice by mail
	userHandheld, err = userhandheldSvc.GetUserhandheldByEmail(ctx, request.Email)

	// Check user password valid
	password := utility.GeneratePassword(userHandheld.Salt, request.Password)
	if password != userHandheld.Password {
		err = errors.WithStack(httpservice.ErrPasswordNotMatch)
		return
	}

	// check active user {
	if !userHandheld.IsActive.Bool {
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
	if err = q.RecordUserHandheldLastLogin(ctx, userHandheld.Guid); err != nil {
		log.FromCtx(ctx).Error(err, "failed record last login")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	// Update token auth record
	if err = q.RecordAuthTokenUserLogin(ctx, sqlc.RecordAuthTokenUserLoginParams{
		UserLogin: sql.NullString{
			String: userHandheld.Guid,
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

func (s *AuthorizationHandheldService) Logout(ctx context.Context, request jwt.RequestJWTToken) (err error) {
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
