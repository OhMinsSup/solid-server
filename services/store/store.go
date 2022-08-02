//go:generate mockgen --build_flags=--mod=mod -destination=mockstore/mockstore.go -package mockstore . Store
//go:generate go run ./generators/main.go
package store

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // sqlite driver
	"solid-server/model"
)

type Store interface {
	// database
	Shutdown() error

	// user
	GetUserByID(userID string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) error

	// post
	InsertPost(post *model.Post, userId string) error
	GetSlugDuplicate(slug, userId string) error
	GetPost(postId string) (*model.Post, error)

	// category
	GetFindOrCreate(names []string) ([]*model.Category, error)

	// etc
	DBType() string
}

// ErrNotFound 는 쿼리가 예기치 않게 레코드를 가져오지 않을 때 스토어
// API 에서 반환할 수 있는 오류 유형입니다.
type ErrNotFound struct {
	resource string
}

// NewErrNotFound creates a new ErrNotFound instance.
func NewErrNotFound(resource string) *ErrNotFound {
	return &ErrNotFound{
		resource: resource,
	}
}

func (nf *ErrNotFound) Error() string {
	return fmt.Sprintf("{%s} not found", nf.resource)
}

// IsErrNotFound 는 `err`이 ErrNotFound이거나 래핑된 경우 true를 반환합니다.
func IsErrNotFound(err error) bool {
	if err == nil {
		return false
	}

	// check if this is a store.ErrNotFound
	var nf *ErrNotFound
	if errors.As(err, &nf) {
		return true
	}

	// check if this is a sql.ErrNotFound
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}
	return false
}
