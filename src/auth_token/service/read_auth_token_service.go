package service

import (
	"context"

	"github.com/pkg/errors"
	"think_warehouse/common/httpservice"
	"think_warehouse/toolkit/log"

	sqlc "think_warehouse/src/repository/pgbo_sqlc"
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
