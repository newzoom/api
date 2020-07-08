package conference

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/newzoom/api/pkg/model"
	"github.com/newzoom/api/pkg/server"
	"github.com/phuwn/tools/errors"
)

// Create - middleman layer to create new conference
func Create(c echo.Context, con *model.Conference) error {
	cfg := server.GetServerCfg()
	err := cfg.Store().Conference.Create(c, con)
	if err != nil {
		return errors.Customize(500, "failed to create new conference", err)
	}
	cur := &model.ConferenceUser{UserID: con.HostID, ConferenceID: con.ID}
	err = cfg.Store().ConferenceUser.Create(c, cur)
	if err != nil {
		return errors.Customize(500, "failed to create conference_user record", err)
	}
	return nil
}

// Join - middleman layer to join a conference
func Join(c echo.Context, conID, userID, password string) (*model.Conference, error) {
	cfg := server.GetServerCfg()
	con, err := cfg.Store().Conference.Get(c, conID, true)
	if err != nil {
		return nil, errors.Customize(404, "invalid conference", err)
	}
	if con.Password != "" && password != con.Password {
		return nil, errors.Customize(403, "wrong password", fmt.Errorf("your password was incorect"))
	}

	cur := &model.ConferenceUser{UserID: userID, ConferenceID: conID}
	err = cfg.Store().ConferenceUser.Create(c, cur)
	if err != nil {
		return nil, errors.Customize(500, "failed to create conference_user record", err)
	}

	u, err := cfg.Store().User.Get(c, userID)
	if err != nil {
		return nil, errors.Customize(500, "failed fetch user's data", err)
	}
	con.Users = append(con.Users, u)

	return con, nil
}

// IsMember - middleman layer to check if user has joined the conference
func IsMember(c echo.Context, conID, userID string) (bool, error) {
	cfg := server.GetServerCfg()
	_, err := cfg.Store().ConferenceUser.Get(c, userID, conID)
	if err != nil {
		if !errors.IsRecordNotFound(err) {
			return false, errors.Customize(500, "unknown error occurs", err)
		}
		return false, nil
	}

	return true, nil
}

// Get - middleman layer to get conference data
func Get(c echo.Context, id string, preload bool) (*model.Conference, error) {
	cfg := server.GetServerCfg()
	con, err := cfg.Store().Conference.Get(c, id, preload)
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.Customize(404, "invalid conference", err)
		}
		return nil, errors.Customize(500, "unknown error", err)
	}

	if len(con.ConferenceUsers) != 0 {
		con.Users = make([]*model.User, len(con.ConferenceUsers))
		for i, v := range con.ConferenceUsers {
			con.Users[i] = v.User
		}
	}

	return con, nil
}
