package domain

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
)

// InstanceContext Структура, описывающая контекст экземпляра сценария.
type InstanceContext struct {
	data map[string]interface{}
}

func NewInstanceContext(data []byte) (*InstanceContext, error) {
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, fmt.Errorf("unmarshal instance context to map: %w", err)
	}
	return &InstanceContext{data: parsed}, nil
}

// MergeData Объединить данные контекста с новыми данными с приоритетом.
func (c *InstanceContext) MergeData(dataToMerge []byte) error {
	var parsed map[string]interface{}
	if err := json.Unmarshal(dataToMerge, &parsed); err != nil {
		return fmt.Errorf("unmarshal data to map: %w", err)
	}

	for key, value := range parsed {
		c.data[key] = value
	}
	return nil
}

// GetValue Получить данные из контекста с помощью пути в JSON. Пример: user_data.common.emails.
func (c *InstanceContext) GetValue(path string) (interface{}, error) {
	jsonStr, err := json.Marshal(c.data)
	if err != nil {
		return nil, err
	}
	return gjson.Get(string(jsonStr), path).Value(), nil
}

// SetRootValue Установить значение ключа на первом уровне.
func (c *InstanceContext) SetRootValue(key string, value interface{}) {
	c.data[key] = value
}

// String Получить JSON-строку контекста.
func (c *InstanceContext) String() string {
	res, _ := json.Marshal(c.data)
	return string(res)
}
