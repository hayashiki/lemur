package docbase

import (
	"github.com/hayashiki/docbase-go"
	"net/http"
)

type client struct {
	cli *docbase.Client
}

func (d *client) Download(attachmentID string) ([]byte, error) {
	panic("implement me")
}

func NewClient(cli *docbase.Client) DocBaser {
	return &client{
		cli: cli,
	}
}

type DocBaser interface {
	PostList() ([]*docbase.Post, error)
	Download(attachmentID string) ([]byte, error)
}

func (d *client) PostList() ([]*docbase.Post, error) {
	opt := &docbase.PostListOptions{
		Q:       "",
		Page:    1,
		PerPage: 10,
	}

	posts, resp, err := d.cli.Posts.List(opt)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return posts, nil
}
