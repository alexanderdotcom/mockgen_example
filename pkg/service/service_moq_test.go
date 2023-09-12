package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/alexanderdotcom/mockgen_example/pkg/service"
	moq_serivce "github.com/alexanderdotcom/mockgen_example/pkg/service/mocks/moq"
)

func TestBankAccountMoq(t *testing.T) {
	tests := []struct {
		name        string
		toAddress   string
		fromAddress string
		bank        func(t *testing.T) service.Bank
		amount      int
		wantErr     bool
	}{
		{
			"Success - Transferring money, each method can only be called once",
			"123ABC",
			"456DEF",
			func(t *testing.T) service.Bank {
				var withdrawTimes int
				var depositTimes int

				mockBank := &moq_serivce.BankMock{
					WithdrawMoneyFunc: func(ctx context.Context, account string, amount int) error {
						withdrawTimes += 1
						require.Equal(t, 1, withdrawTimes)
						return nil
					},
					DepositMoneyFunc: func(ctx context.Context, account string, amount int) error {
						depositTimes += 1
						require.Equal(t, 1, depositTimes)
						return nil
					},
				}

				return mockBank
			},
			100,
			false,
		},
		{
			"Success - Transferring money, withdraw needs to happen before deposit",
			"123ABC",
			"456DEF",
			func(t *testing.T) service.Bank {

				var withdrawTimes int
				var depositTimes int

				mockBank := &moq_serivce.BankMock{
					WithdrawMoneyFunc: func(ctx context.Context, account string, amount int) error {
						withdrawTimes += 1
						require.Equal(t, 1, withdrawTimes)
						return nil
					},
					DepositMoneyFunc: func(ctx context.Context, account string, amount int) error {
						depositTimes += 1
						require.Equal(t, 1, withdrawTimes)
						require.Equal(t, 1, depositTimes)
						return nil
					},
				}
				return mockBank
			},
			100,
			false,
		},
		{
			"Fail - Insufficient funds, return an error",
			"123ABC",
			"456DEF",
			func(t *testing.T) service.Bank {
				mockBank := &moq_serivce.BankMock{
					WithdrawMoneyFunc: func(ctx context.Context, account string, amount int) error {
						return fmt.Errorf("insufficient funds")
					},
					DepositMoneyFunc: func(ctx context.Context, account string, amount int) error {
						require.Fail(t, "We should not call this method")
						return nil
					},
				}
				return mockBank

			},
			100,
			true,
		},
		{
			"Success - Transferring money, checks that the argument values to method are correct",
			"123ABC",
			"456DEF",
			func(t *testing.T) service.Bank {
				mockBank := &moq_serivce.BankMock{
					WithdrawMoneyFunc: func(ctx context.Context, account string, amount int) error {
						require.Equal(t, account, "456DEF")
						return nil
					},
					DepositMoneyFunc: func(ctx context.Context, account string, amount int) error {
						return nil
					},
				}
				return mockBank
			},
			100,
			false,
		},
	}

	var skipBrokenTest = true

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()

			err := service.TransferMoney(ctx, tt.bank(t), tt.fromAddress, tt.fromAddress, tt.amount)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()

			err := service.DepositNeverHappensTransferMoney(ctx, tt.bank(t), tt.fromAddress, tt.fromAddress, tt.amount)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}

	for _, tt := range tests {
		if skipBrokenTest {
			t.Skip()
		}
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := context.TODO()

			err := service.DepositHappensTwiceTransferMoney(ctx, tt.bank(t), tt.fromAddress, tt.fromAddress, tt.amount)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}

	for _, tt := range tests {
		if skipBrokenTest {
			t.Skip()
		}
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := context.TODO()

			err := service.DepositBeforeWithdrawTransferMoney(ctx, tt.bank(t), tt.fromAddress, tt.fromAddress, tt.amount)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}

}
