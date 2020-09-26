package docbase

import (
	"errors"
	"github.com/hayashiki/docbase-go"
	"net/http"
)

type client struct {
	cli *docbase.Client
}

func (d *client) Download(attachmentID string) ([]byte, error) {
	fileBytes, resp, err := d.cli.Attachments.Download(attachmentID)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response status code is not ok")
	}

	return fileBytes.Write(), err
}

func NewClient(cli *docbase.Client) Client {
	return &client{
		cli: cli,
	}
}

//go:generate mockgen -source docbase.go -destination mock/mock_docbase.go
type Client interface {
	PostList(q string) ([]*docbase.Post, error)
	Download(attachmentID string) ([]byte, error)
}

func (d *client) PostList(q string) ([]*docbase.Post, error) {
	opt := &docbase.PostListOptions{
		Q:       q,
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
