package usecase_test

import (
	"github.com/golang/mock/gomock"
	mock_docbase "github.com/hayashiki/lemur/docbase/mock"
	"github.com/hayashiki/lemur/entity"
	mock_entity "github.com/hayashiki/lemur/entity/mock_article"
	mock_github "github.com/hayashiki/lemur/github/mock"
	"github.com/hayashiki/lemur/logger"
	"github.com/hayashiki/lemur/usecase"
	"testing"
	"time"
)

func TestEnqueueArticleWithMock_Do(t *testing.T) {

	article := entity.Article{
		ID:    1530229,
		Title: "dummyTitle",
		Attachments: []*entity.Attachment{
			&entity.Attachment{
				ID:   "1c9c135b-61cd-4f13-8b07-6358691f782d.png",
				Name: "Diagram.png",
				URL:  "https://image.docbase.io/uploads/1c9c135b-61cd-4f13-8b07-6358691f782d.png",
			},
			&entity.Attachment{
				ID:   "b9dea11-a820-45b5-bfc4-e07be2ef6588.png",
				Name: "image.png",
				URL:  "https://image.docbase.io/uploads/b9dea11-a820-45b5-bfc4-e07be2ef6588.png",
			},
		},
		MDBody:    "this is a dummy document",
		CreatedAt: time.Time{},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dcbCli := mock_docbase.NewMockClient(ctrl)
	dcbCli.EXPECT().Download("1b9c135b-61cd-4f13-8b07-6358691f782c.png").Return([]byte("hoge"), nil)
	dcbCli.EXPECT().Download("b9dfea10-a820-45a5-bfc4-e07be2ef6588.png").Return([]byte("hoge"), nil)

	ghCli := mock_github.NewMockClient(ctrl)
	ghCli.EXPECT().CreateNewPullRequest(gomock.Any()).Return(nil).Times(1)
	ghCli.EXPECT().PushNewFileToBranch(gomock.Any()).Return(nil).Times(1)

	artRepo := mock_entity.NewMockArticleRepository(ctrl)

	enqueueArticle := usecase.NewEnqueueArticle(
		logger.NewLogger(),
		dcbCli,
		artRepo,
		ghCli,
	)

	param := usecase.EnqueueArticlesInputParams{Article: article}

	if err := enqueueArticle.Do(param); err != nil {
		t.Logf("enqueueArticle.Do err: %v", err)
		t.Fail()
	}
}
