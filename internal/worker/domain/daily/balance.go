package daily

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	AvailableBalance = "available_balance"
	AvailableSavings = "available_savings"
	Holdfunds        = "hold_funds"
)

type Balance struct {
	Date      time.Time
	AccountID int64
	OrgID     string
	Balances  map[string]decimal.Decimal
	Version   int64
}
