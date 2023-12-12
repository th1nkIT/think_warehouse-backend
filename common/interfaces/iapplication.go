package interfaces

import "github.com/labstack/echo/v4"

type IApplication interface {
	Initialize(e *echo.Echo)
	Destroy()
}
