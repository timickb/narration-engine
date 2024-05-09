package models

import (
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
	"time"
)

// DbTransition Репрезентация таблицы transitions
type DbTransition struct {
	Id         uuid.UUID
	InstanceId uuid.UUID
	StateFrom  string
	StateTo    string
	EventName  string
	Params     string `gorm:"type:jsonb"`
	Failed     bool
	Error      string
	CreatedAt  time.Time
}

func NewDbTransition(id uuid.UUID, dto *domain.SaveTransitionDto) *DbTransition {
	return &DbTransition{
		Id:         id,
		InstanceId: dto.InstanceId,
		StateFrom:  dto.StateFrom,
		StateTo:    dto.StateTo,
		EventName:  dto.EventName,
		Params:     dto.EventParams,
		CreatedAt:  time.Now(),
	}
}

func (t *DbTransition) TableName() string {
	return "transitions"
}

func (t *DbTransition) ToDomain() *domain.SavedTransition {
	return &domain.SavedTransition{
		EventName:   t.EventName,
		EventParams: t.Params,
		StateFrom:   t.StateFrom,
		StateTo:     t.StateTo,
		Error:       t.Error,
	}
}
