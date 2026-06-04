package consumer

import (
	"context"

	"github.com/clodoaldomarques/balances-api/internal/shared/domain/events"
)

//go:generate mockgen -source=interfaces.go -destination=mock.go -package=consumer
type Queue interface {
	Retrieve(ctx context.Context) (map[string]events.Event, error)
	DeleteMessages(receiptHandles []string) error
}
