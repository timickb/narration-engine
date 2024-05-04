package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
	"time"
)

// DbInstance Репрезентация таблицы instances.
type DbInstance struct {
	Id              uuid.UUID
	ParentId        *uuid.UUID
	ScenarioName    string
	ScenarioVersion string
	CurrentState    string
	PreviousState   string
	Context         string `gorm:"type:jsonb"`
	BlockingKey     *string
	LockedBy        *string
	LockedTill      *time.Time
	Retries         int
	Failed          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// FilledInstance Сущность instance, насыщенная вытянутыми записями из pending_events.
type FilledInstance struct {
	DbInstance
	Events []byte
}

// DbInstanceUpdate Структура с полями для обновления сущности.
type DbInstanceUpdate struct {
}

// NewDbInstance Создать запись в таблице instances.
func NewDbInstance(id uuid.UUID, dto *domain.CreateInstanceDto) *DbInstance {
	return &DbInstance{
		Id:              id,
		ParentId:        dto.ParentId,
		ScenarioName:    dto.ScenarioName,
		ScenarioVersion: dto.ScenarioVersion,
		CurrentState:    domain.StateStart.Name,
		PreviousState:   "",
		Context:         string(dto.Context),
		BlockingKey:     dto.BlockingKey,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// TableName Получить имя таблицы.
func (i *DbInstance) TableName() string {
	return "instances"
}

// ToDomain Преобразовать структуру в доменную сущность.
func (i *DbInstance) ToDomain() *domain.Instance {
	return &domain.Instance{
		Id: i.Id,
		Scenario: &domain.Scenario{
			Name:    i.ScenarioName,
			Version: i.ScenarioVersion,
		},
		CurrentState:  &domain.State{Name: i.CurrentState},
		PreviousState: &domain.State{Name: i.PreviousState},
		Context:       i.Context,
		Retries:       i.Retries,
		Failed:        i.Failed,
		LockedBy:      i.LockedBy,
		LockedTill:    i.LockedTill,
		BlockingKey:   i.BlockingKey,
		CreatedAt:     i.CreatedAt,
	}
}

// ToDomain Преобразовать насыщенную структуру в доменную сущность.
func (i *FilledInstance) ToDomain() (*domain.Instance, error) {
	instance := i.DbInstance.ToDomain()

	if len(i.Events) != 0 {
		var rawEvents DbPendingEvents
		if err := json.Unmarshal(i.Events, &rawEvents); err != nil {
			return nil, fmt.Errorf("unmarshall events json: %w", err)
		}

		events := rawEvents.ToDomain()
		eventsQueue := &domain.EventsQueue{}

		for _, event := range events {
			if err := eventsQueue.Enqueue(&domain.EventPushDto{
				EventName: event.EventName,
				Params:    event.EventParams,
				External:  event.External,
			}); err != nil {
				return nil, fmt.Errorf("push to event queue: %w", err)
			}
		}
		instance.PendingEvents = eventsQueue
	}

	return instance, nil
}
