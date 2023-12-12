package application

import (
	"math"
	"net/http"

	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/constants"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/src/middleware"
	"github.com/wit-id/blueprint-backend-go/src/repository/payload"
	"github.com/wit-id/blueprint-backend-go/src/user_backoffice/service"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

func AddRouteUserBackoffice(s *httpservice.Service, cfg config.KVStore, e *echo.Echo) {
	svc := service.NewUserBackofficeService(s.GetDB(), cfg)

	mddw := middleware.NewEnsureToken(s.GetDB(), cfg)

	userBackoffice := e.Group(cfg.GetString(constants.ConfigPrefixRoutesBackoffice) + "user-backoffice")
	userBackoffice.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "user backoffice ok")
	})
	userBackoffice.Use(mddw.ValidateToken)
	userBackoffice.Use(mddw.ValidateUserBackofficeLogin)

	userBackoffice.POST("/create", createUserBackoffice(svc, cfg))
	userBackoffice.PUT("/is-active/:guid", updateIsActiveUserBackoffice(svc))
	userBackoffice.DELETE("/:guid", deleteUserBackoffice(svc))
	userBackoffice.POST("/list", listUserBackoffice(svc))
	userBackoffice.GET("/:guid", getUserBackoffice(svc))

	userBackofficeProfile := userBackoffice.Group("/profile")
	userBackofficeProfile.GET("", getUserBackofficeMyProfile(svc))
	userBackofficeProfile.PUT("", updateUserBackofficeMyProfile(svc))
	userBackofficeProfile.PUT("/change-password", updatePasswordUserBackoffice(svc, cfg))
}

func createUserBackoffice(svc *service.UserBackofficeService, cfg config.KVStore) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.RegisterUserBackofficePayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, roleData, err := svc.CreateUserBackoffice(ctx.Request().Context(), request.ToEntity(cfg, ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadRegisterUserBackoffice(cfg, data, roleData), nil)
	}
}

func updateUserBackofficeMyProfile(svc *service.UserBackofficeService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.UpdateUserBackofficePayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, roleData, err := svc.UpdateUserBackoffice(ctx.Request().Context(), request.ToEntity(ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUpdateUserBackoffice(data, roleData), nil)
	}
}

func updatePasswordUserBackoffice(svc *service.UserBackofficeService, cfg config.KVStore) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.UpdateUserBackofficePasswordPayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Get user backoffice data from context
		userBackoffice := ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)

		// Validate request
		if err := request.Validate(userBackoffice); err != nil {
			return err
		}

		err := svc.UpdateUserBackofficePassword(ctx.Request().Context(), request.ToEntity(cfg, userBackoffice))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, nil, nil)
	}
}

func updateIsActiveUserBackoffice(svc *service.UserBackofficeService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		// Get user backoffice data from context
		userBackoffice := ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)

		err := svc.UpdateUserBackofficeActive(ctx.Request().Context(), guid, userBackoffice)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, nil, nil)
	}
}

func listUserBackoffice(svc *service.UserBackofficeService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.ListUserBackofficePayload

		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		listData, totalData, err := svc.ListUserBackoffice(ctx.Request().Context(), request.ToEntity())
		if err != nil {
			return err
		}

		// TOTAL PAGE
		totalPage := math.Ceil(float64(totalData) / float64(request.Limit))

		return httpservice.ResponsePagination(ctx, payload.ToPayloadListUserBackoffice(listData), nil, int(request.Offset), int(request.Limit), int(totalPage), int(totalData))
	}
}

func deleteUserBackoffice(svc *service.UserBackofficeService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		// Get user backoffice data from context
		userBackoffice := ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow)

		err := svc.DeleteUserBackoffice(ctx.Request().Context(), guid, userBackoffice)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, nil, nil)
	}
}

func getUserBackoffice(svc *service.UserBackofficeService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		data, err := svc.GetUserBackoffice(ctx.Request().Context(), guid)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUserBackoffice(data), nil)
	}
}

func getUserBackofficeMyProfile(svc *service.UserBackofficeService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		data, err := svc.GetUserBackoffice(ctx.Request().Context(), ctx.Get(constants.MddwUserBackoffice).(sqlc.GetUserBackofficeRow).Guid)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUserBackoffice(data), nil)
	}
}
