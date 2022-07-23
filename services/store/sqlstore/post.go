package sqlstore

import (
	sq "github.com/Masterminds/squirrel"
	"solid-server/model"
	"solid-server/utils"
)

func (s *SQLStore) insertPost(db sq.BaseRunner, post model.Post, userID string) error {
	now := utils.GetMillis()

	insertQuery := s.getQueryBuilder(db).Insert(s.tablePrefix+"posts").
		Columns(
			"id",
			"title",
			"slug",
			"sub_title",
			"content",
			"tags",
			"publishing_at",
			"cover_image",
			"disabled_comment",
			"create_at",
			"update_at",
			"delete_at",
			"user_id",
		)

	insertQueryValues := map[string]interface{}{
		"id":               post.ID,
		"title":            post.Title,
		"slug":             post.Slug,
		"sub_title":        post.SubTitle,
		"content":          post.Content,
		"tags":             post.Tags,
		"publishing_at":    0,
		"cover_image":      post.CoverImage,
		"disabled_comment": post.DisabledComment,
		"create_at":        now,
		"update_at":        now,
		"delete_at":        0,
		"user_id":          userID,
	}

	query := insertQuery.SetMap(insertQueryValues).Into(s.tablePrefix + "posts")
	if _, err := query.Exec(); err != nil {
		return err
	}
	return nil
}
