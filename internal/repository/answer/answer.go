package answer

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

func (r *Repository) Create(ctx context.Context, answer *model.Answer) error {
	return r.db.WithContext(ctx).Create(answer).Error
}

func (r *Repository) GetByID(ctx context.Context, id int) (*model.Answer, error) {
	a := new(model.Answer)
	if err := r.db.WithContext(ctx).First(a, id).Error; err != nil {
		return nil, repository.TranslateError(err)
	}
	return a, nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&model.Answer{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrNotFound
	}
	return nil
}
