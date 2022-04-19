//go:generate mockgen --build_flags=--mod=mod -destination=mockstore/mockstore.go -package mockstore . Store
//go:generate go run ./generators/main.go
package store

import "solid-server/model"

type Store interface {
	// session
	GetSession(token string, expireTime int64) (*model.Session, error)

	DBType() string
}
