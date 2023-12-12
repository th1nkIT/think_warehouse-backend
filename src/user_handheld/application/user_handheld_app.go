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
	"github.com/wit-id/blueprint-backend-go/src/user_handheld/service"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

func AddRouteUserHandheld(s *httpservice.Service, cfg config.KVStore, e *echo.Echo) {
	svc := service.NewUserHandheldService(s.GetDB(), cfg)

	mddw := middleware.NewEnsureToken(s.GetDB(), cfg)
	userHandheld := e.Group("/user-handheld")
	userHandheld.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "user backoffice ok")
	})
	userHandheld.Use(mddw.ValidateToken)

	userHandheld.POST("/create", createUserHandheld(svc, cfg))

	userHandheldProfile := userHandheld.Group("/profile", mddw.ValidateUserHandheldLogin)
	userHandheldProfile.GET("", getUserHandheldMyProfile(svc))
	userHandheldProfile.PUT("", updateUserHandheldMyProfile(svc))
	userHandheldProfile.PUT("/fcm", updateUserHandheldFCM(svc))
	userHandheldProfile.PUT("/change-password", updateUserHandheldPassword(svc, cfg))

	userHandheldBO := e.Group(cfg.GetString(constants.ConfigPrefixRoutesBackoffice) + "user-handheld")
	userHandheldBO.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "user backoffice ok")
	})
	userHandheldBO.Use(mddw.ValidateToken)
	userHandheldBO.Use(mddw.ValidateUserBackofficeLogin)

	userHandheldBO.PUT("/is-active/:guid", updateUserHandheldIsActive(svc))
	userHandheldBO.DELETE("/:guid", deleteUserHandheld(svc))
	userHandheldBO.POST("/list", listUserHandheld(svc))
	userHandheldBO.GET("/:guid", getUserHandheld(svc))
}

func createUserHandheld(svc *service.UserHandheldService, cfg config.KVStore) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.RegisterUserHandheldPayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, err := svc.CreateUserHandheld(ctx.Request().Context(), request.ToEntity(cfg))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUserHandheld(data), nil)
	}
}

func updateUserHandheldMyProfile(svc *service.UserHandheldService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.UpdateUserHandheldPayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, err := svc.UpdateUserHandheld(ctx.Request().Context(), request.ToEntity(ctx.Get(constants.MddwUserHandheld).(sqlc.UserHandheld)))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUserHandheld(data), nil)
	}
}

func updateUserHandheldFCM(svc *service.UserHandheldService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.UpdateUserHandheldFCMPayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		data, err := svc.UpdateUserHandheldFCMToken(ctx.Request().Context(), request.ToEntity(ctx.Get(constants.MddwUserHandheld).(sqlc.UserHandheld)))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUserHandheld(data), nil)
	}
}

func updateUserHandheldPassword(svc *service.UserHandheldService, cfg config.KVStore) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.UpdateUserHandheldPasswordPayload
		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Get user backoffice data from context
		userBackoffice := ctx.Get(constants.MddwUserHandheld).(sqlc.UserHandheld)

		// Validate request
		if err := request.Validate(userBackoffice); err != nil {
			return err
		}

		err := svc.UpdateUserHandheldPassword(ctx.Request().Context(), request.ToEntity(cfg, userBackoffice))
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, nil, nil)
	}
}

func updateUserHandheldIsActive(svc *service.UserHandheldService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		err := svc.UpdateUserHandheldIsActive(ctx.Request().Context(), guid)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, nil, nil)
	}
}

func deleteUserHandheld(svc *service.UserHandheldService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		err := svc.DeleteUserHandheld(ctx.Request().Context(), guid)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, nil, nil)
	}
}

func listUserHandheld(svc *service.UserHandheldService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request payload.ListUserHandheldPayload

		if err := ctx.Bind(&request); err != nil {
			log.FromCtx(ctx.Request().Context()).Error(err, "failed to parse request")
			return errors.WithStack(httpservice.ErrBadRequest)
		}

		// Validate request
		if err := request.Validate(); err != nil {
			return err
		}

		listData, totalData, err := svc.ListUserHandheld(ctx.Request().Context(), request.ToEntity())
		if err != nil {
			return err
		}

		// TOTAL PAGE
		totalPage := math.Ceil(float64(totalData) / float64(request.Limit))

		return httpservice.ResponsePagination(ctx, payload.ToPayloadListUserHandheld(listData), nil, int(request.Offset), int(request.Limit), int(totalPage), int(totalData))
	}
}

func getUserHandheld(svc *service.UserHandheldService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		guid := ctx.Param("guid")
		if guid == "" {
			return errors.Wrap(httpservice.ErrBadRequest, httpservice.MsgInvalidIDParam)
		}

		data, err := svc.GetUserHandheld(ctx.Request().Context(), guid)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUserHandheld(data), nil)
	}
}

func getUserHandheldMyProfile(svc *service.UserHandheldService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		data, err := svc.GetUserHandheld(ctx.Request().Context(), ctx.Get(constants.MddwUserHandheld).(sqlc.UserHandheld).Guid)
		if err != nil {
			return err
		}

		return httpservice.ResponseData(ctx, payload.ToPayloadUserHandheld(data), nil)
	}
}
