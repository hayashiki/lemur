package infra

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/hayashiki/docbase-go"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
)

func GetDSClient(projectID string) *datastore.Client {
	client, err := datastore.NewClient(
		context.Background(), projectID,
	)

	if err != nil {
		log.Panic(err)
	}
	return client
}

func DocBaseClient(team, token string) *docbase.Client {
	return docbase.NewClient(nil, team, token)
}

func NewESClient() (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetURL(os.Getenv("ES_URL")),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
}
