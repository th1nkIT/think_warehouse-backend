package service

import (
	"context"

	"github.com/pkg/errors"
	"think_warehouse/common/httpservice"

	"think_warehouse/common/jwt"
	"think_warehouse/toolkit/log"
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
