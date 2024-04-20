package domain

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

// EventsQueue Очередь событий на обработку.
type EventsQueue struct {
	size  int
	front *PendingEvent
	back  *PendingEvent
}

// Size Получить текущий размер очереди.
func (q *EventsQueue) Size() int {
	return q.size
}

// Enqueue Положить событие в очередь на обработку.
func (q *EventsQueue) Enqueue(dto *EventPushDto) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("uuid generate: %w", err)
	}

	event := &PendingEvent{
		Id:          id,
		EventName:   dto.EventName,
		EventParams: dto.Params,
		External:    dto.External,
		CreatedAt:   time.Now(),
	}

	if q.size == 0 {
		q.front = event
		q.back = event
	} else {
		q.back.Next = event
		q.back = event
	}
	q.size++

	return nil
}

// Dequeue Достать первое событие из очереди и удалить его.
func (q *EventsQueue) Dequeue() *PendingEvent {
	if q.size == 0 {
		return nil
	}
	front := q.front
	front.Next = nil
	q.front = q.front.Next
	return front
}

// Front Достать первое событие из очереди.
func (q *EventsQueue) Front() *PendingEvent {
	return q.front
}
