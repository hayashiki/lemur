package usecase

import (
	"github.com/hayashiki/lemur/docbase"
	"github.com/hayashiki/lemur/entity"
	"github.com/hayashiki/lemur/logger"
)

type enqueueArticle struct {
	log     logger.Logger
	docBase docbase.DocBaser
	repo    entity.ArticleRepository
}

func NewEnqueueArticle(
	log logger.Logger,
	docBase docbase.DocBaser,
	repo entity.ArticleRepository,
) *enqueueArticle {
	return &enqueueArticle{
		log,
		docBase,
		repo,
	}
}

func (uc *enqueueArticle) Do(params EnqueueArticlesInputParams) error {
	uc.log.Debug("enqueueArticle.Do called")
	return nil
}
