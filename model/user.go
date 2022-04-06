package model

import (
	"encoding/json"
	"io"
)

const (
	SingleUser   = "single-user"
	GlobalTeamID = "0"
)

// User 유저 모델
// Swagger:model User
type User struct {
	// 유저 아이디
	// required: true
	ID string `json:"id"`

	// 유저 이름
	// required: true
	Username string `json:"username"`

	// 유저 이메일
	// required: true
	Email string `json:"-"`

	// 유저 패스워드
	// swagger:ignore
	Password string `json:"-"`

	// swagger:ignore
	MfaSecret string `json:"-"`

	// swagger:ignore
	AuthService string `json:"-"`

	// swagger:ignore
	AuthData string `json:"-"`

	// 유저 설정 정보
	// required: true
	Props map[string]interface{} `json:"props"`

	// 생성일시
	// required: true
	CreateAt int64 `json:"create_at,omitempty"`

	// 수정일시
	// required: true
	UpdateAt int64 `json:"update_at,omitempty"`

	// 삭제된 시간, 사용자가 삭제되었음을 나타내도록 설정
	// required: true
	DeleteAt int64 `json:"delete_at"`

	// 사용자가 봇인지 아닌지
	// required: true
	IsBot bool `json:"is_bot"`

	// 사용자가 게스트인지 아닌지
	// required: true
	IsGuest bool `json:"is_guest"`
}

// UserPropPatch 사용자 속성 패치
// swagger:model
type UserPropPatch struct {
	// 사용자 정보 업데이트 필드
	// required: false
	UpdatedFields map[string]string `json:"updatedFields"`

	// 사용자 정보 제거 필드
	// required: false
	DeletedFields []string `json:"deletedFields"`
}

// Session 세션 정보
type Session struct {
	ID          string                 `json:"id"`
	Token       string                 `json:"token"`
	UserID      string                 `json:"user_id"`
	AuthService string                 `json:"authService"`
	Props       map[string]interface{} `json:"props"`
	CreateAt    int64                  `json:"create_at,omitempty"`
	UpdateAt    int64                  `json:"update_at,omitempty"`
}

// UserFormJSON 유저 폼 JSON
func UserFormJSON(data io.Reader) (*User, error) {
	var user User
	if err := json.NewDecoder(data).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
