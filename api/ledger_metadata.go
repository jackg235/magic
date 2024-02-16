package api

import (
	"context"
	"encoding/json"
	"fmt"
	"magic-ledger/ledger"
	"magic-ledger/logger"
	"net/http"
)

type LedgerMetadataResponse struct {
	Debits   int64 `json:"debits"`
	Credits  int64 `json:"credits"`
	Expenses int64 `json:"expenses"`
	Assets   int64 `json:"assets"`
	Revenue  int64 `json:"revenue"`
}

// LedgerMetadata serves as a sanity check that debits = credits. Also returns retained earnings info
func LedgerMetadata(w http.ResponseWriter, _ *http.Request) {
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
	var res LedgerMetadataResponse
	debits := int64(0)
	credits := int64(0)
	for _, acct := range accounts {
		acctBalance := balances[acct.Address]
		if acct.Address == assetsAccountName {
			res.Assets = acctBalance
		} else if acct.Address == revenueAccountName {
			res.Revenue = acctBalance
		} else if acct.Address == expensesAccountName {
			res.Expenses = acctBalance
		}
		if balanceType, ok := acct.Metadata[balanceTypeKey]; ok {
			if BalanceType(fmt.Sprintf("%v", balanceType)) == balanceTypeCredit {
				credits += acctBalance
			} else if BalanceType(fmt.Sprintf("%v", balanceType)) == balanceTypeDebit {
				debits += acctBalance
			}
		}
	}
	res.Credits = credits
	res.Debits = debits
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.Error(ctx, err, "error encoding response")
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
