package sqlstore

import (
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"solid-server/model"
	"solid-server/utils"
)

type UserNotFoundError struct {
	id string
}

func (unf UserNotFoundError) Error() string {
	return fmt.Sprintf("user not found (%s)", unf.id)
}

func (s *SQLStore) createUser(db sq.BaseRunner, user *model.User) error {
	now := utils.GetMillis()

	propsBytes, err := json.Marshal(user.Props)
	if err != nil {
		return err
	}

	query := s.getQueryBuilder(db).Insert(s.tablePrefix+"users").
		Columns("id", "username", "email", "password", "mfa_secret", "auth_service", "auth_data", "props", "create_at", "update_at", "delete_at").
		Values(user.ID, user.Username, user.Email, user.Password, user.MfaSecret, user.AuthService, user.AuthData, propsBytes, now, now, 0)

	_, err = query.Exec()
	return err
}
