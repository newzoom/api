package conferenceuser

import (
	"github.com/labstack/echo"

	"github.com/newzoom/api/pkg/model"
)

// Store - conference_user store interface
type Store interface {
	Create(c echo.Context, cur *model.ConferenceUser) error
	Get(c echo.Context, userID, conferenceID string) (*model.ConferenceUser, error)
}
