package model

import (
	"encoding/json"
	"io"
)

type PostCategories struct {
	// 카테고리 아이디
	// required: true
	ID string `json:"id"`

	// 카테고리명
	// required: true
	Name string `json:"name"`

	// 등록되는 포스트 아이디
	// required: true
	PostID string `json:"post_id"`

	// 등록되는 카테고리 아이디
	// required: true
	CategoryId string `json:"category_id"`
}

func PostCategoriesFormJSON(data io.Reader) (*PostCategories, error) {
	var pc PostCategories
	if err := json.NewDecoder(data).Decode(&pc); err != nil {
		return nil, err
	}
	return &pc, nil
}
