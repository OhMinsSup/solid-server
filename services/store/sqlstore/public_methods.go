package sqlstore

import "solid-server/model"

// session

func (s *SQLStore) GetActiveUserCount(updatedSecondsAgo int64) (int, error) {
	return s.getActiveUserCount(s.db, updatedSecondsAgo)
}

func (s *SQLStore) GetSession(token string, expireTime int64) (*model.Session, error) {
	return s.getSession(s.db, token, expireTime)
}

func (s *SQLStore) CreateSession(session *model.Session) error {
	return s.createSession(s.db, session)
}

func (s *SQLStore) UpdateSession(session *model.Session) error {
	return s.updateSession(s.db, session)
}

func (s *SQLStore) DeleteSession(sessionID string) error {
	return s.deleteSession(s.db, sessionID)
}

func (s *SQLStore) CleanUpSessions(expireTime int64) error {
	return s.cleanUpSessions(s.db, expireTime)
}

func (s *SQLStore) RefreshSession(session *model.Session) error {
	return s.refreshSession(s.db, session)
}

// user

func (s *SQLStore) CreateUser(user *model.User) error {
	return s.createUser(s.db, user)
}
