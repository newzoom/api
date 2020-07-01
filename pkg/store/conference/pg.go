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

func (s conferencePGStore) Get(c echo.Context, id int) (*model.Conference, error) {
	tx := db.GetTxFromCtx(c)
	var res model.Conference
	return &res, tx.Where("id = ?", id).First(&res).Error
}

func (s conferencePGStore) Create(c echo.Context, conference *model.Conference) error {
	tx := db.GetTxFromCtx(c)
	return tx.Create(conference).Error
}
