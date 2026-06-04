package daily

import (
	"context"
	"math"
	"time"
)

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{
		r: r,
	}
}

func (s Service) CreateNewBalance(ctx context.Context, b Balance) error {
	return s.r.SaveNewBalance(ctx, b)
}

func (s Service) UpdateExistingBalance(ctx context.Context, b Balance) error {
	old, err := s.r.RetrieveLastBalance(ctx, b.AccountID, b.OrgID)
	if err != nil {
		return err
	}

	if old.Date.Before(b.Date) {
		return s.FillMissingBalances(ctx, old, b)
	}

	return s.r.UpdateExistingBalance(ctx, b)
}

func (s Service) FillMissingBalances(ctx context.Context, old, new Balance) error {
	diff := new.Date.Sub(old.Date).Hours() / 24
	for i := 1; i <= int(math.Round(diff)); i++ {
		f := buildFillBalance(i, old, new)
		if err := s.r.SaveNewBalance(ctx, f); err != nil {
			return err
		}
	}
	return nil
}

func buildFillBalance(i int, old, new Balance) Balance {
	fdate := old.Date.AddDate(0, 0, i)
	if fdate.Before(new.Date) {
		return Balance{
			Date:      fdate,
			AccountID: old.AccountID,
			OrgID:     old.OrgID,
			Balances:  old.Balances,
			Version:   old.Version,
		}
	}
	return Balance{
		Date:      new.Date,
		AccountID: new.AccountID,
		OrgID:     new.OrgID,
		Balances:  new.Balances,
		Version:   new.Version,
	}
}

func (s Service) RetrieveBalanceByDate(ctx context.Context, accountID int64, orgID string, initialDate, finalDate time.Time) ([]Balance, error) {
	return nil, nil
}
