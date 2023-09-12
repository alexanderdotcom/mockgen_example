package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/alexanderdotcom/mockgen_example/pkg/service"
	mock_service "github.com/alexanderdotcom/mockgen_example/pkg/service/mocks/gomock"
)

func TestBankAccount(t *testing.T) {
	tests := []struct {
		name        string
		toAddress   string
		fromAddress string
		bank        func(ctrl *gomock.Controller, t *testing.T) service.Bank
		amount      int
		wantErr     bool
	}{
		{
			"Success - Transferring money, each method can only be called once",
			"123ABC",
			"456DEF",
			func(ctrl *gomock.Controller, t *testing.T) service.Bank {
				mockBank := mock_service.NewMockBank(ctrl)
				b1 := mockBank.EXPECT().WithdrawMoney(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				mockBank.EXPECT().DepositMoney(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil).After(b1)
				return mockBank
			},
			100,
			false,
		},
		{
			"Success - Transferring money, withdraw needs to happen before deposit",
			"123ABC",
			"456DEF",
			func(ctrl *gomock.Controller, t *testing.T) service.Bank {
				mockBank := mock_service.NewMockBank(ctrl)
				b1 := mockBank.EXPECT().WithdrawMoney(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				mockBank.EXPECT().DepositMoney(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil).After(b1)
				return mockBank
			},
			100,
			false,
		},
		{
			"Fail - Insufficient funds, return an error",
			"123ABC",
			"456DEF",
			func(ctrl *gomock.Controller, t *testing.T) service.Bank {
				mockBank := mock_service.NewMockBank(ctrl)
				mockBank.EXPECT().WithdrawMoney(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("insufficient funds"))
				return mockBank
			},
			100,
			true,
		},
		{
			"Success - Transferring money, checks that the argument values to method are correct",
			"123ABC",
			"456DEF",
			func(ctrl *gomock.Controller, t *testing.T) service.Bank {
				mockBank := mock_service.NewMockBank(ctrl)
				mockBank.EXPECT().WithdrawMoney(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				mockBank.EXPECT().DepositMoney(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).DoAndReturn(
					func(ctx context.Context, account string, amount int) error {
						require.Equal(t, account, "456DEF")
						return nil
					})
				return mockBank
			},
			100,
			false,
		},
	}

	var skipBrokenTests = true

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := context.TODO()

			err := service.TransferMoney(ctx, tt.bank(ctrl, t), tt.fromAddress, tt.fromAddress, tt.amount)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if skipBrokenTests {
				t.Skip()
			}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := context.TODO()

			err := service.DepositHappensTwiceTransferMoney(ctx, tt.bank(ctrl, t), tt.fromAddress, tt.fromAddress, tt.amount)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if skipBrokenTests {
				t.Skip()
			}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := context.TODO()

			err := service.DepositBeforeWithdrawTransferMoney(ctx, tt.bank(ctrl, t), tt.fromAddress, tt.fromAddress, tt.amount)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}

}
