package sqlstore

import "solid-server/model"

// session

func (s *SQLStore) GetRegisteredUserCount() (int, error) {
	return s.getRegisteredUserCount(s.db)
}

// user

func (s *SQLStore) GetUserByEmail(email string) (*model.User, error) {
	return s.getUserByEmail(s.db, email)
}

func (s *SQLStore) GetUserByID(userID string) (*model.User, error) {
	return s.getUserByID(s.db, userID)

}

func (s *SQLStore) GetUserByUsername(username string) (*model.User, error) {
	return s.getUserByUsername(s.db, username)

}

func (s *SQLStore) CreateUser(user *model.User) error {
	return s.createUser(s.db, user)
}

// team

func (s *SQLStore) GetTeam(ID string) (*model.Team, error) {
	return s.getTeam(s.db, ID)
}

func (s *SQLStore) UpsertTeamSignupToken(team model.Team) error {
	return s.upsertTeamSignupToken(s.db, team)
}
