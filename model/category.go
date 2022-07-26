package model

import (
	"encoding/json"
	"io"
)

type Category struct {
	// 카테고리 아이디
	// required: true
	ID string `json:"id"`

	// 카테고리명
	// required: true
	Name string `json:"name"`

	// 생성일시
	// required: true
	CreateAt int64 `json:"create_at,omitempty"`

	// 수정일시
	// required: true
	UpdateAt int64 `json:"update_at,omitempty"`

	// 삭제된 시간, 사용자가 삭제되었음을 나타내도록 설정
	// required: true
	DeleteAt int64 `json:"delete_at"`
}

func CategoryFormJSON(data io.Reader) (*Category, error) {
	var category Category
	if err := json.NewDecoder(data).Decode(&category); err != nil {
		return nil, err
	}
	return &category, nil
}
