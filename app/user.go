package app

import (
	"github.com/pkg/errors"
	"solid-server/model"
)

// GetUser id로 활성 사용자를 얻습니다.
func (a *App) GetUser(id string) (*model.User, error) {
	if len(id) < 1 {
		return nil, errors.New("no user ID")
	}

	user, err := a.store.GetUserByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "unable to find user")
	}
	return user, nil
}
