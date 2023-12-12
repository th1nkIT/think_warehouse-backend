package application

import (
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/constants"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/src/middleware"
	"github.com/wit-id/blueprint-backend-go/src/repository/payload"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/src/user_backoffice_role/service"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

func AddRouteUserBackofficeRole(s *httpservice.Service, cfg config.KVStore, e *echo.Echo) {
	svc := service.NewUserBackofficeRoleService(s.GetDB(), cfg)

	mddw := middleware.NewEnsureToken(s.GetDB(), cfg)

	userBackofficeRole := e.Group(cfg.GetString(constants.ConfigPrefixRoutesBackoffice) + "user-backoffice/role")
	userBackofficeRole.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "user backoffice ok")
	})
	userBackofficeRole.Use(mddw.ValidateToken)
	userBackofficeRole.Use(mddw.ValidateUserBackofficeLogin)

	// example role access page
	userBackofficeRole.POST("", mddw.ValidateRole(createUserBackofficeRole(svc), constants.AccessCreate))

	userBackofficeRole.PUT("/:id", updateUserBackofficeRole(svc))
	userBackofficeRole.DELETE("/:id", deleteUserBackofficeRole(svc))
	userBackofficeRole.POST("/list", listUserBackofficeRole(svc))
	userBackofficeRole.GET("/:id", getUserBackofficeRole(svc))
}

func createUserBackofficeRole(svc *service.UserBackofficeRoleService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.UserBackofficeRolePayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, err := svc.CreateUserBackofficeRole(ctx.Request().Context(), request.ToEntityCreate(ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUserBackofficeRole(data), nil)
	}
}

func updateUserBackofficeRole(svc *service.UserBackofficeRoleService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.ParseInt(ctx.Param("id"), constants.DefaultBaseDecimal, constants.DefaultBitSize)
		if id == 0 || err != nil {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		var request payload.UserBackofficeRolePayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, err := svc.UpdateUserBackofficeRole(ctx.Request().Context(), request.ToEntityUpdate(id, ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUserBackofficeRole(data), nil)
	}
}

func deleteUserBackofficeRole(svc *service.UserBackofficeRoleService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.ParseInt(ctx.Param("id"), constants.DefaultBaseDecimal, constants.DefaultBitSize)
		if id == 0 || err != nil {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		if err := svc.DeleteUserBackoffice(ctx.Request().Context(), id, ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)); err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, nil, nil)
	}
}

func listUserBackofficeRole(svc *service.UserBackofficeRoleService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.ListUserBackofficeRolePayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		listData, totalData, err := svc.ListUserBackofficeRole(ctx.Request().Context(), request.ToEntity())
		if err != nil {
			return err
		}

		// TOTAL PAGE
		totalPage := math.Ceil(float64(totalData) / float64(request.Limit))

		return httpservice.ResponsePagination(ctx, payload.ToPayloadListUserBackofficeRole(listData), nil, int(request.Offset), int(request.Limit), int(totalPage), int(totalData))
	}
}

func getUserBackofficeRole(svc *service.UserBackofficeRoleService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.ParseInt(ctx.Param("id"), constants.DefaultBaseDecimal, constants.DefaultBitSize)
		if id == 0 || err != nil {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		data, err := svc.GetUserBackofficeRole(ctx.Request().Context(), id)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUserBackofficeRole(data), nil)
	}
}
