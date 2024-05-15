package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/internal/repository/models"
	"github.com/timickb/narration-engine/pkg/db"
	"github.com/timickb/narration-engine/pkg/utils"
	"gorm.io/gorm"
	"time"
)

type instanceRepo struct {
	db *db.Database
}

// NewInstanceRepo Создать репозиторий над экземплярами сценариев.
func NewInstanceRepo(db *db.Database) *instanceRepo {
	return &instanceRepo{db: db}
}

// Update Обновить экземпляр сценария.
func (r *instanceRepo) Update(ctx context.Context, instance *domain.Instance) error {
	err := r.db.WithContext(ctx).Debug().Transaction(func(tx *gorm.DB) error {
		// Обновить сам экземпляр.
		err := tx.Model(&models.DbInstance{}).
			Where("id = ?", instance.Id).
			Updates(models.NewDbInstanceUpdate(instance)).Error
		if err != nil {
			return err
		}
		// Обновить очередь событий.
		if err := r.updateEvents(tx, instance.PendingEvents, instance.Id); err != nil {
			return err
		}
		// До обновления могли добавиться новые события извне - нужно обновить очередь.
		queue, err := r.fetchNewPendingEvents(tx, instance.Id)
		if err != nil {
			return err
		}

		instance.PendingEvents = queue
		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return fmt.Errorf("update instance: %w", err)
	}
	return nil
}

// Create Создать новый экземпляр сценария.
func (r *instanceRepo) Create(ctx context.Context, dto *domain.CreateInstanceDto) (uuid.UUID, error) {
	instanceId, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil, fmt.Errorf("new uuid: %w", err)
	}
	dbInstance := models.NewDbInstance(instanceId, dto)

	if err := r.db.WithTxSupport(ctx).Create(dbInstance).Error; err != nil {
		return uuid.Nil, fmt.Errorf("create instance: %w", err)
	}
	return instanceId, nil
}

// FetchWithLock Достать экземпляр с блокировкой.
func (r *instanceRepo) FetchWithLock(ctx context.Context, dto *domain.FetchInstanceDto) (*domain.Instance, error) {
	instance := models.FilledInstance{}
	lockedTill := time.Now().Add(dto.LockTimeout)

	err := r.db.WithTxSupport(ctx).Debug().Raw(
		`WITH updated AS (
    	UPDATE instances SET locked_by = ?, locked_till = ?
        	WHERE id = ? AND failed = false AND (locked_by = '' OR locked_by IS NULL OR locked_till < now())
        	RETURNING *
		)
		SELECT updated.*,
       			(SELECT json_agg(row_to_json(pe.*)) FROM pending_events pe
            	WHERE pe.instance_id = updated.id AND pe.executed_at < now()) as events
		FROM instances i LEFT JOIN updated USING (id) WHERE i.id = ?`,
		dto.LockerId, lockedTill, dto.Id, dto.Id,
	).Scan(&instance).Error

	if err != nil {
		return nil, fmt.Errorf("fetch instance with lock: %w", err)
	}

	result, err := instance.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("map instance to domain: %w", err)
	}
	return result, nil
}

// Unlock Снять блокировку с экземпляра.
func (r *instanceRepo) Unlock(ctx context.Context, id uuid.UUID) error {
	query := r.db.WithTxSupport(ctx).Model(&models.DbInstance{}).
		Where("id = ?", id).Updates(map[string]interface{}{
		"locked_by":   "",
		"locked_till": nil,
	})
	if query.Error != nil {
		return fmt.Errorf("unlock instance: %w", query.Error)
	}
	if query.RowsAffected == 0 {
		return fmt.Errorf("instance with id %s not found", id.String())
	}
	return nil
}

// GetWaitingIds Получить список ид экземпляров, ожидающих дальнейшей работы.
func (r *instanceRepo) GetWaitingIds(ctx context.Context, limit int) ([]uuid.UUID, error) {
	var ids []uuid.UUID

	err := r.db.WithTxSupport(ctx).Debug().Table("instances").
		Where("locked_by IS NULL OR locked_by = '' OR locked_till < now()").
		Where("start_after IS NULL OR start_after < now()").
		Where("failed = false").
		Where("current_state != ?", domain.StateEnd.Name).
		Select("id").
		Limit(limit).
		Scan(&ids).Error
	if err != nil {
		return nil, fmt.Errorf("get waiting instances ids: %w", err)
	}
	return ids, nil
}

// GetById Получить данные экземпляра по идентификатору.
func (r *instanceRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Instance, error) {
	var instance models.DbInstance
	if err := r.db.WithTxSupport(ctx).Where("id = ?", id).Take(&instance).Error; err != nil {
		return nil, fmt.Errorf("get instance by id: %w", err)
	}
	mappedInstance, err := instance.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("get instance by id: %w", err)
	}
	return mappedInstance, nil
}

// IsKeyBlocked Проверить, присутствует ли блокировка ключа в каком-нибудь экземпляре.
func (r *instanceRepo) IsKeyBlocked(ctx context.Context, key string) (bool, error) {
	var count int64
	err := r.db.WithTxSupport(ctx).Model(&models.DbInstance{}).
		Where("blocking_key = ?", key).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("get instance by blocking key: %w", err)
	}
	return count > 0, nil
}

func (r *instanceRepo) updateEvents(
	tx *gorm.DB,
	queue *domain.EventsQueue,
	instanceId uuid.UUID,
) error {

	for _, deletedEventId := range queue.GetShiftedIds() {
		err := tx.Delete(&models.DbPendingEvent{}, deletedEventId).Error
		if err != nil {
			return err
		}
	}

	newEvents := utils.MapSlice(
		queue.GetNewEvents(),
		func(e *domain.PendingEvent) *models.DbPendingEvent {
			return models.NewDbPendingEventFromDomain(e, instanceId)
		},
	)

	for _, event := range newEvents {
		if err := tx.Model(&models.DbPendingEvent{}).Create(&event).Error; err != nil {
			return fmt.Errorf("create pending event in db: %w", err)
		}
	}

	return nil
}

func (r *instanceRepo) fetchNewPendingEvents(
	tx *gorm.DB,
	instanceId uuid.UUID,
) (*domain.EventsQueue, error) {

	var newPendingEvents *models.DbPendingEvents
	err := tx.Model(&models.DbPendingEvent{}).
		Where("instance_id = ?", instanceId).
		Find(&newPendingEvents).
		Order("executed_at ASC").
		Error
	if err != nil {
		return nil, err
	}
	pendingEvents := newPendingEvents.ToDomain()
	queue := &domain.EventsQueue{}
	for _, event := range pendingEvents {
		event.FromDb = true
		queue.Enqueue(event)
	}

	return queue, nil
}
