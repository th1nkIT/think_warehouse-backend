package echohttp

import (
	"context"
	"net/http"

	"github.com/wit-id/blueprint-backend-go/common/constants"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
	"github.com/wit-id/blueprint-backend-go/toolkit/echokit"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	authTokenApp "github.com/wit-id/blueprint-backend-go/src/auth_token/application"
	authorizationBackofficeApp "github.com/wit-id/blueprint-backend-go/src/authorization/backoffice/application"
	authorizationHandheldApp "github.com/wit-id/blueprint-backend-go/src/authorization/handheld/application"

	userBackofficeApp "github.com/wit-id/blueprint-backend-go/src/user_backoffice/application"
	userBackofficeRoleApp "github.com/wit-id/blueprint-backend-go/src/user_backoffice_role/application"

	userHandheldApp "github.com/wit-id/blueprint-backend-go/src/user_handheld/application"
)

func RunEchoHTTPService(ctx context.Context, s *httpservice.Service, cfg config.KVStore) {
	e := echo.New()
	e.HTTPErrorHandler = handleEchoError(cfg)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, constants.DefaultAllowHeaderToken, constants.DefaultAllowHeaderRefreshToken},
	}))

	runtimeCfg := echokit.NewRuntimeConfig(cfg, "restapi")
	runtimeCfg.HealthCheckFunc = s.GetServiceHealth

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	authTokenApp.AddRouteAuthToken(s, cfg, e)
	authorizationBackofficeApp.AddRouteAuthorizationBackoffice(s, cfg, e)
	authorizationHandheldApp.AddRouteAuthorizationHandheld(s, cfg, e)

	userBackofficeRoleApp.AddRouteUserBackofficeRole(s, cfg, e)
	userBackofficeApp.AddRouteUserBackoffice(s, cfg, e)

	userHandheldApp.AddRouteUserHandheld(s, cfg, e)

	// set config routes for role access
	httpservice.SetRouteConfig(ctx, s, cfg, e)

	// run actual server
	echokit.RunServerWithContext(ctx, e, runtimeCfg)
}
