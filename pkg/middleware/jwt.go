package middleware

import (
	"strings"

	"github.com/labstack/echo"
	"github.com/phuwn/tools/errors"

	"github.com/newzoom/api/pkg/model"
)

var noAuthPath = []string{
	"/page",
	"/healthz",
	"/auth",
}

func shouldAuth(c echo.Context) bool {
	requestURL := c.Request().RequestURI
	for _, v := range noAuthPath {
		if strings.Contains(requestURL, v) {
			return false
		}
	}
	return true
}

func authenticate(c echo.Context) error {
	token, err := model.GetTokenFromRequest(c)
	if err != nil {
		return err
	}
	uid, err := model.VerifyUserSession(token)
	if err != nil {
		return err
	}
	model.SetUserIDToCtx(c, uid)
	return nil
}

// WithAuth - authentication middleware
func WithAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !shouldAuth(c) {
			return next(c)
		}
		err := authenticate(c)
		if err != nil {
			return errors.Customize(401, "invalid token", err)
		}
		return next(c)
	}
}
