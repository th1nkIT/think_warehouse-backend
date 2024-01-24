package application

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/constants"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/src/middleware"
	"github.com/wit-id/blueprint-backend-go/src/product/product_category/service"
	"github.com/wit-id/blueprint-backend-go/src/repository/payload"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
	"math"
	"net/http"
)

func AddRouteProductCategory(s *httpservice.Service, cfg config.KVStore, e *echo.Echo) {
	svc := service.NewProductCategoryService(s.GetDB(), cfg)

	mddw := middleware.NewEnsureToken(s.GetDB(), cfg)

	product := e.Group("/product-category")
	product.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "product ok")
	})

	product.POST("", listProductCategory(svc))
	product.GET("/:guid", getProductCategory(svc))
	product.POST("/create", createProductCategory(svc), mddw.ValidateToken, mddw.ValidateUserBackofficeLogin)
	//product.PUT("/:guid", updateProductCategory(svc), mddw.ValidateToken, mddw.ValidateUserBackofficeLogin)
	//product.DELETE("/:guid", deleteProductCategory(svc), mddw.ValidateToken, mddw.ValidateUserBackofficeLogin)
	//product.GET("/reactive/:guid", reactiveProductCategory(svc), mddw.ValidateToken, mddw.ValidateUserBackofficeLogin)
}

func createProductCategory(svc *service.ProductCategoryService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.RegisterProductCategoryPayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, userBackoffice, err := svc.CreateProductCategory(ctx.Request().Context(), request.ToEntity(ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadRegisterProductCategory(data, userBackoffice), nil)
	}
}

func getProductCategory(svc *service.ProductCategoryService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		data, err := svc.GetProductCategory(ctx.Request().Context(), guid)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadProductCategory(data), nil)
	}
}

func listProductCategory(svc *service.ProductCategoryService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.ListProductCategoryPayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		listData, totalData, err := svc.ListProductCategory(ctx.Request().Context(), request.ToEntity())
		if err != nil {
			return err
		}

		// Total Page
		totalPage := math.Ceil(float64(totalData) / float64(request.Limit))

		return httpservice.ResponsePagination(ctx, payload.ToPayloadListProductCategory(listData), nil, int(request.Offset), int(request.Limit), int(totalPage), int(totalData))
	}
}
