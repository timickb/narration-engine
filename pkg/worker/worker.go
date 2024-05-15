package worker

import (
	"context"
	"encoding/json"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

// Worker Контракт обработчика состояний.
type Worker interface {
	HandleState(ctx context.Context, req *schema.HandleRequest) (*schema.HandleResponse, error)
	Name() string
}

func UnmarshallRequestBody[T any](req *schema.HandleRequest) (*T, error) {
	var result T
	merged, err := jsonpatch.MergePatch([]byte(req.Context), []byte(req.EventParams))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(merged, &result); err != nil {
		return nil, fmt.Errorf("unmarshall request context: %w", err)
	}
	return &result, nil
}
