package eventtask

import (
	"github.com/hayashiki/lemur/entity"
	"github.com/hayashiki/lemur/event"
)

var TaskName = "docBase"


func NewEnqueueArticle(v *entity.Article) event.Task {
	return event.Task{
		Queue:  TaskName,
		Path:   "/enqueue/articles",
		Object: v,
	}
}
