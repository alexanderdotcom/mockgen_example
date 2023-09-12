package service

import (
	"context"
	"fmt"
)

//go:generate go run go.uber.org/mock/mockgen -destination mocks/gomock/service.go github.com/alexanderdotcom/mockgen_example/pkg/service Bank
//go:generate go run github.com/matryer/moq -pkg moq_serivce -out ./mocks/moq/serivce.go . Bank

type Bank interface {
	CreateAccount(ctx context.Context, account string) error
	DepositMoney(ctx context.Context, account string, amount int) error
	WithdrawMoney(ctx context.Context, account string, amount int) error
}

func TransferMoney(ctx context.Context, ba Bank, fromAccount string, toAccount string, amount int) error {
	if err := ba.WithdrawMoney(ctx, fromAccount, amount); err != nil {
		return fmt.Errorf("unable to withdraw money: %w", err)
	}
	if err := ba.DepositMoney(ctx, toAccount, amount); err != nil {
		return fmt.Errorf("unable to deposit money: %w", err)
	}
	return nil
}

func DepositBeforeWithdrawTransferMoney(ctx context.Context, ba Bank, fromAccount string, toAccount string, amount int) error {
	if err := ba.DepositMoney(ctx, toAccount, amount); err != nil {
		return fmt.Errorf("unable to deposit money: %w", err)
	}
	if err := ba.WithdrawMoney(ctx, fromAccount, amount); err != nil {
		return fmt.Errorf("unable to withdraw money: %w", err)
	}
	return nil
}

func DepositHappensTwiceTransferMoney(ctx context.Context, ba Bank, fromAccount string, toAccount string, amount int) error {
	if err := ba.WithdrawMoney(ctx, fromAccount, amount); err != nil {
		return fmt.Errorf("unable to withdraw money: %w", err)
	}
	if err := ba.DepositMoney(ctx, toAccount, amount); err != nil {
		return fmt.Errorf("unable to deposit money: %w", err)
	}
	if err := ba.DepositMoney(ctx, toAccount, amount); err != nil {
		return fmt.Errorf("unable to deposit money: %w", err)
	}
	return nil
}

func DepositNeverHappensTransferMoney(ctx context.Context, ba Bank, fromAccount string, toAccount string, amount int) error {
	if err := ba.WithdrawMoney(ctx, fromAccount, amount); err != nil {
		return fmt.Errorf("unable to withdraw money: %w", err)
	}
	return nil
}
