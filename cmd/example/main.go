package main

import (
	"encoding/json"
	"github.com/hayashiki/lemur/elasticsearch"
	"github.com/hayashiki/lemur/entity"
	"github.com/hayashiki/lemur/infra"
	"github.com/olivere/elastic/v7"
	"log"
	"time"
)

var esClient *elastic.Client

func init() {
	esClient, _ = infra.NewESClient()
}

func main() {
	get("something")
}

func query() {

	testCase := func(q map[string]string) {
		query := elasticsearch.CreateSearchQuery(q)

		s, err := query.Query.Source()

		if err != nil {

		}

		j, err := json.Marshal(s)

		if err != nil {

		}

		source := string(j)

		log.Printf("source %s", source)
	}

	testCase(map[string]string{
		"category": "tech",
	})
}

func update() {
	esRepo := elasticsearch.NewRepository(esClient)

	art := &entity.Article{
		ID:        1350176,
		Title:     "hogezzzzz",
		Category:  "tech",
		MDBody:    "example ## body ## hoge",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := esRepo.UpdateArticleDocument(art)

	if err != nil {
		log.Printf("err is %v", err)

	}
}

func get(keyword string) {
	query := map[string]string{
		"keyword": keyword,
	}

	esRepo := elasticsearch.NewRepository(esClient)

	ids, err := esRepo.SearchDocuments(query)

	if err != nil {
		log.Printf("err is %v", err)

	}

	for _, id := range ids {
		log.Printf("id is %d", id)
	}
}
