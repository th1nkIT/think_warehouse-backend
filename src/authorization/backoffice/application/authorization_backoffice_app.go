package application

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"think_warehouse/common/httpservice"
	"think_warehouse/common/jwt"
	"think_warehouse/src/authorization/backoffice/service"
	"think_warehouse/src/middleware"
	"think_warehouse/src/repository/payload"
	"think_warehouse/toolkit/config"
	"think_warehouse/toolkit/log"
)

func AddRouteAuthorizationBackoffice(s *httpservice.Service, cfg config.KVStore, e *echo.Echo) {
	svc := service.NewAuthorizationBackofficeService(s.GetDB(), cfg)

	mddw := middleware.NewEnsureToken(s.GetDB(), cfg)
	authorizationBackoffice := e.Group("/authorization/backoffice")
	authorizationBackoffice.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "authorization backoffice ok")
	})
	authorizationBackoffice.Use(mddw.ValidateToken)

	authorizationBackoffice.POST("/login", loginBackoffice(svc))
	authorizationBackoffice.POST("/logout", logoutBackoffice(svc))
}

func loginBackoffice(svc *service.AuthorizationBackofficeService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.AuthorizationBackofficePayload

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

		return httpservice.ResponseData(ctx, payload.ToPayloadUserBackofficeByMail(data), nil)
	}
}

func logoutBackoffice(svc *service.AuthorizationBackofficeService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := svc.Logout(ctx.Request().Context(), ctx.Get("token-data").(jwt.RequestJWTToken))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, nil, nil)
	}
}
