package api

import (
	"context"
	"log"
	"magic-ledger/ledger"
	"magic-ledger/logger"
)

func InitializeInternalAccounts() {
	ctx := context.Background()
	assetsAccount, err := ledger.GetAccount(ctx, assetsAccountName)
	if err != nil {
		log.Fatal(err)
	}
	if assetsAccount != nil && len(assetsAccount.Metadata) != 0 {
		logger.Info(ctx, "internal accounts already created")
		return
	}
	// initialize assets, revenue, expenses accounts. API does not provide a convenient way
	// to create an account so we need to create a NOOP transaction
	metadata := map[string]interface{}{
		transactionTypeKey: createInternalAccountsTransaction,
	}
	postings := []ledger.TransactionPosting{
		{
			Src:    worldAccountName,
			Dest:   assetsAccountName,
			Amount: 0,
		},
		{
			Src:    worldAccountName,
			Dest:   revenueAccountName,
			Amount: 0,
		},
		{
			Src:    worldAccountName,
			Dest:   expensesAccountName,
			Amount: 0,
		},
	}
	_, err = ledger.CreateTransactionWithPostings(ctx, metadata, postings)
	if err != nil {
		log.Fatal(err)
	}
	assetsMetadata := map[string]interface{}{
		balanceTypeKey:    balanceTypeDebit,
		ledgerableTypeKey: ledgerableTypeInternal,
	}
	if err = ledger.AddMetaDataToAccount(ctx, assetsAccountName, assetsMetadata); err != nil {
		log.Fatal(err)
	}

	revenueMetadata := map[string]interface{}{
		balanceTypeKey:    balanceTypeCredit,
		ledgerableTypeKey: ledgerableTypeInternal,
	}
	if err = ledger.AddMetaDataToAccount(ctx, revenueAccountName, revenueMetadata); err != nil {
		log.Fatal(err)
	}

	expensesMetadata := map[string]interface{}{
		balanceTypeKey:    balanceTypeDebit,
		ledgerableTypeKey: ledgerableTypeInternal,
	}
	if err = ledger.AddMetaDataToAccount(ctx, expensesAccountName, expensesMetadata); err != nil {
		log.Fatal(err)
	}
	logger.Info(ctx, "successfully created internal accounts")
}
