package model

import (
	"encoding/json"
	"io"
)

type Post struct {
	// 포스트 아아디
	// required: true
	ID string `json:"id"`

	// 포스트 제목
	// required: true
	Title string `json:"title"`

	// 포스트 Slug
	// required: true
	Slug string `json:"slug"`

	// 포스트 sub title
	// required: false
	SubTitle string `json:"sub_title"`

	// 포스트 내용
	// required: true
	Content string `json:"content"`

	// 포스트에 연관된 태그
	// required: false
	// 임시 타입으로 설정
	Tags []string `json:"tags"`

	// 포스트 공개일
	// required: false
	PublishingAt int64 `json:"publishing_at"`

	// 포스트 배경 이미지
	// required: false
	CoverImage string `json:"cover_image"`

	// 댓글 작성 가능 여부
	// required: false
	DisabledComment bool `json:"disabled_comment"`

	// 생성일시
	// required: true
	CreateAt int64 `json:"create_at,omitempty"`

	// 수정일시
	// required: true
	UpdateAt int64 `json:"update_at,omitempty"`

	// 삭제된 시간, 사용자가 삭제되었음을 나타내도록 설정
	// required: true
	DeleteAt int64 `json:"delete_at"`

	// 포스트를 작성한 유저 아이디
	// required: true
	UserID string `json:"user_id"`
}

// PostFormJSON 포스트 폼 JSON
func PostFormJSON(data io.Reader) (*Post, error) {
	var post Post
	if err := json.NewDecoder(data).Decode(&post); err != nil {
		return nil, err
	}
	return &post, nil
}
