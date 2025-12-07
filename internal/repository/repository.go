package repository

import (
	"context"

	"github.com/Rasulikus/qaservice/internal/model"
)

type AnswerRepository interface {
	Create(ctx context.Context, answer *model.Answer) error
	GetByID(ctx context.Context, id int) (*model.Answer, error)
	Delete(ctx context.Context, id int) error
}

type QuestionRepository interface {
	Create(ctx context.Context, question *model.Question) error
	GetByID(ctx context.Context, id int) (*model.Question, error)
	List(ctx context.Context) ([]*model.Question, error)
	Delete(ctx context.Context, id int) error
}
