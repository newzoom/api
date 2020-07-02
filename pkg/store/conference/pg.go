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

func (s conferencePGStore) Get(c echo.Context, id string) (*model.Conference, error) {
	tx := db.GetTxFromCtx(c)
	var res model.Conference
	err := tx.Where("id = ?", id).Preload("ConferenceUsers").Preload("ConferenceUsers.User").First(&res).Error
	if err != nil {
		return nil, err
	}

	if len(res.ConferenceUsers) != 0 {
		res.Users = make([]*model.User, len(res.ConferenceUsers))
		for i, v := range res.ConferenceUsers {
			res.Users[i] = v.User
		}
	}

	return &res, nil
}

func (s conferencePGStore) Create(c echo.Context, conference *model.Conference) error {
	tx := db.GetTxFromCtx(c)
	return tx.Create(conference).Error
}
