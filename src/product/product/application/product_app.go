package application

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"math"
	"net/http"
	"think_warehouse/common/constants"
	"think_warehouse/common/httpservice"
	"think_warehouse/src/middleware"
	"think_warehouse/src/product/product/service"
	"think_warehouse/src/repository/payload"
	sqlc "think_warehouse/src/repository/pgbo_sqlc"
	"think_warehouse/toolkit/config"
	"think_warehouse/toolkit/log"
)

func AddRouteProduct(s *httpservice.Service, cfg config.KVStore, e *echo.Echo) {
	svc := service.NewProductService(s.GetDB(), cfg)

	mddw := middleware.NewEnsureToken(s.GetDB(), cfg)

	product := e.Group("/product")
	product.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "product ok")
	})

	product.POST("", listProduct(svc))
	product.GET("/:guid", getProduct(svc))
	product.POST("/create", createProduct(svc), mddw.ValidateToken, mddw.ValidateUserBackofficeLogin)
	product.PUT("/:guid", updateProduct(svc), mddw.ValidateToken, mddw.ValidateUserBackofficeLogin)
	product.DELETE("/:guid", deleteProduct(svc), mddw.ValidateToken, mddw.ValidateUserBackofficeLogin)
	product.GET("/reactive/:guid", reactiveProduct(svc), mddw.ValidateToken, mddw.ValidateUserBackofficeLogin)
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

		data, userBackoffice, err := svc.CreateProduct(ctx.Request().Context(), request.ToEntity(ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadRegisterProduct(data, userBackoffice), nil)
	}
}

func updateProduct(svc *service.ProductService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		var request payload.UpdateProductPayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, userBackoffice, err := svc.UpdateProduct(ctx.Request().Context(), request.ToEntity(ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow), guid))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUpdateProduct(data, userBackoffice), nil)
	}
}

func deleteProduct(svc *service.ProductService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		userData := ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)

		if err := svc.DeleteProduct(ctx.Request().Context(), guid, userData); err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, nil, nil)
	}
}

func reactiveProduct(svc *service.ProductService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		userData := ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)

		if err := svc.ReactiveProduct(ctx.Request().Context(), guid, userData); err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, nil, nil)
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
