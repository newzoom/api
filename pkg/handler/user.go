package handler

import (
	"github.com/labstack/echo"
)

func userRoutes(r *echo.Echo) {
	g := r.Group("/user")
	{
		g.POST("", newUser)
	}
}

func newUser(c echo.Context) error {
	return nil
}
