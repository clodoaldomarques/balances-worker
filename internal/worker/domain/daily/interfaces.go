package daily

import (
	"context"
	"time"
)

//go:generate mockgen -source=interfaces.go -destination=mock.go -package=daily
type Repository interface {
	SaveNewBalance(ctx context.Context, b Balance) error
	UpdateExistingBalance(ctx context.Context, b Balance) error
	RetrieveLastBalance(ctx context.Context, accountID int64, orgID string) (Balance, error)
	RetrieveBalanceByPeriod(ctx context.Context, accountID int64, orgID string, initialDate, finalDate time.Time) ([]Balance, error)
}
