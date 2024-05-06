package domain

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var (
	testPushEventDto1 = &EventPushDto{
		EventName: "test",
		Params:    "{}",
		External:  false,
		FromDb:    true,
	}
	testPushEventDto2 = &EventPushDto{
		EventName: "test2",
		Params:    "{}",
		External:  false,
		FromDb:    false,
	}
	testPushEventDto3 = &EventPushDto{
		EventName: "test3",
		Params:    "{}",
		External:  false,
		FromDb:    false,
	}
)

func TestEventsQueue_Enqueue(t *testing.T) {
	q := EventsQueue{}

	err := q.Enqueue(testPushEventDto1)
	assert.NoError(t, err)
	assert.Equal(t, testPushEventDto1.EventName, q.front.EventName)
	assert.Equal(t, testPushEventDto1.EventName, q.back.EventName)
	assert.Equal(t, 1, q.size)

	err = q.Enqueue(testPushEventDto2)
	assert.NoError(t, err)
	assert.Equal(t, testPushEventDto1.EventName, q.front.EventName)
	assert.Equal(t, testPushEventDto2.EventName, q.back.EventName)
	assert.Equal(t, 2, q.size)
}

func TestEventsQueue_PushToFront(t *testing.T) {
	q := EventsQueue{}
	_ = q.Enqueue(testPushEventDto1)
	_ = q.Enqueue(testPushEventDto2)

	assert.NoError(t, q.PushToFront(testPushEventDto3))
	assert.Equal(t, testPushEventDto3.EventName, q.front.EventName)
	assert.Equal(t, 3, q.size)
}

func TestEventsQueue_Front(t *testing.T) {
	q := EventsQueue{}
	assert.Nil(t, q.Front())

	_ = q.Enqueue(testPushEventDto1)
	assert.Equal(t, testPushEventDto1.EventName, q.Front().EventName)

	_ = q.Enqueue(testPushEventDto2)
	assert.Equal(t, testPushEventDto1.EventName, q.Front().EventName)
}

func TestEventsQueue_Dequeue(t *testing.T) {
	q := EventsQueue{}
	assert.Nil(t, q.Dequeue())

	_ = q.Enqueue(testPushEventDto1)
	_ = q.Enqueue(testPushEventDto2)
	_ = q.Enqueue(testPushEventDto3)
	assert.Equal(t, 3, q.size)

	event := q.Dequeue()
	assert.Equal(t, testPushEventDto1.EventName, event.EventName)
	assert.Equal(t, testPushEventDto2.EventName, q.front.EventName)
	assert.Equal(t, 1, len(q.shifted))
	assert.Equal(t, 2, q.size)

	event = q.Dequeue()
	assert.Equal(t, testPushEventDto2.EventName, event.EventName)
	assert.Equal(t, testPushEventDto3.EventName, q.front.EventName)
	assert.Equal(t, 2, len(q.shifted))
	assert.Equal(t, 1, q.size)

	event = q.Dequeue()
	assert.Equal(t, testPushEventDto3.EventName, event.EventName)
	assert.Nil(t, q.front)
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
	_ = q.Enqueue(testPushEventDto1)
	_ = q.Enqueue(testPushEventDto2)
	_ = q.Enqueue(testPushEventDto3)

	result := q.GetNewEvents()
	assert.Equal(t, 2, len(result))
	assert.Equal(t, testPushEventDto2.EventName, result[0].EventName)
	assert.Equal(t, testPushEventDto3.EventName, result[1].EventName)
}
