package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/pkg/utils"
	"sort"
	"time"
)

// DbInstance Репрезентация таблицы instances.
type DbInstance struct {
	Id                 uuid.UUID
	ParentId           *uuid.UUID
	ScenarioName       string
	ScenarioVersion    string
	CurrentState       string
	PreviousState      string
	CurrentStateStatus string
	Context            string `gorm:"type:jsonb"`
	BlockingKey        *string
	LockedBy           *string
	LockedTill         *time.Time
	StartAfter         *time.Time
	Retries            int
	Failed             bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// FilledInstance Сущность instance, насыщенная вытянутыми записями из pending_events.
type FilledInstance struct {
	DbInstance
	Events []byte
}

// DbInstanceUpdate Структура с полями для обновления сущности.
type DbInstanceUpdate struct {
	Failed             *bool
	Retries            *int
	Context            *string
	CurrentState       *string
	PreviousState      *string
	CurrentStateStatus *string
	StartAfter         **time.Time
	UpdatedAt          *time.Time
	LastTransitionAt   *time.Time
}

// NewDbInstance Создать запись в таблице instances.
func NewDbInstance(id uuid.UUID, dto *domain.CreateInstanceDto) *DbInstance {
	return &DbInstance{
		Id:                 id,
		ParentId:           dto.ParentId,
		ScenarioName:       dto.ScenarioName,
		ScenarioVersion:    dto.ScenarioVersion,
		Context:            string(dto.Context),
		CurrentStateStatus: string(domain.StateStatusHandlerExecuted),
		BlockingKey:        dto.BlockingKey,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

// NewDbInstanceUpdate Создать структуру для обновления экземпляра.
func NewDbInstanceUpdate(instance *domain.Instance) *DbInstanceUpdate {
	return &DbInstanceUpdate{
		Failed:             utils.Ptr(instance.Failed),
		Retries:            utils.Ptr(instance.Retries),
		Context:            utils.Ptr(instance.Context.String()),
		CurrentState:       utils.Ptr(instance.CurrentState.Name),
		PreviousState:      utils.Ptr(instance.PreviousState.Name),
		CurrentStateStatus: utils.Ptr(string(instance.CurrentStateStatus)),
		StartAfter:         utils.Ptr(instance.GetStartAfter()),
		UpdatedAt:          utils.Ptr(time.Now()),
		LastTransitionAt:   instance.LastTransitionAt,
	}
}

// TableName Получить имя таблицы.
func (i *DbInstance) TableName() string {
	return "instances"
}

// ToDomain Преобразовать структуру в доменную сущность.
func (i *DbInstance) ToDomain() (*domain.Instance, error) {
	instanceCtx, err := domain.NewInstanceContext([]byte(i.Context))
	if err != nil {
		return nil, fmt.Errorf("NewInstanceContext: %w", err)
	}
	return &domain.Instance{
		Id: i.Id,
		Scenario: &domain.Scenario{
			Name:    i.ScenarioName,
			Version: i.ScenarioVersion,
		},
		CurrentState:       &domain.State{Name: i.CurrentState},
		PreviousState:      &domain.State{Name: i.PreviousState},
		CurrentStateStatus: domain.StateStatus(i.CurrentStateStatus),
		Context:            instanceCtx,
		Retries:            i.Retries,
		Failed:             i.Failed,
		LockedBy:           i.LockedBy,
		LockedTill:         i.LockedTill,
		BlockingKey:        i.BlockingKey,
		CreatedAt:          i.CreatedAt,
	}, nil
}

// ToDomain Преобразовать насыщенную структуру в доменную сущность.
func (i *FilledInstance) ToDomain() (*domain.Instance, error) {
	instance, err := i.DbInstance.ToDomain()
	if err != nil {
		return nil, err
	}
	instance.PendingEvents = &domain.EventsQueue{}

	if len(i.Events) != 0 {
		var rawEvents DbPendingEvents
		if err := json.Unmarshal(i.Events, &rawEvents); err != nil {
			return nil, fmt.Errorf("unmarshall events json: %w", err)
		}

		events := rawEvents.ToDomain()

		sort.Slice(events, func(i, j int) bool {
			return events[i].ExecutedAt.Before(events[j].ExecutedAt)
		})

		for _, event := range events {
			event.FromDb = true
			instance.PendingEvents.Enqueue(event)
		}
	}

	return instance, nil
}
