package daily

import (
	"time"

	"github.com/clodoaldomarques/balances-api/internal/shared/domain/events"

	"github.com/shopspring/decimal"
)

const (
	MaxLimit       = "max_limit"
	TotalLimit     = "total_limit"
	OverdraftLimit = "overdraft_limit"
)

const (
	Available = "available"
	Savings   = "savings"
	Blocked   = "blocked"
)

var (
	CREATE_ACCOUNT events.Type = "create_account"
	UPDATE_ACCOUNT events.Type = "update_account"
	PROCESS_ENTRY  events.Type = "process_entry"
	layout                     = "2006-01-02"
)

type CreateAccountEvent struct {
	AccountID int64                      `json:"account_id"`
	OrgID     string                     `json:"org_id"`
	Limits    map[string]decimal.Decimal `json:"limits"`
	Balances  map[string]decimal.Decimal `json:"balances"`
	CreatedAt time.Time                  `json:"created_at"`
	Status    string                     `json:"status"`
	Version   int64                      `json:"version"`
}

func (cae CreateAccountEvent) ToDailyBalance() Balance {
	return Balance{
		Date:      parseBalanceDate(cae.CreatedAt.Format(layout)),
		AccountID: cae.AccountID,
		OrgID:     cae.OrgID,
		Balances:  calculate(cae.Limits, cae.Balances),
		Version:   cae.Version,
	}
}

type UpdateAccountEvent struct {
	AccountID int64                      `json:"account_id"`
	OrgID     string                     `json:"org_id"`
	Limits    map[string]decimal.Decimal `json:"limits"`
	Balances  map[string]decimal.Decimal `json:"balances"`
	Status    string                     `json:"status"`
	UpdatedAt time.Time                  `json:"updated_at"`
	Version   int64                      `json:"version"`
}

func (uae UpdateAccountEvent) ToDailyBalance() Balance {
	return Balance{
		Date:      parseBalanceDate(uae.UpdatedAt.Format(layout)),
		AccountID: uae.AccountID,
		OrgID:     uae.OrgID,
		Balances:  calculate(uae.Limits, uae.Balances),
		Version:   uae.Version,
	}
}

type ProcessEntryEvent struct {
	AccountID  int64                      `json:"account_id"`
	OrgID      string                     `json:"org_id"`
	TrackingID string                     `json:"tracking_id"`
	Impacts    []ImpactEvent              `json:"impacts"`
	Limits     map[string]decimal.Decimal `json:"limits"`
	Balances   map[string]decimal.Decimal `json:"balances"`
	Version    int64                      `json:"version"`
	CreatedAt  time.Time                  `json:"created_at"`
}

func (pee ProcessEntryEvent) ToDailyBalance() Balance {
	return Balance{
		Date:      parseBalanceDate(pee.CreatedAt.Format(layout)),
		AccountID: pee.AccountID,
		OrgID:     pee.OrgID,
		Balances:  calculate(pee.Limits, pee.Balances),
		Version:   pee.Version,
	}
}

type ImpactEvent struct {
	Balance   string          `json:"balance"`
	Operation string          `json:"operation"`
	Amount    decimal.Decimal `json:"amount"`
	Rules     []string        `json:"rules,omitempty"`
}

func parseBalanceDate(d string) time.Time {
	dt, _ := time.Parse(layout, d)
	return dt.UTC()
}

func calculate(limits, balances map[string]decimal.Decimal) map[string]decimal.Decimal {
	calculated := make(map[string]decimal.Decimal)
	calculated[AvailableBalance] = limits[TotalLimit].Add(balances[Available])
	calculated[AvailableSavings] = balances[Savings]
	calculated[Holdfunds] = balances[Blocked]
	return calculated
}
