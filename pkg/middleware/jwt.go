package middleware

import (
	"fmt"
	"strings"

	"github.com/labstack/echo"
	"github.com/phuwn/tools/errors"

	"github.com/newzoom/api/pkg/model"
)

var authPath = []string{
	"/conferences",
	"/ws",
}

func shouldAuth(c echo.Context) bool {
	requestURL := c.Request().RequestURI
	for _, v := range authPath {
		if strings.Contains(requestURL, v) {
			return true
		}
	}
	return false
}

func authenticate(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return errors.Customize(401, "missing access_token", fmt.Errorf("empty access_token"))
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
