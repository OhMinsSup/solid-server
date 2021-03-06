package sqlstore

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"solid-server/model"
	"solid-server/utils"
	"strings"
)

func postFields(prefix string) []string {
	fields := []string{
		"id",
		"title",
		"slug",
		"sub_title",
		"content",
		"publishing_at",
		"cover_image",
		"disabled_comment",
		"create_at",
		"update_at",
		"delete_at",
		"user_id",
	}

	if prefix == "" {
		return fields
	}

	prefixedFields := make([]string, len(fields))
	for i, field := range fields {
		if strings.HasPrefix(field, "COALESCE(") {
			prefixedFields[i] = strings.Replace(field, "COALESCE(", "COALESCE("+prefix, 1)
		} else {
			prefixedFields[i] = prefix + field
		}
	}
	return prefixedFields
}

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

func (s *SQLStore) getPost(db sq.BaseRunner, postId string) (*model.Post, error) {
	query := s.getQueryBuilder(db).
		Select(postFields("p.")...).
		From(s.tablePrefix + "posts as p").
		LeftJoin(s.tablePrefix + "users as u on p.user_id=u.id").
		Where(sq.Eq{"p.id": postId})

	rows, err := query.Query()
	if err != nil {
		s.logger.Error(`getPost ERROR`, mlog.Err(err))
		return nil, err
	}
	defer s.CloseRows(rows)

	posts, err := s.postsFromRows(rows)
	if err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, sql.ErrNoRows
	}

	return posts[0], nil
}

func (s *SQLStore) insertPost(db sq.BaseRunner, post *model.Post, userID string) error {
	now := utils.GetMillis()

	insertQuery := s.getQueryBuilder(db).Insert("").
		Columns(
			"id",
			"title",
			"slug",
			"sub_title",
			"content",
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
		"publishing_at":    post.PublishingAt,
		"cover_image":      post.CoverImage,
		"disabled_comment": post.DisabledComment,
		"create_at":        now,
		"update_at":        now,
		"delete_at":        0,
		"user_id":          userID,
	}

	postQuery := insertQuery.SetMap(insertQueryValues).Into(s.tablePrefix + "posts")
	if _, err := postQuery.Exec(); err != nil {
		s.logger.Error(
			"insertPost insert error",
			mlog.String("userID", userID),
			mlog.Err(err),
		)
		return err
	}

	if len(post.Categories) > 0 {
		categories, err := s.GetFindOrCreate(post.Categories)
		if err != nil {
			s.logger.Error(
				"insertPost categories find and created err",
				mlog.String("userID", userID),
				mlog.Err(err),
			)
			return err
		}

		for _, category := range categories {
			err = s.syncPostCategory(db, post.ID, category.ID)
			if err != nil {
				s.logger.Error(
					"insertPost syncPostCategory err",
					mlog.String("userID", userID),
					mlog.Err(err),
				)
				return err
			}
		}
	}

	return nil
}

func (s *SQLStore) syncPostCategory(db sq.BaseRunner, postId, categoryId string) error {
	likedQuery := s.getQueryBuilder(db).
		Select("id", "post_id", "category_id").
		From(s.tablePrefix + "post_categories").
		Where(sq.Eq{"post_id": postId, "category_id": categoryId})

	rows, err := likedQuery.Query()

	if err != nil {
		s.logger.Error(`syncPostCategory ERROR`, mlog.Err(err))
		return err
	}
	defer s.CloseRows(rows)

	list, err := s.postCategoriesFromRows(rows)
	if err != nil {
		s.logger.Error(`syncPostCategory ERROR`, mlog.Err(err))
		return err
	}

	if len(list) == 0 {
		insertBuilder := s.getQueryBuilder(db).
			Insert(s.tablePrefix+"post_categories").
			Columns("id", "post_id", "category_id").
			Values(utils.NewID(utils.IDTypePostCategories), postId, categoryId)

		_, err = insertBuilder.Exec()
		if err != nil {
			s.logger.Error(`syncPostCategory post_categories table insert ERROR`, mlog.Err(err))
			return err
		}
	} else {
		var others []string
		for _, like := range list {
			if like.PostID == postId && like.CategoryId == categoryId {
				// ????????? ??? ??????
			} else {
				// ????????? ?????? ??????
				others = append(others, categoryId)
			}
		}

		if len(others) != 0 {
			insertBuilder := s.getQueryBuilder(db).Insert(s.tablePrefix+"post_categories").
				Columns("id", "post_id", "category_id")
			for _, other := range others {
				insertBuilder = insertBuilder.Values(utils.NewID(utils.IDTypePostCategories), postId, other)
			}

			_, err = insertBuilder.Exec()
			if err != nil {
				s.logger.Error(`getFindOrCreate post_categories table insert ERROR`, mlog.Err(err))
				return err
			}
		}
	}

	return nil
}

func (s *SQLStore) postCategoriesFromRows(rows *sql.Rows) ([]*model.PostCategories, error) {
	results := []*model.PostCategories{}

	for rows.Next() {
		var result model.PostCategories

		err := rows.Scan(
			&result.ID,
			&result.PostID,
			&result.CategoryId,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &result)
	}

	return results, nil
}

func (s *SQLStore) postsFromRows(rows *sql.Rows) ([]*model.Post, error) {
	posts := []*model.Post{}

	for rows.Next() {
		var post model.Post

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Slug,
			&post.SubTitle,
			&post.Content,
			&post.PublishingAt,
			&post.CoverImage,
			&post.DisabledComment,
			&post.CreateAt,
			&post.UpdateAt,
			&post.DeleteAt,
			&post.UserID,
		)
		if err != nil {
			s.logger.Error("postsFromRows scan error", mlog.Err(err))
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}
