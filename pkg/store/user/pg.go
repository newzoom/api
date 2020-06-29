package user

import (
	"github.com/labstack/echo"

	"github.com/phuwn/tools/db"

	"github.com/newzoom/api/pkg/model"
)

type userPGStore struct{}

// NewStore - create new user store
func NewStore() Store {
	return &userPGStore{}
}

func (s userPGStore) Get(c echo.Context, id string) (*model.User, error) {
	tx := db.GetTxFromCtx(c)
	var res model.User
	return &res, tx.Where("id = ?", id).First(&res).Error
}

func (s userPGStore) Create(c echo.Context, user *model.User) error {
	tx := db.GetTxFromCtx(c)
	return tx.Create(user).Error
}
