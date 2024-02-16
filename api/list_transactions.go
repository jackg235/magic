package api

import (
	"context"
	"encoding/json"
	"magic-ledger/ledger"
	"magic-ledger/logger"
	"net/http"
)

type ListTransactionsResponse struct {
	Transactions interface{} `json:"transactions"`
}

func ListTransactions(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	transactions, err := ledger.ListTransactions(ctx)
	if err != nil {
		http.Error(w, "error listing ledger account", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(
		ListTransactionsResponse{
			Transactions: transactions,
		},
	)
	if err != nil {
		logger.Error(ctx, err, "error encoding response")
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
