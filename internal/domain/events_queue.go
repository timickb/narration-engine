package domain

import (
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/pkg/utils"
)

// EventsQueue Очередь событий на обработку.
type EventsQueue struct {
	size    int
	front   *PendingEvent
	back    *PendingEvent
	shifted map[uuid.UUID]*PendingEvent
}

// Size Получить текущий размер очереди.
func (q *EventsQueue) Size() int {
	return q.size
}

// Enqueue Добавить событие в конец очереди.
func (q *EventsQueue) Enqueue(event *PendingEvent) {
	if q.size == 0 {
		q.front = event
		q.back = event
	} else {
		q.back.Next = event
		q.back = event
	}
	q.size++
}

// PushToFront Добавить событие в начало очереди (приоритетно).
func (q *EventsQueue) PushToFront(event *PendingEvent) {
	event.Next = q.front
	q.front = event
	q.size++
}

// Dequeue Достать первое событие из очереди и удалить его.
func (q *EventsQueue) Dequeue() *PendingEvent {
	if q.size == 0 {
		return nil
	}
	front := q.front
	q.front = q.front.Next
	front.Next = nil
	if q.front == nil {
		q.back = nil
	}

	if q.shifted == nil {
		q.shifted = make(map[uuid.UUID]*PendingEvent)
	}
	q.shifted[front.Id] = front
	q.size--

	return front
}

// Front Достать первое событие из очереди.
func (q *EventsQueue) Front() *PendingEvent {
	return q.front
}

// GetShiftedIds Получить идентификаторы убранных из очереди событий.
func (q *EventsQueue) GetShiftedIds() []uuid.UUID {
	return utils.MapToKeysSlice(q.shifted)
}

// GetNewEvents Получить события, которых еще нет в БД.
func (q *EventsQueue) GetNewEvents() []*PendingEvent {
	result := make([]*PendingEvent, 0)

	for event := q.front; event != nil; event = event.Next {
		if !event.FromDb {
			result = append(result, event)
		}
	}
	return result
}
