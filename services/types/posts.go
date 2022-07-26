package types

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
	Categories []string `json:"categories"`

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
