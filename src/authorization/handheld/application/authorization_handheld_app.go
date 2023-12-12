package application

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/common/jwt"
	"github.com/wit-id/blueprint-backend-go/src/authorization/handheld/service"
	"github.com/wit-id/blueprint-backend-go/src/middleware"
	"github.com/wit-id/blueprint-backend-go/src/repository/payload"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

func AddRouteAuthorizationHandheld(s *httpservice.Service, cfg config.KVStore, e *echo.Echo) {
	svc := service.NewAuthorizationHandheldService(s.GetDB(), cfg)

	mddw := middleware.NewEnsureToken(s.GetDB(), cfg)
	authorizationHandheld := e.Group("/authorization/handheld")
	authorizationHandheld.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "authorization handheld ok")
	})
	authorizationHandheld.Use(mddw.ValidateToken)

	authorizationHandheld.POST("/login", loginHandheld(svc))
	authorizationHandheld.POST("/logout", logoutHandheld(svc))
}

func loginHandheld(svc *service.AuthorizationHandheldService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.AuthorizationHandheldPayload

		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, err := svc.Login(ctx.Request().Context(), request, ctx.Get("token-data").(jwt.RequestJWTToken))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUserHandheld(data), nil)
	}
}

func logoutHandheld(svc *service.AuthorizationHandheldService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := svc.Logout(ctx.Request().Context(), ctx.Get("token-data").(jwt.RequestJWTToken))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, nil, nil)
	}
}
