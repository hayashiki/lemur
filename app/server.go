package app

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/hayashiki/lemur/config"
	"github.com/hayashiki/lemur/docbase"
	"github.com/hayashiki/lemur/elasticsearch"
	"github.com/hayashiki/lemur/entity"
	"github.com/hayashiki/lemur/event"
	"github.com/hayashiki/lemur/event/eventtask"
	"github.com/hayashiki/lemur/github"
	"github.com/hayashiki/lemur/infra"
	"github.com/hayashiki/lemur/logger"
	"github.com/hayashiki/lemur/usecase"
	"log"
	"net/http"
	"os"
)

type server struct {
	logger      logger.Logger
	docBaseSvc  docbase.Client
	taskQueue   event.TaskQueue
	articleRepo entity.ArticleRepository
	githubSvc   github.Client
	esRepo      elasticsearch.Repository
}

func (s *server) Router() http.Handler {
	r := chi.NewRouter()

	r.Route("/cron", func(r chi.Router) {
		r.Get("/docbase", s.cronFetchDocBase)
	})

	r.Route("/enqueue", func(r chi.Router) {
		r.Post("/articles", s.enqueueArticles)
	})

	r.Get("/", healthCheckHandler)

	return r
}

type Server interface {
	Router() http.Handler
}

func NewServer(c *config.Config) Server {

	esClient, err := infra.NewESClient()
	if err != nil {
	}

	esRepo := elasticsearch.NewRepository(esClient)

	return &server{
		logger:      logger.NewLogger(),
		docBaseSvc:  docbase.NewClient(infra.DocBaseClient(c.DocBaseTeam, c.DocBaseToken)),
		taskQueue:   event.NewTasksClient(c.GCPProjectID, c.GCPLocationID),
		articleRepo: entity.NewArticleRepository(infra.GetDSClient(c.GCPProjectID)),
		githubSvc:   github.NewClient(c.GithubOrg, c.GithubRepo, c.GithubSecret, "hayashiki"),
		esRepo:      esRepo,
	}
}

func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "ok")
}

func (s *server) cronFetchDocBase(w http.ResponseWriter, r *http.Request) {

	err := validateCron(r, os.Getenv("AUTHORIZATION"))

	if err != nil {
		return
	}

	article := usecase.NewArticles(
		s.logger,
		s.docBaseSvc,
		s.taskQueue,
		s.articleRepo,
	)

	if err := article.Do(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// enqueueされたarticleを処理する
// usecaseにわたすために、articleを引数にする
// そのためにtaskをパースする必要がある
func (s *server) enqueueArticles(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	article := entity.Article{}

	if err := event.ParseTask(r, &article); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	enqueueArticle := usecase.NewEnqueueArticle(
		s.logger,
		s.docBaseSvc,
		s.articleRepo,
		s.githubSvc,
		s.esRepo,
	)

	log.Printf("article %+v", article.ID)

	// docBaseSvc
	taskName := r.Header.Get("X-Appengine-Taskname")

	if taskName != eventtask.TaskName {
		http.Error(w, "Invalid task name", http.StatusInternalServerError)
	}

	//TODO: handle Queuename if need
	// eg.244011105321548948x
	log.Printf("header X-Appengine-Queuename %+v", r.Header.Get("X-Appengine-Queuename"))

	params := usecase.EnqueueArticlesInputParams{Article: article}
	if err := enqueueArticle.Do(params); err != nil {
		log.Printf("enqueueArticle.Do %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateCron(r *http.Request, secrets string) error {
	if r.Header.Get("X-Appengine-Cron") == "" {
		return fmt.Errorf("cron request does not have X-Appengine-Cron header")
	}

	if r.Header.Get("Authorization") == "secrets" {
		return fmt.Errorf("cron request does not have Authorization header")
	}
	return nil
}
