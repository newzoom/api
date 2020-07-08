package user

import (
	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
	"github.com/phuwn/tools/errors"

	"github.com/newzoom/api/pkg/model"
	"github.com/newzoom/api/pkg/server"
)

func createUser(c echo.Context, u *model.User) error {
	cfg := server.GetServerCfg()
	err := cfg.Store().User.Create(c, u)
	if err != nil {
		return errors.Customize(500, "create user failed", err)
	}
	return nil
}

// VerifyGoogleUser - verify user's google auth code, response google info of user
func VerifyGoogleUser(c echo.Context, code string) (*model.User, error) {
	cfg := server.GetServerCfg()
	token, err := cfg.Service().Google.GetOauth2Token(code)
	if err != nil {
		return nil, err
	}

	return cfg.Service().Google.GetUserGoogleInfo(token)
}

// FirstOrCreate - middleman layer to get the first record that match user's email or create new user if it doesn't exist
func FirstOrCreate(c echo.Context, u *model.User) error {
	cfg := server.GetServerCfg()
	res, err := cfg.Store().User.GetByEmail(c, u.Email)
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return createUser(c, u)
		}
		return errors.Customize(500, "failed to get user with email: "+u.Email, err)
	}

	err = copier.Copy(u, res)
	if err != nil {
		return errors.Customize(500, "failed to copy value", err)
	}
	return nil
}
