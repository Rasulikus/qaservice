package answer

import (
	"context"

	"github.com/Rasulikus/qaservice/internal/model"
	"github.com/Rasulikus/qaservice/internal/repository"
	"github.com/Rasulikus/qaservice/internal/service"
)

type Service struct {
	answers   repository.AnswerRepository
	questions repository.QuestionRepository
}

func NewService(answerRepo repository.AnswerRepository, questionRepo repository.QuestionRepository) *Service {
	return &Service{
		answers:   answerRepo,
		questions: questionRepo,
	}
}

func (s *Service) Create(ctx context.Context, input service.CreateAnswerInput) (*model.Answer, error) {
	// проверка существует ли вопрос
	if _, err := s.questions.GetByID(ctx, input.QuestionID); err != nil {
		return nil, err
	}

	a := &model.Answer{
		QuestionID: input.QuestionID,
		UserID:     input.UserID,
		Text:       input.Text,
	}
	if err := s.answers.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*model.Answer, error) {
	a, err := s.answers.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if err := s.answers.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
