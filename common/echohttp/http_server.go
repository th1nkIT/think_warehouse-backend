package echohttp

import (
	"context"
	"net/http"
	productHandledApp "think_warehouse/src/product/product/application"
	productCategoryHandledApp "think_warehouse/src/product/product_category/application"

	"think_warehouse/common/constants"
	"think_warehouse/common/httpservice"
	"think_warehouse/toolkit/config"
	"think_warehouse/toolkit/echokit"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	authTokenApp "think_warehouse/src/auth_token/application"
	authorizationBackofficeApp "think_warehouse/src/authorization/backoffice/application"
	authorizationHandheldApp "think_warehouse/src/authorization/handheld/application"

	userBackofficeApp "think_warehouse/src/user_backoffice/application"
	userBackofficeRoleApp "think_warehouse/src/user_backoffice_role/application"

	userHandheldApp "think_warehouse/src/user_handheld/application"

	warehouseHandledApp "think_warehouse/src/warehouse/application"
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

	// Product
	productHandledApp.AddRouteProduct(s, cfg, e)
	// Product Category
	productCategoryHandledApp.AddRouteProductCategory(s, cfg, e)

	// Warehouse
	warehouseHandledApp.AddRouteWarehouse(s, cfg, e)

	// set config routes for role access
	httpservice.SetRouteConfig(ctx, s, cfg, e)

	// run actual server
	echokit.RunServerWithContext(ctx, e, runtimeCfg)
}
