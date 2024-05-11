package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/internal/repository/models"
	"github.com/timickb/narration-engine/pkg/db"
)

type transitionRepo struct {
	db *db.Database
}

func NewTransitionRepo(db *db.Database) *transitionRepo {
	return &transitionRepo{db: db}
}

// Save Сохранить переход в таблицу с историей переходов.
func (r *transitionRepo) Save(ctx context.Context, dto *domain.SaveTransitionDto) (uuid.UUID, error) {
	id := uuid.New()
	err := r.db.WithTxSupport(ctx).Create(models.NewDbTransition(id, dto)).Error
	if err != nil {
		return uuid.Nil, fmt.Errorf("save transition: %w", err)
	}
	return id, nil
}

// SetError Обновить поле error в сохраненном переходе.
func (r *transitionRepo) SetError(ctx context.Context, transitionId uuid.UUID, errText string) error {
	query := r.db.WithTxSupport(ctx).Model(&models.DbTransition{}).
		Where("id = ?", transitionId).
		UpdateColumn("error", errText).
		UpdateColumn("failed", true)
	if query.Error != nil {
		return fmt.Errorf("set error to transition: %w", query.Error)
	}
	if query.RowsAffected == 0 {
		return fmt.Errorf("transition %s not found", transitionId.String())
	}
	return nil
}

// GetLastForInstance Получить последний сохраненный переход для экземпляра.
func (r *transitionRepo) GetLastForInstance(ctx context.Context, instanceId uuid.UUID) (*domain.SavedTransition, error) {
	var transition models.DbTransition

	err := r.db.WithTxSupport(ctx).Model(&models.DbTransition{}).
		Where("instance_id = ?", instanceId).
		Order("created_at DESC").
		Limit(1).
		Take(&transition).Error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no transitions found for instance %s", instanceId.String())
		}
		return nil, fmt.Errorf("get last transition for instance: %w", err)
	}

	return transition.ToDomain(), nil
}
