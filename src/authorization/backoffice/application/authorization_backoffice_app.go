package application

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/common/jwt"
	"github.com/wit-id/blueprint-backend-go/src/authorization/backoffice/service"
	"github.com/wit-id/blueprint-backend-go/src/middleware"
	"github.com/wit-id/blueprint-backend-go/src/repository/payload"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
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
