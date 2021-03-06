package conference

import (
	"github.com/labstack/echo"

	"github.com/newzoom/api/pkg/model"
)

// Store - conference store interface
type Store interface {
	Get(c echo.Context, id string, preload bool) (*model.Conference, error)
	Create(c echo.Context, conference *model.Conference) error
}
