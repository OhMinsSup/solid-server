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

func (a *Auth) VerifyAccessToken(tokenString string) (*DecodedToken, error) {
	type MyCustomClaims struct {
		UserID string `json:"user_id"`
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.config.SessionSecretKey), nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return &DecodedToken{
			UserID: claims.UserID,
			Exp:    claims.ExpiresAt.Unix(),
		}, nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, errors.Wrap(err, "That's not even a token")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return nil, errors.Wrap(err, "Token is expired")
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, errors.Wrap(err, "Token has invalid signature")
	} else {
		return nil, errors.Wrap(err, "Token is invalid")
	}
}
