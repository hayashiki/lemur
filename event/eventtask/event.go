package eventtask

import (
	"github.com/hayashiki/lemur/entity"
	"github.com/hayashiki/lemur/event"
)

var TaskName = "docBase"
var TaskPath = "/enqueue/articles"

func NewEnqueueArticle(v *entity.Article) event.Task {
	return event.Task{
		Queue:  TaskName,
		Path:   TaskPath,
		Object: v,
	}
}
