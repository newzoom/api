package conference

import (
	"github.com/labstack/echo"

	"github.com/phuwn/tools/db"

	"github.com/newzoom/api/pkg/model"
)

type conferencePGStore struct{}

// NewStore - create new conference store
func NewStore() Store {
	return &conferencePGStore{}
}

func (s conferencePGStore) Get(c echo.Context, id string, preload bool) (*model.Conference, error) {
	tx := db.GetTxFromCtx(c)
	var res model.Conference
	tx = tx.Where("id = ?", id)
	if preload {
		tx = tx.Preload("ConferenceUsers").Preload("ConferenceUsers.User")
	}
	return &res, tx.First(&res).Error
}

func (s conferencePGStore) Create(c echo.Context, conference *model.Conference) error {
	tx := db.GetTxFromCtx(c)
	return tx.Create(conference).Error
}
