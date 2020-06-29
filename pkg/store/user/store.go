package user

import (
	"github.com/labstack/echo"

	"github.com/newzoom/api/pkg/model"
)

// Store - user store interface
type Store interface {
	Get(c echo.Context, id string) (*model.User, error)
	Create(c echo.Context, user *model.User) error
}
