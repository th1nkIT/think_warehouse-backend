package application

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/constants"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/src/middleware"
	"github.com/wit-id/blueprint-backend-go/src/repository/payload"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/src/warehouse/service"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
	"math"
	"net/http"
)

func AddRouteWarehouse(s *httpservice.Service, cfg config.KVStore, e *echo.Echo) {
	svc := service.NewWarehouseService(s.GetDB(), cfg)

	mddw := middleware.NewEnsureToken(s.GetDB(), cfg)

	warehouse := e.Group("/warehouse")
	warehouse.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "warehouse ok")
	})

	warehouse.POST("", listWarehouse(svc))
	warehouse.GET("/:guid", getWarehouse(svc))
	warehouse.POST("/create", createWarehouse(svc), mddw.ValidateToken, mddw.ValidateUserBackofficeLogin)
	warehouse.PUT("/:guid", updateWarehouse(svc), mddw.ValidateToken, mddw.ValidateUserBackofficeLogin)
	warehouse.DELETE("/:guid", deleteWarehouse(svc), mddw.ValidateToken, mddw.ValidateUserBackofficeLogin)
	warehouse.GET("/reactive/:guid", reactiveWarehouse(svc), mddw.ValidateToken, mddw.ValidateUserBackofficeLogin)
}

func createWarehouse(svc *service.WarehouseService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.RegisterWarehousePayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, userBackoffice, err := svc.CreateWarehouse(ctx.Request().Context(), request.ToEntity(ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadRegisterWarehouse(data, userBackoffice), nil)
	}
}

func updateWarehouse(svc *service.WarehouseService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		var request payload.UpdateWarehousePayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, userBackoffice, err := svc.UpdateWarehouse(ctx.Request().Context(), request.ToEntity(ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow), guid))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUpdateWarehouse(data, userBackoffice), nil)
	}
}

func deleteWarehouse(svc *service.WarehouseService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		userData := ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)

		err := svc.DeleteWarehouse(ctx.Request().Context(), guid, userData)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, nil, nil)
	}
}

func reactiveWarehouse(svc *service.WarehouseService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		userData := ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)

		err := svc.ReactiveWarehouse(ctx.Request().Context(), guid, userData)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, nil, nil)
	}
}

func listWarehouse(svc *service.WarehouseService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.ListWarehousePayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		listData, totalData, err := svc.ListWarehouse(ctx.Request().Context(), request.ToEntity())
		if err != nil {
			return err
		}

		// TOTAL PAGE
		totalPage := math.Ceil(float64(totalData) / float64(request.Limit))

		return httpservice.ResponsePagination(ctx, payload.ToPayloadListWarehouse(listData), nil, int(request.Offset), int(request.Limit), int(totalPage), int(totalData))
	}
}

func getWarehouse(svc *service.WarehouseService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		data, err := svc.GetWarehouse(ctx.Request().Context(), guid)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadWarehouse(data), nil)
	}
}
