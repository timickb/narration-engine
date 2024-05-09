package worker

import (
	"context"
	"encoding/json"
	"fmt"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

// Worker Контракт обработчика состояний.
type Worker interface {
	HandleState(ctx context.Context, req *schema.HandleRequest) (*schema.HandleResponse, error)
	Name() string
}

func UnmarshallRequestBody[T any](req *schema.HandleRequest) (*T, error) {
	var result T
	if err := json.Unmarshal([]byte(req.Context), &result); err != nil {
		return nil, fmt.Errorf("unmarshall request context: %w", err)
	}
	return &result, nil
}
