package middleware

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/constants"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/common/jwt"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

type EnsureToken struct {
	mainDB *sql.DB
	config config.KVStore
}

type AccessPage struct {
	Page    string   `json:"page"`
	KeyPage string   `json:"key_page"`
	Access  []string `json:"access"`
}

func NewEnsureToken(db *sql.DB, cfg config.KVStore) *EnsureToken {
	return &EnsureToken{
		mainDB: db,
		config: cfg,
	}
}

func (v *EnsureToken) ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		request := ctx.Request()

		headerDataToken := request.Header.Get(v.config.GetString("header.token-param"))
		if headerDataToken == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, httpservice.MsgHeaderTokenNotFound).SetInternal(errors.Wrap(httpservice.ErrMissingHeaderData, httpservice.MsgHeaderTokenNotFound))
		}

		jwtResponse, err := jwt.ClaimsJwtToken(v.config, headerDataToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, httpservice.MsgHeaderTokenUnauthorized).SetInternal(errors.Wrap(err, httpservice.MsgHeaderTokenUnauthorized))
		}

		// Set data jwt response to ...
		ctx.Set(constants.MddwTokenKey, jwtResponse)

		return next(ctx)
	}
}

func (v *EnsureToken) ValidateRefreshToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		request := ctx.Request()

		headerDataToken := request.Header.Get(v.config.GetString("header.refresh-token-param"))
		if headerDataToken == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, httpservice.MsgHeaderRefreshTokenNotFound).SetInternal(errors.Wrap(httpservice.ErrMissingHeaderData, httpservice.MsgHeaderRefreshTokenNotFound))
		}

		jwtResponse, err := jwt.ClaimsJwtToken(v.config, headerDataToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, httpservice.MsgHeaderRefreshTokenUnauthorized).SetInternal(errors.Wrap(err, httpservice.MsgHeaderRefreshTokenUnauthorized))
		}

		// Set data jwt response to ...
		ctx.Set(constants.MddwTokenKey, jwtResponse)

		return next(ctx)
	}
}

func (v *EnsureToken) ValidateUserBackofficeLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// Get data token session
		tokenAuth := ctx.Get(constants.MddwTokenKey).(jwt.RequestJWTToken)

		q := sqlc.New(v.mainDB)

		tokenData, err := q.GetAuthToken(ctx.Request().Context(), sqlc.GetAuthTokenParams{
			Name:       tokenAuth.AppName,
			DeviceID:   tokenAuth.DeviceID,
			DeviceType: tokenAuth.DeviceType,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, httpservice.MsgHeaderTokenUnauthorized).SetInternal(errors.Wrap(httpservice.ErrUnauthorizedTokenData, httpservice.MsgHeaderTokenUnauthorized))
		}

		if !tokenData.IsLogin {
			return echo.NewHTTPError(http.StatusUnauthorized, httpservice.MsgIsNotLogin).SetInternal(errors.WithMessage(httpservice.ErrUnauthorizedUser, httpservice.MsgIsNotLogin))
		}

		// Get user backoffice
		userBackofficeData, err := q.GetUserBackoffice(ctx.Request().Context(), tokenData.UserLogin.String)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, httpservice.MsgUnauthorizedUser).SetInternal(errors.Wrap(httpservice.ErrUnauthorizedUser, httpservice.MsgUnauthorizedUser))
		}

		// check active user {
		if !userBackofficeData.IsActive.Bool {
			return echo.NewHTTPError(http.StatusUnauthorized, httpservice.MsgUserNotActive).SetInternal(errors.WithMessage(httpservice.ErrUnauthorizedUser, httpservice.MsgUserNotActive))
		}

		// Set data user response to ...
		ctx.Set(constants.MddwUserBackoffice, userBackofficeData)

		return next(ctx)
	}
}

func (v *EnsureToken) ValidateUserHandheldLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// Get data token session
		tokenAuth := ctx.Get(constants.MddwTokenKey).(jwt.RequestJWTToken)

		q := sqlc.New(v.mainDB)

		tokenData, err := q.GetAuthToken(ctx.Request().Context(), sqlc.GetAuthTokenParams{
			Name:       tokenAuth.AppName,
			DeviceID:   tokenAuth.DeviceID,
			DeviceType: tokenAuth.DeviceType,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, httpservice.MsgHeaderTokenUnauthorized).SetInternal(errors.Wrap(httpservice.ErrUnauthorizedTokenData, httpservice.MsgHeaderTokenUnauthorized))
		}

		if !tokenData.IsLogin {
			return echo.NewHTTPError(http.StatusUnauthorized, httpservice.MsgIsNotLogin).SetInternal(errors.WithMessage(httpservice.ErrUnauthorizedUser, httpservice.MsgIsNotLogin))
		}

		// Get user backoffice
		userHandheldData, err := q.GetUserHandheld(ctx.Request().Context(), tokenData.UserLogin.String)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, httpservice.MsgUnauthorizedUser).SetInternal(errors.Wrap(httpservice.ErrUnauthorizedUser, httpservice.MsgUnauthorizedUser))
		}

		// check active user {
		if !userHandheldData.IsActive.Bool {
			return echo.NewHTTPError(http.StatusUnauthorized, httpservice.MsgUserNotActive).SetInternal(errors.WithMessage(httpservice.ErrUnauthorizedUser, httpservice.MsgUserNotActive))
		}

		// Set data user response to ...
		ctx.Set(constants.MddwUserHandheld, userHandheldData)

		return next(ctx)
	}
}

func (v *EnsureToken) ValidateRole(next echo.HandlerFunc, access string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userBoAccess := ctx.Get(constants.MddwKeyRole).(sqlc.GetUserBackofficeRow)

		if !userBoAccess.IsAllAccess.Bool {
			path := ctx.Path()

			var (
				listAccess  []AccessPage
				accessRoute string
				pageAccess  bool
			)

			err := json.Unmarshal([]byte(userBoAccess.RoleAccess.String), &listAccess)
			if err != nil {
				log.FromCtx(ctx.Request().Context()).Info("failed unmarshal user backoffice access")
			}

			for v := range listAccess {
				if strings.Contains(path, listAccess[v].KeyPage) {
					pageAccess = true
					accessRoute = strings.Join(listAccess[v].Access, "|")

					break
				}
			}

			hasAccess := pageAccess && strings.Contains(accessRoute, access)
			if !hasAccess {
				return echo.NewHTTPError(http.StatusUnauthorized, httpservice.MsgUnauthorizedUser).SetInternal(errors.WithMessage(httpservice.ErrUnauthorizedUser, httpservice.MsgUnauthorizedUser))
			}
		}

		return next(ctx)
	}
}
