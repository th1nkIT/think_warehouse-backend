package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"

	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
)

func (s *AuthTokenService) ReadAuthToken(ctx context.Context, request sqlc.GetAuthTokenParams) (authToken sqlc.AuthToken, err error) {
	q := sqlc.New(s.mainDB)

	authToken, err = q.GetAuthToken(ctx, request)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed get auth token")
		err = errors.WithStack(httpservice.ErrInvalidToken)

		return
	}

	return
}
