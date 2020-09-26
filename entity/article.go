package entity

import "time"

const ArticleKind = "Article"

type Article struct {
	ID          int64
	UserId      int64
	Title       string
	Number      int64
	Message     string
	Url         string
	MdGcsPath   string
	Attachments []*Attachment
	MDBody      string `datastore:",noindex"`
	CreatedAt   time.Time
}

type Attachment struct {
	ID        string
	Name      string
	Content   []byte
	OriginURL string
	URL       string
}

func NewArticle(id int64, title, body string, createdAt time.Time, attachments []*Attachment) *Article {
	return &Article{
		ID:          id,
		Title:       title,
		MDBody:      body,
		CreatedAt:   createdAt,
		Attachments: attachments,
	}
}

func NewAttachment(id, name string, url string) *Attachment {
	return &Attachment{
		ID:   id,
		Name: name,
		URL:  url,
	}
}
