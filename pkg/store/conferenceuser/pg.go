package conferenceuser

import (
	"github.com/labstack/echo"
	"github.com/newzoom/api/pkg/model"
	"github.com/phuwn/tools/db"
)

type conferenceUserPGStore struct{}

// NewStore - create new conference_user store
func NewStore() Store {
	return &conferenceUserPGStore{}
}

func (s conferenceUserPGStore) Create(c echo.Context, cur *model.ConferenceUser) error {
	tx := db.GetTxFromCtx(c)
	return tx.Create(cur).Error
}

func (s conferenceUserPGStore) Get(c echo.Context, userID, conferenceID string) (*model.ConferenceUser, error) {
	tx := db.GetTxFromCtx(c)
	cur := &model.ConferenceUser{}
	return cur, tx.Where("user_id = ? and conference_id = ?", userID, conferenceID).First(cur).Error
}
