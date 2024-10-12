package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">>> before:", account1.Balance, account2.Balance)

	// run a concurrent transfer transactions
	maxRetries := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for retries := 0; retries < maxRetries; retries++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// checking results
	existed := make(map[int]bool)
	for i := 0; i < maxRetries; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries

		// from entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// to entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		fmt.Println(">>> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount ..., maxRetries * amount

		k := int(diff1 / amount)

		require.True(t, k >= 1 && k <= maxRetries)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balances
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">>> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-(int64(maxRetries)*amount), updatedAccount1.Balance)
	require.Equal(t, account2.Balance+(int64(maxRetries)*amount), updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	/*
		Checking the deadlock transactions, for instance such scenario

		1 user wants to send money to the 2
		steps:
			- decrease account from 1
			- increase account to 2

		But we get the deadlock if this operation handles twice and reversed at the same time
		Then we get:
		steps:
			- decrease account from 1
			- decrease account from 2

			// HERE IS A DEADLOCK SINCE ONE QUERY IS WAITING FOR ANOTHER ONE AND VICE VERSA
			- increase account to 2
			- increase account to 1


		SOLUTION: Change order of the second reversed query like that:
		steps:
			- decrease amount from 1
			- increase account to 1

			- increase amount to 2
			- decrease amount from 2

	*/
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">>> before:", account1.Balance, account2.Balance)

	// run a concurrent transfer transactions
	maxRetries := 10
	amount := int64(10)
	errs := make(chan error)

	for retries := 0; retries < maxRetries; retries++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if retries%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}
	// check errors

	for i := 0; i < maxRetries; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	// checking final balance equality
	fmt.Println(">>> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
