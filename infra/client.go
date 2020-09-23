package infra

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/hayashiki/docbase-go"
	"log"
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
