package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/hayashiki/lemur/entity"
	"github.com/olivere/elastic/v7"
	"log"
	"strconv"
	"strings"
	"time"
)

const articleIndex = "article"

func newTermString(name string, input []string) *elastic.TermsQuery {
	values := make([]interface{}, len(input))

	for i, s := range input {
		values[i] = s
	}

	return elastic.NewTermsQuery(name, values...)
}

// ElasticQuery struct
type ElasticQuery struct {
	Index    string
	Query    *elastic.BoolQuery
	SortInfo elastic.SortInfo
	From     int
	Size     int
}

func CreateSearchQuery(q map[string]string) *ElasticQuery {
	query := elastic.NewBoolQuery()
	if category, ok := q["category"]; ok && len(category) > 0 {
		query = query.Filter(newTermString("category", strings.Split(category, ",")))
	}

	if keywords, ok := q["keywords"]; ok && len(keywords) > 0 {
		for _, keyword := range keywords {
			query = query.Must(elastic.NewMatchPhraseQuery("search_text", keyword))
		}
	}

	from := 0
	if offset, ok := q["offset"]; ok {
		if v, err := strconv.Atoi(offset); err == nil {
			from = v
		}
	}
	size := 36
	if limit, ok := q["limit"]; ok {
		if v, err := strconv.Atoi(limit); err == nil {
			size = v
		}
	}

	sort := elastic.SortInfo{Field: "title", Ascending: false}

	return &ElasticQuery{
		Index:    articleIndex,
		Query:    query,
		SortInfo: sort,
		From:     from,
		Size:     size,
	}
}

//go:generate mockgen -source repository.go -destination mock/mock_repository.go
type Repository interface {
	SaveArticleDocument(input *entity.Article) error
}

type repository struct {
	esClient *elastic.Client
}

func NewRepository(esClient *elastic.Client) *repository {
	return &repository{esClient: esClient}
}

type articleDocuments struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *repository) SaveArticleDocument(input *entity.Article) error {
	ctx := context.Background()
	doc := articleDocuments{
		Title:       input.Title,
		Description: input.MDBody,
	}
	_, err := r.esClient.Index().Index(articleIndex).
		Id(strconv.Itoa(int(input.ID))).
		BodyJson(doc).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateArticleDocument(input *entity.Article) error {
	ctx := context.Background()
	doc := articleDocuments{
		Title:       input.Title,
		Description: input.MDBody,
		Category:    input.Category,
		UpdatedAt:   time.Now(),
	}
	res, err := r.esClient.Update().Index(articleIndex).
		Id(strconv.Itoa(int(input.ID))).
		Doc(doc).
		Do(ctx)
	if err != nil {
		return err
	}
	log.Printf("res is %v", res)
	return nil

}

func (r *repository) DeleteLinkDocument(input *entity.Article) error {
	ctx := context.Background()
	_, err := r.esClient.Delete().Index(articleIndex).
		Id(strconv.Itoa(int(input.ID))).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SearchDocuments(q map[string]string) ([]int, error) {
	ctx := context.Background()
	//query := elastic.NewMultiMatchQuery(word, "title", "description")

	qq := CreateSearchQuery(q)
	result, err := r.esClient.Search().
		Index(articleIndex).
		Query(qq.Query).
		Do(ctx)

	if err != nil {
		return nil, err
	}
	var ids []int
	if result.Hits.TotalHits.Value > 0 {
		for _, hit := range result.Hits.Hits {
			id, err := strconv.Atoi(hit.Id)
			if err != nil {
				return nil, err
			}
			var ad articleDocuments
			ids = append(ids, id)
			if err := json.Unmarshal(hit.Source, &ad); err != nil {
				log.Printf("failed to json unmarshal")
			}
			log.Printf("articleDocuments", ad.Title)
		}
	}
	return ids, nil
}
