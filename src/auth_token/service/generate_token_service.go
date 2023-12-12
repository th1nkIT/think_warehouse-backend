package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/wit-id/blueprint-backend-go/common/utility"

	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/common/jwt"
	"github.com/wit-id/blueprint-backend-go/src/repository/payload"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

func (s *AuthTokenService) AuthToken(ctx context.Context, request payload.AuthTokenPayload) (authToken sqlc.AuthToken, err error) {
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

	// validate app key
	if err = s.validateAppKey(ctx, q, payload.ValidateAppKeyPayload{
		AppName: request.AppName,
		AppKey:  request.AppKey,
	}); err != nil {
		return
	}

	// generate jwt token
	jwtResponse, err := jwt.CreateJWTToken(s.cfg, jwt.RequestJWTToken{
		AppName:    request.AppName,
		DeviceID:   request.DeviceID,
		DeviceType: request.DeviceType,
	})
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed generate token")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	authToken, err = s.recordToken(ctx, q, jwtResponse, false)

	if err = tx.Commit(); err != nil {
		log.FromCtx(ctx).Error(err, "error commit")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}

func (s *AuthTokenService) RefreshToken(ctx context.Context, request jwt.RequestJWTToken) (authToken sqlc.AuthToken, err error) {
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
				log.FromCtx(ctx).Error(err, "error rollback=%s", rollBackErr)
				err = errors.WithStack(httpservice.ErrUnknownSource)

				return
			}
		}
	}()

	fmt.Println(utility.PrettyPrint(request))

	jwtResponse, err := jwt.CreateJWTToken(s.cfg, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed generate token (refresh)")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	authToken, err = s.recordToken(ctx, q, jwtResponse, true)

	if err = tx.Commit(); err != nil {
		log.FromCtx(ctx).Error(err, "error commit")
		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}

func (s *AuthTokenService) validateAppKey(ctx context.Context, q *sqlc.Queries, request payload.ValidateAppKeyPayload) (err error) {
	appKeyData, err := q.GetAppKeyByName(ctx, request.AppName)
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed get app key data by name")
		err = errors.WithStack(httpservice.ErrInvalidAppKey)

		return
	}

	if request.AppKey != appKeyData.Key {
		log.FromCtx(ctx).Info("app key is not match")

		err = errors.WithStack(httpservice.ErrInvalidAppKey)

		return
	}

	return
}

func (s *AuthTokenService) recordToken(ctx context.Context, q *sqlc.Queries, token jwt.ResponseJwtToken, isRefreshToken bool) (authToken sqlc.AuthToken, err error) {
	if !isRefreshToken {
		authToken, err = q.InsertAuthToken(ctx, sqlc.InsertAuthTokenParams{
			Name:                token.AppName,
			DeviceID:            token.DeviceID,
			DeviceType:          token.DeviceType,
			Token:               token.Token,
			TokenExpired:        token.TokenExpired,
			RefreshToken:        token.RefreshToken,
			RefreshTokenExpired: token.RefreshTokenExpired,
		})
	} else {
		// Get record
		authData, errGetRecord := s.ReadAuthToken(ctx, sqlc.GetAuthTokenParams{
			Name:       token.AppName,
			DeviceID:   token.DeviceID,
			DeviceType: token.DeviceType,
		})
		if errGetRecord != nil {
			err = errGetRecord
			return
		}

		authToken, err = q.InsertAuthToken(ctx, sqlc.InsertAuthTokenParams{
			Name:                authData.Name,
			DeviceID:            authData.DeviceID,
			DeviceType:          authData.DeviceType,
			Token:               token.Token,
			TokenExpired:        token.TokenExpired,
			RefreshToken:        token.RefreshToken,
			RefreshTokenExpired: token.RefreshTokenExpired,
			IsLogin:             authData.IsLogin,
			UserLogin:           authData.UserLogin,
		})
	}

	if err != nil {
		log.FromCtx(ctx).Error(err, "failed record token")

		err = errors.WithStack(httpservice.ErrUnknownSource)

		return
	}

	return
}
