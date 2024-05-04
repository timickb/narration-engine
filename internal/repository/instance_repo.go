package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/internal/repository/models"
	"github.com/timickb/narration-engine/pkg/db"
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
	return nil
}

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
        	WHERE id = ? AND (locked_by = '' OR locked_by IS NULL OR locked_till < now())
        	RETURNING *
		)
		SELECT updated.*,
       			(SELECT json_agg(row_to_json(pe.*)) FROM pending_events pe
            	WHERE pe.instance_id = updated.id) as events
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
	if err := r.db.WithTxSupport(ctx).Take(&instance).Error; err != nil {
		return nil, fmt.Errorf("get instance by id: %w", err)
	}
	return &domain.Instance{}, nil
}
