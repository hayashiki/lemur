package usecase

import (
	"fmt"
	"github.com/hayashiki/lemur/docbase"
	"github.com/hayashiki/lemur/entity"
	"github.com/hayashiki/lemur/github"
	"github.com/hayashiki/lemur/logger"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
)

type enqueueArticle struct {
	log        logger.Logger
	docBaseSvc docbase.Client
	repo       entity.ArticleRepository
	githubSvc  github.Client
}

func NewEnqueueArticle(
	log logger.Logger,
	docBase docbase.Client,
	repo entity.ArticleRepository,
	githubSvc github.Client,
) *enqueueArticle {
	return &enqueueArticle{
		log,
		docBase,
		repo,
		githubSvc,
	}
}

func (uc *enqueueArticle) Do(params EnqueueArticlesInputParams) error {
	uc.log.Debug("enqueueArticle.Do called")

	files, err := uc.getAsync(params.Article.Attachments)

	if err != nil {
		return err
	}

	files = append(files, &github.File{
		Path:    "",
		Name:    params.Article.Title,
		Content: []byte(params.Article.MDBody),
	})

	message   := fmt.Sprintf("Add %s", params.Article.Title)
	gitSvc := github.NewGitCommit(files, message)
	if err := uc.githubSvc.PushNewFileToBranch(gitSvc); err != nil {
		return err
	}

	if err := uc.githubSvc.CreateNewPullRequest(gitSvc); err != nil {
		return err
	}

	return nil
}

func (uc *enqueueArticle) getAsync(attachments []*entity.Attachment) ([]*github.File, error) {
	var files []*github.File

	eg := errgroup.Group{}
	mutex := &sync.Mutex{}
	for _, attachment := range attachments {
		attachment := attachment
		eg.Go(func() error {
			file, err := uc.get(attachment.ID, attachment.Name)
			if err != nil {
				log.Printf("GetMultiAsync err: %w", err)
				return err
			}
			mutex.Lock()
			files = append(files, file)
			mutex.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		log.Printf("eg.Wait err: %w", err)
		return files, err
	}
	return files, nil
}

func (uc *enqueueArticle) get(id, name string) (*github.File, error) {
	fileBytes, err := uc.docBaseSvc.Download(id)
	if err != nil {
		log.Printf("GetMultiAsync err: %w", err)
		return nil, err
	}
	file := &github.File{
		Path:    "attachments/",
		Name:    name,
		Content: fileBytes,
	}

	return file, nil
}
