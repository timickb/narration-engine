package worker

import (
	"context"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

// Worker Контракт обработчика состояний.
type Worker interface {
	HandleState(ctx context.Context, req *schema.HandleRequest) (*schema.HandleResponse, error)
	Name() string
}
