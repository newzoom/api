package handler

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/newzoom/api/pkg/handler/conference"
	"github.com/newzoom/api/pkg/handler/ws"
	"github.com/newzoom/api/pkg/model"
	"github.com/phuwn/tools/errors"
)

func wsRoutes(r *echo.Echo) {
	r.GET("/ws/:id", serveWs)
}

func serveWs(c echo.Context) error {
	userID := model.GetUserIDFromCtx(c)
	if userID == "" {
		return errors.New("failed to get user id")
	}

	conID := c.Param("id")

	isMember, err := conference.IsMember(c, conID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.Customize(403, "permission denied", fmt.Errorf("invalid request"))
	}

	return ws.Serve(c, conID, userID)
}
