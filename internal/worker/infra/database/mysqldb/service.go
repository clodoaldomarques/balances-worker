package mysqldb

import (
	"context"
	"database/sql"
	"time"

	"github.com/clodoaldomarques/balances-api/internal/worker/domain/daily"
	"github.com/clodoaldomarques/core-sdk/pkg/logger"
)

var (
	INSERT_BALANCE            = `insert into daily_balances (date, account_id, org_id, balances, version) values (?, ?, ?, ?, ?)`
	UPDATE_BALANCE            = `update daily_balances set balances = ?, version = ? where date = ? and  account_id = ? and org_id = ? and version < ?`
	SELECT_LAST_BALANCE       = `select db.date, db.account_id, db.org_id, db.balances, db.version from daily_balances db WHERE db.account_id = ? and db.org_id = ? order by db.date DESC limit 1`
	SELECT_BALANCES_BY_PERIOD = `select db.date, db.account_id, db.org_id, db.balances, db.version from daily_balances db where db.account_id = ? and db.org_id = ? and db.date BETWEEN ? and ? order by db.date`
)

type Repository struct {
	ctx context.Context
	db  *sql.DB
}

func NewRepository(ctx context.Context) *Repository {
	db, err := Connect()
	if err != nil {
		logger.Error(ctx, "error on connect to database", logger.Fields{"error": err.Error()})
		return &Repository{ctx: ctx}
	}

	return &Repository{ctx: ctx, db: db}
}

func (r Repository) Close() {
	r.db.Close()
}

func (r Repository) SaveNewBalance(ctx context.Context, b daily.Balance) error {
	bt := buildBalanceTable(b)
	statement, err := r.db.Prepare(INSERT_BALANCE)
	if err != nil {
		logger.Error(ctx, "error on save new balance", logger.Fields{"balance": b, "error": err.Error(), "sql": INSERT_BALANCE})
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(
		bt.Date,
		bt.AccountID,
		bt.OrgID,
		bt.Balances,
		bt.Version,
	)
	if err != nil {
		return err
	}
	return nil
}
func (r Repository) UpdateExistingBalance(ctx context.Context, b daily.Balance) error {
	bt := buildBalanceTable(b)
	statement, err := r.db.Prepare(UPDATE_BALANCE)
	if err != nil {
		logger.Error(ctx, "error on update existing balance", logger.Fields{"balance": b, "error": err.Error(), "sql": UPDATE_BALANCE})
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(
		bt.Balances,
		bt.Version,
		bt.Date,
		bt.AccountID,
		bt.OrgID,
		bt.Version,
	)
	if err != nil {
		return err
	}
	return nil
}
func (r Repository) RetrieveLastBalance(ctx context.Context, accountID int64, orgID string) (daily.Balance, error) {
	var b Balance
	err := r.db.QueryRow(SELECT_LAST_BALANCE, accountID, orgID).Scan(
		&b.Date,
		&b.AccountID,
		&b.OrgID,
		&b.Balances,
		&b.Version,
	)
	if err != nil {
		logger.Error(ctx, "error on retrieve existing balance", logger.Fields{"account": accountID, "error": err.Error(), "sql": SELECT_LAST_BALANCE})
		return daily.Balance{}, err
	}
	return b.toEntity(), nil

}

func (r Repository) RetrieveBalanceByPeriod(ctx context.Context, accountID int64, orgID string, initialDate, finalDate time.Time) ([]daily.Balance, error) {
	return nil, nil
}
