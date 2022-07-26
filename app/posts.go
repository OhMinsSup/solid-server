package app

import (
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/teris-io/shortid"
	"solid-server/model"
	"solid-server/services/types"
	"solid-server/utils"
)

func generateUrlSlug(title string) string {
	urlSlug := slug.Make(title)
	sid, _ := shortid.New(1, shortid.DefaultABC, 20)
	id, _ := sid.Generate()
	result := fmt.Sprintf("%v-%v", urlSlug, id)
	return result
}

func (a *App) CreatePost(body types.CreatePostRequest, userId string) error {
	if len(body.Title) <= 0 {
		return errors.New("the title is empty")
	}

	if len(body.Content) <= 0 {
		return errors.New("the content is empty")
	}
	//
	processedUrlSlug := body.Slug
	//err := a.store.GetSlugDuplicate(processedUrlSlug, userId)
	//if err != nil {
	//	processedUrlSlug = generateUrlSlug(body.Title)
	//}

	err := a.store.InsertPost(&model.Post{
		ID:              utils.NewID(utils.IDTypePost),
		Title:           body.Title,
		SubTitle:        body.SubTitle,
		Slug:            processedUrlSlug,
		Content:         body.Content,
		CoverImage:      body.CoverImage,
		DisabledComment: body.DisabledComment,
		PublishingAt:    body.PublishingAt,
		Categories:      body.Categories,
	}, userId)

	fmt.Println(err)

	return nil
}
