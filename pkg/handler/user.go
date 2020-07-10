package handler

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/labstack/echo"
	"github.com/phuwn/tools/errors"

	"github.com/newzoom/api/pkg/handler/user"
	"github.com/newzoom/api/pkg/model"
)

func userRoutes(r *echo.Echo) {
	r.POST("/auth", signIn)
	// g := r.Group("/users")
	// {
	// }
}

// SignInRequest - data form to sign in to auth
type SignInRequest struct {
	Code        string `json:"code"`
	RedirectURL string `json:"redirect_url"`
}

func signIn(c echo.Context) error {
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return errors.Customize(400, "unable to read the request body", err)
	}

	req := &SignInRequest{}
	err = json.Unmarshal(b, req)
	if err != nil {
		return errors.Customize(400, "wrong sign in data form", err)
	}

	u, err := user.VerifyGoogleUser(c, req.Code, req.RedirectURL)
	if err != nil {
		return err
	}

	err = user.FirstOrCreate(c, u)
	if err != nil {
		return err
	}

	jwt, err := model.GenerateJWTToken(&model.TokenInfo{UserID: u.ID}, time.Now().Add(24*time.Hour).Unix())
	if err != nil {
		return err
	}

	u.AccessToken = &jwt
	return JSON(c, 200, u)
}
