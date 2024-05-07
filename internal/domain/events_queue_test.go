package domain

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

var (
	testPendingEvent1 = &PendingEvent{
		Id:          uuid.New(),
		EventName:   "test",
		EventParams: "{}",
		External:    false,
		FromDb:      true,
		CreatedAt:   time.Now().Add(-time.Minute * 2),
		ExecutedAt:  time.Now().Add(-time.Minute * 2),
	}
	testPendingEvent2 = &PendingEvent{
		Id:          uuid.New(),
		EventName:   "test2",
		EventParams: "{}",
		External:    false,
		FromDb:      false,
		CreatedAt:   time.Now().Add(-time.Minute),
		ExecutedAt:  time.Now().Add(-time.Minute),
	}
	testPendingEvent3 = &PendingEvent{
		Id:          uuid.New(),
		EventName:   "test3",
		EventParams: "{}",
		External:    false,
		FromDb:      false,
		CreatedAt:   time.Now(),
		ExecutedAt:  time.Now(),
	}
)

func TestEventsQueue_Enqueue(t *testing.T) {
	q := EventsQueue{}

	q.Enqueue(testPendingEvent1)
	assert.Equal(t, testPendingEvent1.EventName, q.front.EventName)
	assert.Equal(t, testPendingEvent1.EventName, q.back.EventName)
	assert.Equal(t, 1, q.size)

	q.Enqueue(testPendingEvent2)
	assert.Equal(t, testPendingEvent1.EventName, q.front.EventName)
	assert.Equal(t, testPendingEvent2.EventName, q.back.EventName)
	assert.Equal(t, 2, q.size)
}

func TestEventsQueue_PushToFront(t *testing.T) {
	q := EventsQueue{}
	q.Enqueue(testPendingEvent1)
	q.Enqueue(testPendingEvent2)
	q.PushToFront(testPendingEvent3)

	assert.Equal(t, testPendingEvent3.EventName, q.front.EventName)
	assert.Equal(t, 3, q.size)
}

func TestEventsQueue_Front(t *testing.T) {
	q := EventsQueue{}
	assert.Nil(t, q.Front())

	q.Enqueue(testPendingEvent1)
	assert.Equal(t, testPendingEvent1.EventName, q.Front().EventName)

	q.Enqueue(testPendingEvent2)
	assert.Equal(t, testPendingEvent1.EventName, q.Front().EventName)
}

func TestEventsQueue_Dequeue(t *testing.T) {
	q := EventsQueue{}
	assert.Nil(t, q.Dequeue())

	q.Enqueue(testPendingEvent1)
	q.Enqueue(testPendingEvent2)
	q.Enqueue(testPendingEvent3)
	assert.Equal(t, 3, q.size)

	event := q.Dequeue()
	assert.Equal(t, testPendingEvent1.EventName, event.EventName)
	assert.Equal(t, testPendingEvent2.EventName, q.front.EventName)
	assert.Equal(t, testPendingEvent3.EventName, q.back.EventName)
	assert.Equal(t, 1, len(q.shifted))
	assert.Equal(t, 2, q.size)

	event = q.Dequeue()
	assert.Equal(t, testPendingEvent2.EventName, event.EventName)
	assert.Equal(t, testPendingEvent3.EventName, q.front.EventName)
	assert.Equal(t, testPendingEvent3.EventName, q.back.EventName)
	assert.Equal(t, 2, len(q.shifted))
	assert.Equal(t, 1, q.size)

	event = q.Dequeue()
	assert.Equal(t, testPendingEvent3.EventName, event.EventName)
	assert.Nil(t, q.front)
	assert.Nil(t, q.back)
	assert.Equal(t, 3, len(q.shifted))
	assert.Equal(t, 0, q.size)
}

func TestEventsQueue_GetShiftedIds(t *testing.T) {
	ids := []uuid.UUID{
		uuid.MustParse("3c197784-0acc-11ef-a889-8389c5594477"),
		uuid.MustParse("3c197784-0acc-11ef-a889-8389c5594478"),
		uuid.MustParse("3c197784-0acc-11ef-a889-8389c5594479"),
	}

	q := EventsQueue{}
	q.shifted = map[uuid.UUID]*PendingEvent{
		ids[0]: {
			Id:        ids[0],
			EventName: "test",
		},
		ids[1]: {
			Id:        ids[1],
			EventName: "test",
		},
		ids[2]: {
			Id:        ids[2],
			EventName: "test",
		},
	}

	actual := q.GetShiftedIds()
	reflect.DeepEqual(ids, actual)
}

func TestEventsQueue_GetNewEvents(t *testing.T) {
	q := EventsQueue{}
	q.Enqueue(testPendingEvent1)
	q.Enqueue(testPendingEvent2)
	q.Enqueue(testPendingEvent3)

	result := q.GetNewEvents()
	assert.Equal(t, 2, len(result))
	assert.Equal(t, testPendingEvent2.EventName, result[0].EventName)
	assert.Equal(t, testPendingEvent3.EventName, result[1].EventName)
}
