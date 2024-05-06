package models

import (
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/pkg/utils"
	"time"
)

// DbPendingEvent Репрезентация таблицы pending_events.
type DbPendingEvent struct {
	Id         uuid.UUID `json:"id,omitempty"`
	InstanceId uuid.UUID `json:"instance_id,omitempty"`
	EventName  string    `json:"event_name,omitempty"`
	Params     string    `gorm:"type:jsonb" json:"params,omitempty"`
	External   bool      `json:"external,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	ExecutedAt time.Time `json:"executed_at,omitempty"`
}

// DbPendingEvents Список репрезентаций таблицы pending_events.
type DbPendingEvents []*DbPendingEvent

// NewDbPendingEvent Добавить запись в таблицу pending_events.
func NewDbPendingEvent(id uuid.UUID, dto *domain.CreatePendingEventDto) *DbPendingEvent {
	return &DbPendingEvent{
		Id:         id,
		InstanceId: dto.InstanceId,
		EventName:  dto.Name,
		Params:     string(dto.Params),
		External:   dto.External,
		CreatedAt:  time.Now(),
	}
}

func NewDbPendingEventFromDomain(event *domain.PendingEvent, instanceId uuid.UUID) *DbPendingEvent {
	return &DbPendingEvent{
		Id:         event.Id,
		InstanceId: instanceId,
		EventName:  event.EventName,
		Params:     event.EventParams,
		External:   event.External,
		CreatedAt:  event.CreatedAt,
		ExecutedAt: event.ExecutedAt,
	}
}

// TableName Получить имя таблицы.
func (e *DbPendingEvent) TableName() string {
	return "pending_events"
}

// ToDomain Преобразовать в доменную структуру.
func (e *DbPendingEvent) ToDomain() *domain.PendingEvent {
	return &domain.PendingEvent{
		Id:          e.Id,
		EventName:   e.EventName,
		EventParams: e.Params,
		External:    e.External,
		CreatedAt:   e.CreatedAt,
		ExecutedAt:  e.ExecutedAt,
	}
}

// ToDomain Преобразовать в список доменных структур.
func (events DbPendingEvents) ToDomain() []*domain.PendingEvent {
	return utils.MapSlice(events, (*DbPendingEvent).ToDomain)
}
