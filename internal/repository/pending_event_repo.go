package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/internal/repository/models"
	"github.com/timickb/narration-engine/pkg/db"
)

type pendingEventRepo struct {
	db *db.Database
}

func NewPendingEventRepo(db *db.Database) *pendingEventRepo {
	return &pendingEventRepo{db: db}
}

// Create Поставить событие в очередь.
func (r *pendingEventRepo) Create(ctx context.Context, dto *domain.CreatePendingEventDto) (uuid.UUID, error) {
	eventId, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil, fmt.Errorf("new uuid: %w", err)
	}

	event := models.NewDbPendingEvent(eventId, dto)
	if err := r.db.WithTxSupport(ctx).Create(event).Error; err != nil {
		return uuid.Nil, fmt.Errorf("create pending event: %w", err)
	}
	return eventId, nil
}
