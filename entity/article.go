package entity

import "time"

const ArticleKind = "Article"

type Article struct {
	ID        int64
	UserId    int64
	Title     string
	Number    int64
	Message   string
	Url       string
	MdGcsPath string
	ImageURLs []*Image
	MDBody    string `datastore:",noindex"`
	CreatedAt time.Time
}

type Image struct {
	ID        string
	Name      string
	Content   []byte
	OriginURL string
	URL       string
}

func NewArticle(id int64, title, body string, createdAt time.Time, images []*Image) *Article {
	return &Article{
		ID:        id,
		Title:     title,
		MDBody:    body,
		CreatedAt: createdAt,
		ImageURLs: images,
	}
}

func NewImage(id, name string, url string) *Image {
	return &Image{
		ID:   id,
		Name: name,
		URL:  url,
	}
}
