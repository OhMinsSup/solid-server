package app

import (
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"solid-server/model"
	"solid-server/services/auth"
	"solid-server/utils"

	"github.com/pkg/errors"
)

const (
	DaysPerMonth     = 30
	DaysPerWeek      = 7
	HoursPerDay      = 24
	MinutesPerHour   = 60
	SecondsPerMinute = 60
)

// Login 인증 데이터가 유효한 경우 로그인하여 새 사용자 세션을 만듭니다.
func (a *App) Login(username, email, password, mfaToken string) (string, error) {
	var user *model.User
	if username != "" {
		var err error
		user, err = a.store.GetUserByUsername(username)
		if err != nil {
			// TODO: metrics
			a.logger.Debug("Invalid username for user")
			return "", errors.New("invalid username")
		}
	}

	if user == nil && email != "" {
		var err error
		user, err = a.store.GetUserByEmail(email)
		if err != nil {
			// TODO: metrics
			a.logger.Debug("Invalid email for user")
			return "", errors.New("invalid email")
		}
	}

	if user == nil {
		// TODO: metrics
		return "", errors.New("invalid username or password")
	}

	if !auth.ComparePassword(user.Password, password) {
		// TODO: metrics
		a.logger.Debug("Invalid password for user", mlog.String("userID", user.ID))
		return "", errors.New("invalid password")
	}

	authService := user.AuthService
	if authService == "" {
		authService = "native"
	}

	accessToken, err := a.auth.CreateAccessToken(user.ID)
	if err != nil {
		return "", errors.Wrap(err, "Unable to create access token")
	}

	// TODO: metrics
	// TODO: MFA verification
	return accessToken, nil
}

// RegisterUser 제공된 데이터가 유효한 경우 새 사용자를 생성
func (a *App) RegisterUser(username, email, password string) error {
	var user *model.User
	if username != "" {
		var err error
		user, err = a.store.GetUserByUsername(username)
		if err == nil && user != nil {
			return errors.New("The username already exists")
		}
	}

	if user == nil && email != "" {
		var err error
		user, err = a.store.GetUserByEmail(email)
		if err == nil && user != nil {
			return errors.New("The email already exists")
		}
	}

	passwordSettings := auth.PasswordSettings{
		MinimumLength: 6,
	}

	err := auth.IsPasswordValid(password, passwordSettings)
	if err != nil {
		return errors.Wrap(err, "Invalid password")
	}

	err = a.store.CreateUser(&model.User{
		ID:          utils.NewID(utils.IDTypeUser),
		Username:    username,
		Email:       email,
		Password:    auth.HashPassword(password),
		MfaSecret:   "",
		AuthService: a.config.AuthMode,
		AuthData:    "",
		Props:       map[string]interface{}{},
	})
	if err != nil {
		return errors.Wrap(err, "Unable to create the new user")
	}

	return nil
}
