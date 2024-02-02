package application

import (
	"encoding/json"
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
		var request interface{}
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		reqMap := request.(map[string]interface{})

		isVariant := reqMap["is_variant"].(bool)
		if !isVariant {
			return createProductWithoutVariant(svc, ctx, request)
		} else {
			return createProductWithVariant(svc, ctx, request)
		}
	}
}

func createProductWithoutVariant(svc *service.ProductService, ctx echo.Context, rq interface{}) (err error) {
	var request payload.CreateProductWithoutVariantRequest

	byteData, err := json.Marshal(rq)
	if err != nil {
		log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
		return errors.WithStack(httpservice.ErrBadRequest)
	}

	if err = json.Unmarshal(byteData, &request); err != nil {
		log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
		return errors.WithStack(httpservice.ErrBadRequest)
	}

	// Validate request
	if err = request.Validate(); err != nil {
		return
	}

	productParams, productPriceParams, stockLogParams, stockParams := request.ToEntity(ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow))

	product, productCategory, productPrice, stockLog, stock, err := svc.CreateProductWithoutVariant(ctx.Request().Context(), productParams, productPriceParams, stockLogParams, stockParams)
	if err != nil {
		return
	}

	return httpservice.ResponseData(ctx, payload.ToResponsePayloadProductWithoutVariant(product, productCategory, productPrice, stockLog, stock), nil)
}

func createProductWithVariant(svc *service.ProductService, ctx echo.Context, rq interface{}) (err error) {
	var request payload.CreateProductWithVariantRequest

	byteData, err := json.Marshal(rq)
	if err != nil {
		log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
		return errors.WithStack(httpservice.ErrBadRequest)
	}

	if err = json.Unmarshal(byteData, &request); err != nil {
		log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
		return errors.WithStack(httpservice.ErrBadRequest)
	}

	// Validate Request
	if err := request.Validate(); err != nil {
		return err
	}

	productParams, productVariantParams, productPriceParams, stockLogParams, stockParams := request.ToEntity(ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow))

	product, productCategory, productVariant, productPrice, stockLog, stock, err := svc.CreateProductWithVariant(ctx.Request().Context(), productParams, productVariantParams, productPriceParams, stockLogParams, stockParams)
	if err != nil {
		return err
	}

	return httpservice.ResponseData(ctx, payload.ToResponsePayloadProductWithVariant(product, productCategory, productVariant, productPrice, stockLog, stock), nil)
}

func updateProduct(svc *service.ProductService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request interface{}
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		product, err := svc.GetProduct(ctx.Request().Context(), guid)
		if err != nil {
			return err
		}

		isVariant := request.(map[string]interface{})["is_variant"].(bool)

		if !isVariant {
			return updateProductWithoutVariant(svc, ctx, request, product)
		} else {
			return updateProductWithVariant(svc, ctx, request, product)
		}
	}
}

func updateProductWithoutVariant(svc *service.ProductService, ctx echo.Context, rq interface{}, productIn sqlc.GetProductRow) (err error) {
	var request payload.UpdateProductWithoutVariantRequest

	byteData, err := json.Marshal(rq)
	if err != nil {
		log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
		return errors.WithStack(httpservice.ErrBadRequest)
	}

	if err = json.Unmarshal(byteData, &request); err != nil {
		log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
		return errors.WithStack(httpservice.ErrBadRequest)
	}

	// Validate Request
	if err := request.Validate(); err != nil {
		return err
	}

	productParams, productPriceParams, stockLogParams, stockParams := request.ToEntity(ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow), productIn.Guid)

	product, productCategory, productPrice, stockLog, stock, err := svc.UpdateProductWithoutVariant(ctx.Request().Context(), productParams, productPriceParams, stockLogParams, stockParams)
	if err != nil {
		return err
	}

	return httpservice.ResponseData(ctx, payload.ToResponsePayloadProductWithoutVariant(product, productCategory, productPrice, stockLog, stock), nil)
}

func updateProductWithVariant(svc *service.ProductService, ctx echo.Context, rq interface{}, productIn sqlc.GetProductRow) (err error) {
	var request payload.UpdateProductWithVariantRequest

	byteData, err := json.Marshal(rq)
	if err != nil {
		log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
		return errors.WithStack(httpservice.ErrBadRequest)
	}

	if err = json.Unmarshal(byteData, &request); err != nil {
		log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
		return errors.WithStack(httpservice.ErrBadRequest)
	}

	// Validate Request
	if err := request.Validate(); err != nil {
		return err
	}

	productParams, productVariantParams, productPriceParams, stockLogParams, stockParams := request.ToEntity(ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow), productIn.Guid)

	product, productCategory, productVariants, productPrices, stockLogs, stocks, err := svc.UpdateProductWithVariant(ctx.Request().Context(), productParams, productVariantParams, productPriceParams, stockLogParams, stockParams)
	if err != nil {
		return err
	}

	return httpservice.ResponseData(ctx, payload.ToResponsePayloadProductWithVariant(product, productCategory, productVariants, productPrices, stockLogs, stocks), nil)
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
