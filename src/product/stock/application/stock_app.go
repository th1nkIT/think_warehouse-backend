package application

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"math"
	"net/http"
	"think_warehouse/common/constants"
	"think_warehouse/common/httpservice"
	"think_warehouse/src/middleware"
	"think_warehouse/src/product/stock/service"
	"think_warehouse/src/repository/payload"
	sqlc "think_warehouse/src/repository/pgbo_sqlc"
	"think_warehouse/toolkit/config"
	"think_warehouse/toolkit/log"
)

func AddRouteStock(s *httpservice.Service, cfg config.KVStore, e *echo.Echo) {
	svc := service.NewStockService(s.GetDB(), cfg)

	mddw := middleware.NewEnsureToken(s.GetDB(), cfg)

	stock := e.Group("/stock")
	stock.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "stock ok")
	})

	stock.Use(mddw.ValidateToken)
	stock.Use(mddw.ValidateUserBackofficeLogin)

	stock.PUT("", updateStock(svc, cfg))
	stock.POST("/list", listStock(svc, cfg))
	stock.GET("/:guid", getStock(svc, cfg))
}

func updateStock(svc *service.StockService, cfg config.KVStore) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request []payload.UpdateStockPayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		user := ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)

		data, err := svc.UpdateStock(ctx.Request().Context(), request, user)
		if err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to update stock")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUpdateStock(data), nil)
	}
}

func listStock(svc *service.StockService, cfg config.KVStore) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.ListStockPayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		listData, totalData, err := svc.ListStock(ctx.Request().Context(), request.ToEntity())
		if err != nil {
			return err
		}

		// TOTAL PAGE
		totalPage := math.Ceil(float64(totalData) / float64(request.Limit))

		return httpservice.ResponsePagination(ctx, payload.ToPayloadListStock(listData), nil, int(request.Offset), int(request.Limit), int(totalPage), int(totalData))

	}
}

func getStock(svc *service.StockService, cfg config.KVStore) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		data, err := svc.GetStock(ctx.Request().Context(), guid)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadStockSingle(data), nil)
	}
}
