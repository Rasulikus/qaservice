package question

import (
	"context"

	"github.com/Rasulikus/qaservice/internal/model"
	"github.com/Rasulikus/qaservice/internal/repository"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, question *model.Question) error {
	return r.db.WithContext(ctx).Create(question).Error
}

func (r *Repository) GetByID(ctx context.Context, id int) (*model.Question, error) {
	q := new(model.Question)
	if err := r.db.WithContext(ctx).Preload("Answers").First(q, id).Error; err != nil {
		return nil, repository.TranslateError(err)
	}
	return q, nil
}

func (r *Repository) List(ctx context.Context) ([]*model.Question, error) {
	var qs []*model.Question
	if err := r.db.WithContext(ctx).Find(&qs).Error; err != nil {
		return nil, err
	}
	return qs, nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&model.Question{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrNotFound
	}
	return nil
}
