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

func (s userPGStore) Create(c echo.Context, user *model.User) error {
	tx := db.GetTxFromCtx(c)
	return tx.Create(user).Error
}

func (s userPGStore) GetByEmail(c echo.Context, email string) (*model.User, error) {
	tx := db.GetTxFromCtx(c)
	u := &model.User{}
	return u, tx.Where("email = ?", email).First(u).Error
}
