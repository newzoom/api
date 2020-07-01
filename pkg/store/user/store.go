package user

import (
	"github.com/labstack/echo"

	"github.com/newzoom/api/pkg/model"
)

// Store - user store interface
type Store interface {
	Create(c echo.Context, user *model.User) error
	GetByEmail(c echo.Context, email string) (*model.User, error)
}
