package sqlstore

import "solid-server/model"

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

// post

func (s *SQLStore) InsertPost(post model.Post, userId string) error {
	return s.insertPost(s.db, post, userId)
}
