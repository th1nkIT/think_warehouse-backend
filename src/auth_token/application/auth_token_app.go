package application

import (
	"net/http"

	"think_warehouse/common/jwt"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"think_warehouse/common/httpservice"
	"think_warehouse/src/auth_token/service"
	"think_warehouse/src/middleware"
	"think_warehouse/src/repository/payload"
	"think_warehouse/toolkit/config"
	"think_warehouse/toolkit/log"
)

func AddRouteAuthToken(s *httpservice.Service, cfg config.KVStore, e *echo.Echo) {
	svc := service.NewAuthTokenService(s.GetDB(), cfg)

	mddw := middleware.NewEnsureToken(s.GetDB(), cfg)

	token := e.Group("/token")
	token.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "auth token ok")
	})

	token.POST("/auth", authToken(svc))
	token.GET("/refresh", refreshToken(svc), mddw.ValidateRefreshToken)
}

func authToken(svc *service.AuthTokenService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.AuthTokenPayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, err := svc.AuthToken(ctx.Request().Context(), request)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadAuthToken(data), nil)
	}
}

func refreshToken(svc *service.AuthTokenService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		data, err := svc.RefreshToken(ctx.Request().Context(), ctx.Get("token-data").(jwt.RequestJWTToken))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadAuthToken(data), nil)
	}
}
