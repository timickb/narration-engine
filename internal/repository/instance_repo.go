package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/timickb/go-stateflow/internal/domain"
	"github.com/timickb/go-stateflow/pkg/db"
)

type instanceRepo struct {
	db *db.Database
}

// NewInstanceRepo Создать репозиторий над экземплярами сценариев.
func NewInstanceRepo(db *db.Database) *instanceRepo {
	return &instanceRepo{db: db}
}

// Save Сохранить экземпляр сценария.
func (r *instanceRepo) Save(ctx context.Context, instance *domain.Instance) error {
	// TODO: implement
	panic("implement me")
}

// GetWaitingIds Получить список ид экземпляров, ожидающих дальнейшей работы.
func (r *instanceRepo) GetWaitingIds(ctx context.Context, limit int) ([]uuid.UUID, error) {
	// TODO: implement
	panic("implement me")
}

// GetById Получить данные экземпляра по идентификатору.
func (r *instanceRepo) GetById(ctx context.Context, id uuid.UUID) ([]*domain.Instance, error) {
	// TODO: implement
	panic("implement me")
}
