package app

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/hayashiki/lemur/config"
	"github.com/hayashiki/lemur/docbase"
	"github.com/hayashiki/lemur/entity"
	"github.com/hayashiki/lemur/event"
	"github.com/hayashiki/lemur/event/eventtask"
	"github.com/hayashiki/lemur/infra"
	"github.com/hayashiki/lemur/logger"
	"github.com/hayashiki/lemur/usecase"
	"log"
	"net/http"
)

type server struct {
	logger      logger.Logger
	docBase     docbase.DocBaser
	taskQueue   event.TaskQueue
	articleRepo entity.ArticleRepository
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
	return &server{
		logger:      logger.NewLogger(),
		docBase:     docbase.NewClient(infra.DocBaseClient(c.DocBaseTeam, c.DocBaseToken)),
		taskQueue:   event.NewTasksClient(c.GCPProjectID, c.GCPLocationID),
		articleRepo: entity.NewArticleRepository(infra.GetDSClient(c.GCPProjectID)),
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func (s *server) cronFetchDocBase(w http.ResponseWriter, r *http.Request) {
	article := usecase.NewArticles(
		s.logger,
		s.docBase,
		s.taskQueue,
		s.articleRepo,
		)

	if err := article.Do(); err != nil {
		// curryの参考にちゃんとする
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
		s.docBase,
		s.articleRepo,
	)

	log.Printf("article %+v", article)

	// docBase
	taskName := r.Header.Get("X-Appengine-Taskname")

	if taskName != eventtask.TaskName {
		http.Error(w, "Invalid task name", http.StatusInternalServerError)
	}

	//TODO: handle Queuename if need
	//244011105321548948x
	log.Printf("header X-Appengine-Queuename %+v", r.Header.Get("X-Appengine-Queuename"))

	params := usecase.EnqueueArticlesInputParams{Article: article}
	if err := enqueueArticle.Do(params); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
