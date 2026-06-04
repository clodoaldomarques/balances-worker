package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	http *echo.Echo
}

func New() *echo.Echo {
	e := echo.New()
	routes(e)
	return e
}

func routes(e *echo.Echo) {
	e.GET("/", HealthCheck)
}

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
