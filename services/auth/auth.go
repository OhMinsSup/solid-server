package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"solid-server/services/config"
	"solid-server/services/store"
	"time"
)

type AuthInterface interface {
	CreateAccessToken(userId string) (string, error)
	VerifyAccessToken(token string) (*DecodedToken, error)
}

// Auth authenticates 인증
type Auth struct {
	config      *config.Configuration
	store       store.Store
	permissions interface{}
}

type DecodedToken struct {
	UserID string `json:"user_id"`
	Exp    int64  `json:"exp"`
}

// New authenticates 반환.
func New(config *config.Configuration, store store.Store, permissions interface{}) *Auth {
	return &Auth{config: config, store: store, permissions: permissions}
}

func (a *Auth) CreateAccessToken(userId string) (string, error) {
	type MyCustomClaims struct {
		UserID string `json:"user_id"`
		jwt.RegisteredClaims
	}

	expiry := time.Now().Add(time.Duration(a.config.SessionExpireTime) * time.Second).Unix()

	claims := MyCustomClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiry, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "solid-server",
			Subject:   "oauthToken",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(a.config.SessionSecretKey))
	if err != nil {
		return "", errors.Wrap(err, "failed to create access token")
	}
	return ss, nil
}

func (a *Auth) VerifyAccessToken(token string) (*DecodedToken, error) {
	clamis := jwt.MapClaims{}
	data, err := jwt.ParseWithClaims(token, &clamis, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.config.SessionSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !data.Valid {
		return nil, errors.New("invalid token")
	}

	result := DecodedToken{}
	result.UserID = clamis["user_id"].(string)
	result.Exp = clamis["exp"].(int64)

	return &result, nil
}
