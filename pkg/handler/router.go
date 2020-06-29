package handler

import (
	"github.com/labstack/echo"

	"github.com/phuwn/tools/db"
	"github.com/phuwn/tools/handler"
	mw "github.com/phuwn/tools/middleware"
)

// JSON - shorcut for handler.JSON function
var JSON = handler.JSON

// Router - handling routes for incoming request
func Router() *echo.Echo {
	r := echo.New()
	r.HTTPErrorHandler = handler.JSONError
	r.Pre(mw.RemoveTrailingSlash)
	{
		r.Use(mw.CorsConfig())
		r.Use(mw.AddTransaction)
	}

	r.GET("/healthz", healthz)
	{
		userRoutes(r)
	}

	return r
}

func healthz(c echo.Context) error {
	err := db.Healthz()
	if err != nil {
		return err
	}
	return JSON(c, 200, "ok")
}
