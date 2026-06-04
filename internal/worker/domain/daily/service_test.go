package daily

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestService_UpdateExistingBalance(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctrl *gomock.Controller) *Service
		fields func() Balance
		want   func(t *testing.T, err error)
	}{
		{
			name: "when update existing balance without error",
			setup: func(ctrl *gomock.Controller) *Service {
				m := NewMockRepository(ctrl)
				m.EXPECT().RetrieveLastBalance(gomock.Any(), gomock.Eq(int64(1234)), gomock.Eq("TN-Test")).Return(buildFakeBalance("2024-12-01"), nil).Times(1)
				m.EXPECT().UpdateExistingBalance(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				return NewService(m)
			},
			fields: func() Balance {
				return buildFakeBalance("2024-12-01")
			},
			want: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "when fill latest balances",
			setup: func(ctrl *gomock.Controller) *Service {
				m := NewMockRepository(ctrl)
				m.EXPECT().RetrieveLastBalance(gomock.Any(), gomock.Eq(int64(1234)), gomock.Eq("TN-Test")).Return(buildFakeBalance("2024-12-01"), nil).Times(1)
				m.EXPECT().SaveNewBalance(gomock.Any(), gomock.Any()).Return(nil).Times(9)
				return NewService(m)
			},
			fields: func() Balance {
				return buildFakeBalance("2024-12-10")
			},
			want: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "when update existing balance with retrieve error",
			setup: func(ctrl *gomock.Controller) *Service {
				m := NewMockRepository(ctrl)
				m.EXPECT().RetrieveLastBalance(gomock.Any(), gomock.Eq(int64(1234)), gomock.Eq("TN-Test")).Return(Balance{}, errors.New("not found")).Times(1)
				return NewService(m)
			},
			fields: func() Balance {
				return buildFakeBalance("2024-12-01")
			},
			want: func(t *testing.T, err error) {
				assert.Equal(t, "not found", err.Error())
			},
		},
		{
			name: "when update existing balance with update error",
			setup: func(ctrl *gomock.Controller) *Service {
				m := NewMockRepository(ctrl)
				m.EXPECT().RetrieveLastBalance(gomock.Any(), gomock.Eq(int64(1234)), gomock.Eq("TN-Test")).Return(buildFakeBalance("2024-12-01"), nil).Times(1)
				m.EXPECT().UpdateExistingBalance(gomock.Any(), gomock.Any()).Return(errors.New("fail")).Times(1)
				return NewService(m)
			},
			fields: func() Balance {
				return buildFakeBalance("2024-12-01")
			},
			want: func(t *testing.T, err error) {
				assert.Equal(t, "fail", err.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := tt.setup(ctrl)
			err := s.UpdateExistingBalance(context.Background(), tt.fields())
			tt.want(t, err)
		})
	}
}

func buildFakeBalance(d string) Balance {
	date, _ := time.Parse("2006-01-02", d)
	return Balance{
		Date:      date,
		AccountID: 1234,
		OrgID:     "TN-Test",
		Balances: map[string]decimal.Decimal{
			"Saldo": decimal.NewFromInt(10),
		},
		Version: 1,
	}
}
