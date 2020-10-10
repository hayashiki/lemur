package entity

import (
	"cloud.google.com/go/datastore"
	"context"
)

type repository struct {
	ctx      context.Context
	dsClient *datastore.Client
}

//go:generate mockgen -source article_repository.go -destination mock_article/mock_article_repository.go
type ArticleRepository interface {
	Get(id int64) (*Article, error)
	Put(*Article) (int64, error)
	Delete(id int64) error
}

func NewArticleRepository(dsClient *datastore.Client) ArticleRepository {
	ctx := context.Background()
	return &repository{ctx, dsClient}
}

func (r *repository) key(id int64) *datastore.Key {
	return datastore.IDKey(ArticleKind, id, nil)
}

func (r *repository) Get(id int64) (*Article, error) {
	var a Article
	key := datastore.IDKey(ArticleKind, id, nil)
	err := r.dsClient.Get(r.ctx, key, &a)
	return &a, err
}

func (r *repository) Put(a *Article) (int64, error) {
	key := datastore.IDKey(ArticleKind, a.ID, nil)
	_, err := r.dsClient.Put(r.ctx, key, a)
	if err != nil {
		return int64(0), err
	}
	return key.ID, err
}

func (r *repository) Delete(id int64) error {
	//TODO: refine
	//key := datastore.NewKey(r.ctx, article.Kind, "", id, nil)
	//err := datastore.Delete(r.ctx, key)
	return nil
}
