package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/phuwn/tools/errors"

	"github.com/newzoom/api/pkg/handler/conference"
	"github.com/newzoom/api/pkg/handler/ws"
	"github.com/newzoom/api/pkg/model"
)

func conferenceRoutes(r *echo.Echo) {
	g := r.Group("/conferences")
	{
		g.POST("", createConference)
		g.POST("/join/:id", joinConference)
		g.GET("/:id", getConference)
	}
}

// NewConferenceRequest - data form to create new conference
type NewConferenceRequest struct {
	Topic       string `json:"topic"`
	Description string `json:"description"`
	Password    string `json:"password"`
}

func createConference(c echo.Context) error {
	hostID := model.GetUserIDFromCtx(c)
	if hostID == "" {
		return errors.New("failed to get host id")
	}

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return errors.Customize(400, "invalid json body", err)
	}

	reqData := &NewConferenceRequest{}
	err = json.Unmarshal(b, reqData)
	if err != nil {
		return errors.Customize(400, "wrong structure json", err)
	}

	var password string
	if reqData.Password != "" {
		password = fmt.Sprintf("%x", md5.Sum([]byte(reqData.Password)))
	}
	con := &model.Conference{
		Topic:       reqData.Topic,
		Description: reqData.Description,
		Password:    password,
		HostID:      hostID,
		IsActive:    true,
	}
	err = conference.Create(c, con)
	if err != nil {
		return err
	}
	ws.NewHub(con.ID)

	return JSON(c, 201, con)
}

// JoinConferenceRequest - data form to join a conference
type JoinConferenceRequest struct {
	Password string `json:"password"`
}

func joinConference(c echo.Context) error {
	userID := model.GetUserIDFromCtx(c)
	if userID == "" {
		return errors.New("failed to get user id")
	}

	conID := c.Param("id")

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return errors.Customize(400, "invalid json body", err)
	}

	reqData := &JoinConferenceRequest{}
	err = json.Unmarshal(b, reqData)
	if err != nil {
		return errors.Customize(400, "wrong structure json", err)
	}

	var password string
	if reqData.Password != "" {
		password = fmt.Sprintf("%x", md5.Sum([]byte(reqData.Password)))
	}

	isMember, err := conference.IsMember(c, conID, userID)
	if err != nil {
		return err
	}
	if isMember {
		con, err := conference.Get(c, conID, true)
		if err != nil {
			return err
		}
		return JSON(c, 200, con)
	}

	con, err := conference.Join(c, conID, userID, password)
	if err != nil {
		return err
	}

	return JSON(c, 200, con)
}

func getConference(c echo.Context) error {
	userID := model.GetUserIDFromCtx(c)
	if userID == "" {
		return errors.New("failed to get user id")
	}
	cid := c.Param("id")
	con, err := conference.Get(c, cid, false)
	if err != nil {
		return err
	}
	isMember, err := conference.IsMember(c, cid, userID)
	if err != nil {
		return err
	}
	if isMember {
		return JSON(c, 200, con)
	}
	if con.Password != "" {
		con.HavePassword = true
	}
	return JSON(c, 200, con)
}
