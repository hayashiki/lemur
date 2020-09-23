package usecase

import (
	"github.com/hayashiki/lemur/docbase"
	"github.com/hayashiki/lemur/entity"
	"github.com/hayashiki/lemur/event"
	"github.com/hayashiki/lemur/event/eventtask"
	"github.com/hayashiki/lemur/logger"
	"log"
)

type Articles struct {
	log       logger.Logger
	docBase   docbase.DocBaser
	taskQueue event.TaskQueue
	repo      entity.ArticleRepository
}

type EnqueueArticlesInputParams struct {
	Article entity.Article
}

func NewArticles(
	log logger.Logger,
	docBase docbase.DocBaser,
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

	// 昨日の日付を引数にしたい
	posts, err := a.docBase.PostList()

	if err != nil {
		return err
	}

	if len(posts) == 0 {
		return nil
	}

	//var tasks []event.Task
	for _, p := range posts {

		var images []*entity.Image
		for _, at := range p.Attachments {
			img := entity.NewImage(at.ID, at.Name, at.URL)
			images = append(images, img)
		}

		article := entity.NewArticle(int64(p.ID), p.Title, p.Body, p.CreatedAt, images)
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
