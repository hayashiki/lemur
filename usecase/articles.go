package usecase

import (
	"fmt"
	"github.com/hayashiki/lemur/docbase"
	"github.com/hayashiki/lemur/entity"
	"github.com/hayashiki/lemur/event"
	"github.com/hayashiki/lemur/event/eventtask"
	"github.com/hayashiki/lemur/logger"
	"log"
	"time"
)

type Articles struct {
	log       logger.Logger
	docBase   docbase.Client
	taskQueue event.TaskQueue
	repo      entity.ArticleRepository
}

type EnqueueArticlesInputParams struct {
	Article entity.Article
}

func NewArticles(
	log logger.Logger,
	docBase docbase.Client,
	taskQueue event.TaskQueue,
	repo entity.ArticleRepository,
) *Articles {
	return &Articles{
		log,
		docBase,
		taskQueue,
		repo,
	}
}

func (a *Articles) Do() error {

	posts, err := a.docBase.PostList(createQuery())

	if err != nil {
		return err
	}

	if len(posts) == 0 {
		return nil
	}

	for _, p := range posts {

		var attachments []*entity.Attachment
		for _, at := range p.Attachments {
			img := entity.NewAttachment(at.ID, at.Name, at.URL)
			attachments = append(attachments, img)
		}

		article := entity.NewArticle(int64(p.ID), p.Title, p.Body, p.CreatedAt, attachments)
		// TODO PutMultiにする
		_, err := a.repo.Put(article)

		if err != nil {
			log.Printf("error is %v", err)
			a.log.Error("Article.Do failed artcleID: %d, err: %v", article.ID, err)
		}

		err = a.taskQueue.CreateTask(eventtask.NewEnqueueArticle(article))

		if err != nil {
			log.Printf("error is %v", err)
		}
	}

	return nil
}

func createQuery() string {
	strTime := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	return fmt.Sprintf("changed_at: %s", strTime)
}
