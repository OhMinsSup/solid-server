package sqlstore

import (
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"solid-server/model"
)

func (s *SQLStore) getSlugDuplicate(db sq.BaseRunner, slug, userId string) error {
	query := s.getQueryBuilder(db).
		Select("id", "slug", "user_id").
		From(s.tablePrefix + "posts").
		Where(sq.Eq{"user_id": userId, "slug": slug})

	rows, err := query.Query()
	if err != nil {
		s.logger.Error(`getUsersByCondition ERROR`, mlog.Err(err))
		return err
	}
	defer s.CloseRows(rows)

	var processedUrlSlug string

	err = rows.Scan(nil, &processedUrlSlug, nil)
	if err != nil {
		s.logger.Error(
			"getSlugDuplicate insert error",
			mlog.Err(err),
		)
		return err
	}

	if processedUrlSlug == "" {
		return errors.New("processedUrlSlug is Empty")
	}

	return nil
}

func (s *SQLStore) insertPost(db sq.BaseRunner, post *model.Post, userID string) error {
	//now := utils.GetMillis()
	//
	//insertQuery := s.getQueryBuilder(db).Insert("").
	//	Columns(
	//		"id",
	//		"title",
	//		"slug",
	//		"sub_title",
	//		"content",
	//		"publishing_at",
	//		"cover_image",
	//		"disabled_comment",
	//		"create_at",
	//		"update_at",
	//		"delete_at",
	//		"user_id",
	//	)
	//
	//insertQueryValues := map[string]interface{}{
	//	"id":               post.ID,
	//	"title":            post.Title,
	//	"slug":             post.Slug,
	//	"sub_title":        post.SubTitle,
	//	"content":          post.Content,
	//	"publishing_at":    post.PublishingAt,
	//	"cover_image":      post.CoverImage,
	//	"disabled_comment": post.DisabledComment,
	//	"create_at":        now,
	//	"update_at":        now,
	//	"delete_at":        0,
	//	"user_id":          userID,
	//}
	//
	//postQuery := insertQuery.SetMap(insertQueryValues).Into(s.tablePrefix + "posts")
	//if _, err := postQuery.Exec(); err != nil {
	//	s.logger.Error(
	//		"insertPost insert error",
	//		mlog.String("userID", userID),
	//		mlog.Err(err),
	//	)
	//	return err
	//}

	if len(post.Categories) > 0 {
		s.GetFindOrCreate(post.Categories)
	}

	return nil
}
