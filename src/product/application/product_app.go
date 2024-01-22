package application

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/constants"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/src/middleware"
	"github.com/wit-id/blueprint-backend-go/src/product/service"
	"github.com/wit-id/blueprint-backend-go/src/repository/payload"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
	"math"
	"net/http"
)

func AddRouteProduct(s *httpservice.Service, cfg config.KVStore, e *echo.Echo) {
	svc := service.NewProductService(s.GetDB(), cfg)

	mddw := middleware.NewEnsureToken(s.GetDB(), cfg)

	product := e.Group("/product")
	product.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "product ok")
	})
	product.Use(mddw.ValidateToken)
	product.Use(mddw.ValidateUserBackofficeLogin)

	product.POST("/create", createProduct(svc))
	//product.PUT("/is-active/:guid", updateIsActiveUserBackoffice(svc))
	//product.DELETE("/:guid", deleteUserBackoffice(svc))
	product.POST("/list", listProduct(svc))
	product.GET("/:guid", getProduct(svc))
}

func createProduct(svc *service.ProductService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.RegisterProductPayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, err := svc.CreateProduct(ctx.Request().Context(), request.ToEntity(ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadRegisterProduct(data), nil)
	}
}

func listProduct(svc *service.ProductService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.ListProductPayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		listData, totalData, err := svc.ListProduct(ctx.Request().Context(), request.ToEntity())
		if err != nil {
			return err
		}

		// TOTAL PAGE
		totalPage := math.Ceil(float64(totalData) / float64(request.Limit))

		return httpservice.ResponsePagination(ctx, payload.ToPayloadListProduct(listData), nil, int(request.Offset), int(request.Limit), int(totalPage), int(totalData))
	}
}

func getProduct(svc *service.ProductService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		data, err := svc.GetProduct(ctx.Request().Context(), guid)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadProduct(data), nil)
	}
}
