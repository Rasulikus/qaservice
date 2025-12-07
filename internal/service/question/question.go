package question

import (
	"context"

	"github.com/Rasulikus/qaservice/internal/model"
	"github.com/Rasulikus/qaservice/internal/repository"
	"github.com/Rasulikus/qaservice/internal/service"
)

type Service struct {
	questions repository.QuestionRepository
}

func NewService(repo repository.QuestionRepository) *Service {
	return &Service{
		questions: repo,
	}
}

func (s *Service) Create(ctx context.Context, input service.CreateQuestionInput) (*model.Question, error) {
	q := &model.Question{Text: input.Text}
	if err := s.questions.Create(ctx, q); err != nil {
		return nil, err
	}
	return q, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*model.Question, error) {
	q, err := s.questions.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (s *Service) List(ctx context.Context) ([]*model.Question, error) {
	qs, err := s.questions.List(ctx)
	if err != nil {
		return nil, err
	}
	return qs, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if err := s.questions.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
