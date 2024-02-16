package api

import (
	"context"
	"encoding/json"
	"fmt"
	"magic-ledger/ledger"
	"magic-ledger/logger"
	"net/http"
)

type ListAccountsResponse struct {
	Accounts interface{} `json:"accounts"`
}

type Account struct {
	Address        string `json:"address"`
	Name           string `json:"name"`
	MerchantId     string `json:"merchant_id"`
	Balance        int64  `json:"balance"`
	BalanceType    string `json:"balance_type"`
	LedgerableType string `json:"ledgerable_type"`
}

func ListAccounts(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	accounts, err := ledger.ListAccounts(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("error listing ledger accounts: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	balances, err := ledger.ListBalances(ctx)
	if err != nil {
		http.Error(w, "error listing ledger balances", http.StatusInternalServerError)
		return
	}
	accountsWithBalances := make([]Account, len(accounts))
	for i, acct := range accounts {
		accountWithBalance := Account{
			Address: acct.Address,
			Balance: balances[acct.Address],
		}
		if name, ok := acct.Metadata[nameKey]; ok {
			accountWithBalance.Name = fmt.Sprintf("%v", name)
		}
		if merchantId, ok := acct.Metadata[merchantIdKey]; ok {
			accountWithBalance.MerchantId = fmt.Sprintf("%v", merchantId)
		}
		if balanceType, ok := acct.Metadata[balanceTypeKey]; ok {
			accountWithBalance.BalanceType = fmt.Sprintf("%v", balanceType)
		}
		if ledgerableType, ok := acct.Metadata[ledgerableTypeKey]; ok {
			accountWithBalance.LedgerableType = fmt.Sprintf("%v", ledgerableType)
		}
		accountsWithBalances[i] = accountWithBalance
	}

	err = json.NewEncoder(w).Encode(
		ListAccountsResponse{
			Accounts: accountsWithBalances,
		},
	)
	if err != nil {
		logger.Error(ctx, err, "error encoding response")
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
