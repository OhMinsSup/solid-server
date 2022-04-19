package sqlstore

import (
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	"solid-server/model"
	"solid-server/utils"
)

// getActiveUserCount 는 N초 전 활성 세션이 있는 사용자 수를 반환합니다.
func (s *SQLStore) getActiveUserCount(db sq.BaseRunner, updatedSecondsAgo int64) (int, error) {
	var updateAt = utils.GetMillis() - utils.SecondsToMillis(updatedSecondsAgo)
	query := s.getQueryBuilder(db).
		Select("count(distinct user_id)").
		From(s.tablePrefix + "sessions").
		Where(sq.Gt{"update_at": updateAt})

	row := query.QueryRow()

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *SQLStore) getSession(db sq.BaseRunner, token string, expireTimeSeconds int64) (*model.Session, error) {
	var updateAt = utils.GetMillis() - utils.SecondsToMillis(expireTimeSeconds)
	query := s.getQueryBuilder(db).
		Select("id", "token", "user_id", "auth_service", "props").
		From(s.tablePrefix + "sessions").
		Where(sq.Eq{"token": token}).
		Where(sq.Gt{"update_at": updateAt})

	row := query.QueryRow()
	session := model.Session{}

	var propsBytes []byte
	err := row.Scan(&session.ID, &session.Token, &session.UserID, &session.AuthService, &propsBytes)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(propsBytes, &session.Props)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SQLStore) createSession(db sq.BaseRunner, session *model.Session) error {
	now := utils.GetMillis()

	propsBytes, err := json.Marshal(session.Props)
	if err != nil {
		return err
	}

	query := s.getQueryBuilder(db).Insert(s.tablePrefix+"sessions").
		Columns("id", "token", "user_id", "props", "create_at", "update_at").
		Values(session.ID, session.Token, session.UserID, session.AuthService, propsBytes, now, now)

	_, err = query.Exec()
	return err
}
