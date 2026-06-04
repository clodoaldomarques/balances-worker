package mysqldb

import (
	"time"

	"github.com/clodoaldomarques/balances-api/internal/shared/commons"
	"github.com/clodoaldomarques/balances-api/internal/worker/domain/daily"
)

type Balance struct {
	Date      time.Time          `json:"date"`
	AccountID int64              `json:"account_id"`
	OrgID     string             `json:"org_id"`
	Balances  commons.DecimalMap `json:"balances"`
	Version   int64              `json:"version"`
}

func (b Balance) toEntity() daily.Balance {
	return daily.Balance{
		Date:      b.Date,
		AccountID: b.AccountID,
		OrgID:     b.OrgID,
		Balances:  b.Balances,
		Version:   b.Version,
	}
}

func buildBalanceTable(b daily.Balance) Balance {
	return Balance{
		Date:      b.Date,
		AccountID: b.AccountID,
		OrgID:     b.OrgID,
		Balances:  b.Balances,
		Version:   b.Version,
	}
}
