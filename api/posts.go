package api

import (
	"encoding/json"
	"net/http"
	"solid-server/model"
)

// CreatePostRequest 포스트 등록 요청 body

type CreatePostRequest struct {
	// Title Name
	// required: true
	Title string `json:"title"`

	// Slug url path "ID" string
	// required: false
	Slug string `json:"slug"`

	// SubTitle post sub title info
	// required: false
	SubTitle string `json:"sub_title"`

	// Content post main content data
	// required: true
	Content string `json:"content"`

	// Tags 포스트와 연관된 키워드
	// required: false
	Tags []string `json:"tags"`

	// PublishingAt 포스트 공개일
	// required: false
	PublishingAt int64 `json:"publishing_at"`

	// CoverImage 포스트 배경 이미지
	// required: false
	CoverImage string `json:"cover_image"`

	// DisabledComment 댓글 작성 가능 여부
	// required: false
	DisabledComment bool `json:"disabled_comment"`
}

func (a *API) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	// swagger:operation Post /posts/ create post
	//
	// Post Create
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   description: Post request
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/LoginRequest"
	// responses:
	//   '200':
	//     description: success
	//     schema:
	//       "$ref": "#/definitions/LoginResponse"
	//   '401':
	//     description: invalid login
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   '500':
	//     description: internal error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	ctx := r.Context()
	user := ctx.Value(sessionContextKey).(*model.User)

	userData, err := json.Marshal(user)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "", err)
		return
	}

	jsonBytesResponse(w, http.StatusOK, userData)
}
