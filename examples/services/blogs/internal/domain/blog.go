package domain

import (
	"github.com/google/uuid"
	"time"
)

// Blog Сущность блога.
type Blog struct {
	Id                uuid.UUID
	AuthorId          uuid.UUID
	Name              string
	SubscribersCount  int
	PublicationsCount int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// BlogUpdateDto Структура для обновления информации о блоге.
type BlogUpdateDto struct {
	Id                uuid.UUID
	Name              *string
	SubscribersCount  *int
	PublicationsCount *int
}
