package service

import (
	"context"

	"github.com/Rasulikus/qaservice/internal/model"
)

type CreateQuestionInput struct {
	Text string
}

type QuestionService interface {
	Create(ctx context.Context, input CreateQuestionInput) (*model.Question, error)
	GetByID(ctx context.Context, id int) (*model.Question, error)
	List(ctx context.Context) ([]*model.Question, error)
	Delete(ctx context.Context, id int) error
}

type CreateAnswerInput struct {
	QuestionID int
	UserID     string
	Text       string
}

type AnswerService interface {
	Create(ctx context.Context, input CreateAnswerInput) (*model.Answer, error)
	GetByID(ctx context.Context, id int) (*model.Answer, error)
	Delete(ctx context.Context, id int) error
}
