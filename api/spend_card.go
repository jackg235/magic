package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"magic-ledger/ledger"
	"magic-ledger/logger"
	"net/http"
	"strings"
)

type SpendCardRequest struct {
	CardAddress *string `json:"card_address"`

	// the amount spent
	Amount *int64 `json:"amount,string"`
}

type SpendCardResponse struct {
	Transaction interface{} `json:"transaction"`
}

func SpendCard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	decoder := json.NewDecoder(r.Body)
	var req SpendCardRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "unable to decode SpendCard request", http.StatusBadRequest)
		return
	}
	logger.Info(ctx, "got SpendCard request %v", req)

	if req.CardAddress == nil || req.Amount == nil {
		http.Error(w, "cardAddress and amount cannot be null", http.StatusBadRequest)
		return
	}
	account, err := ledger.GetAccount(ctx, *req.CardAddress)
	if err != nil {
		http.Error(w, "error getting ledger account", http.StatusInternalServerError)
		return
	}
	if account == nil {
		http.Error(w, fmt.Sprintf("no ledger account associated with address %s", *req.CardAddress), http.StatusBadRequest)
		return
	}
	merchantId, ok := account.Metadata[merchantIdKey]
	if !ok {
		http.Error(w, fmt.Sprintf("no merchant id associated with account address: %s", *req.CardAddress), http.StatusBadRequest)
		return
	}
	userName, ok := account.Metadata[nameKey]
	if !ok {
		http.Error(w, fmt.Sprintf("no user id associated with account address: %s", *req.CardAddress), http.StatusBadRequest)
		return
	}
	purchaseId := fmt.Sprintf("purchase:%s", strings.Replace(uuid.NewString(), "-", "", -1))

	metadata := map[string]interface{}{
		transactionTypeKey: spendCardTransaction,
		cardIdKey:          *req.CardAddress,
		nameKey:            userName,
		merchantIdKey:      merchantId,
		purchaseIdKey:      purchaseId,
	}
	postings := []ledger.TransactionPosting{
		{
			Src:    *req.CardAddress,
			Dest:   fmt.Sprintf("%v", merchantId),
			Amount: *req.Amount,
		},
	}
	txn, err := ledger.CreateTransactionWithPostings(ctx, metadata, postings)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating transaction: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(
		SpendCardResponse{
			Transaction: txn,
		},
	)
	if err != nil {
		logger.Error(ctx, err, "error encoding response")
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
