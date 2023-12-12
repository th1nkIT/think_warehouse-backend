package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"

	"github.com/wit-id/blueprint-backend-go/common/jwt"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

func (s *AuthTokenService) ValidateAuthToken(ctx context.Context, token string) (claimsJwt jwt.RequestJWTToken, err error) {
	claimsJwt, err = jwt.ClaimsJwtToken(s.cfg, token)
	if err != nil {
		log.FromCtx(ctx).Error(err, "failed validate claims token")
		err = errors.WithStack(httpservice.ErrInvalidToken)

		return
	}

	return
}
