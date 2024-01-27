package echohttp

import (
	"net/http"

	"think_warehouse/common/httpservice"
	"think_warehouse/toolkit/config"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func handleEchoError(_ config.KVStore) echo.HTTPErrorHandler {
	return func(err error, ctx echo.Context) {
		var echoError *echo.HTTPError

		// if *echo.HTTPError, let echokit middleware handles it
		if errors.As(err, &echoError) {
			return
		}

		statusCode := http.StatusInternalServerError
		// message := "mohon maaf, terjadi kesalahan pada server"
		message := err.Error()

		switch {
		case errors.Is(err, httpservice.ErrBadRequest) || errors.Is(err, httpservice.ErrPasswordNotMatch) || errors.Is(err, httpservice.ErrConfirmPasswordNotMatch):
			statusCode = http.StatusBadRequest
			message = err.Error()
		case errors.Is(err, httpservice.ErrInvalidAppKey) || errors.Is(err, httpservice.ErrInvalidOTP) || errors.Is(err, httpservice.ErrUnauthorizedUser) || errors.Is(err, httpservice.ErrInActiveUser) || errors.Is(err, httpservice.ErrUnauthorizedTokenData):
			statusCode = http.StatusUnauthorized
			message = err.Error()
		case errors.Is(err, httpservice.ErrUserNotFound):
			statusCode = http.StatusNotFound
			message = err.Error()
		case errors.Is(err, httpservice.ErrNoResultData):
			statusCode = http.StatusOK
			message = err.Error()
		}

		_ = ctx.JSON(statusCode, echo.NewHTTPError(statusCode, message))
	}
}
