package api

import (
	"context"
	"encoding/json"
	"fmt"
	"magic-ledger/ledger"
	"magic-ledger/logger"
	"net/http"
)

type PayoutMerchantRequest struct {
	MerchantId *string `json:"merchant_id"`

	// the amount to payout
	Amount *int64 `json:"amount,string"`
}

type PayoutMerchantResponse struct {
	Transaction interface{} `json:"transaction"`
}

func PayoutMerchant(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	decoder := json.NewDecoder(r.Body)
	var req PayoutMerchantRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "unable to decode request", http.StatusBadRequest)
		return
	}

	if req.MerchantId == nil || req.Amount == nil {
		http.Error(w, "merchantId and amount cannot be null", http.StatusBadRequest)
		return
	}
	account, err := ledger.GetAccount(ctx, *req.MerchantId)
	if err != nil {
		http.Error(w, "error getting ledger account", http.StatusInternalServerError)
		return
	}
	if account == nil {
		http.Error(w, fmt.Sprintf("no ledger account associated with address %s", *req.MerchantId), http.StatusBadRequest)
		return
	}

	metadata := map[string]interface{}{
		transactionTypeKey: payoutMerchantTransaction,
		merchantIdKey:      *req.MerchantId,
	}
	postings := []ledger.TransactionPosting{
		{
			Src:    *req.MerchantId,
			Dest:   worldAccountName,
			Amount: *req.Amount,
		},
		{
			Src:    assetsAccountName,
			Dest:   worldAccountName,
			Amount: *req.Amount,
		},
	}
	txn, err := ledger.CreateTransactionWithPostings(ctx, metadata, postings)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating transaction: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(
		PayoutMerchantResponse{
			Transaction: txn,
		},
	)
	if err != nil {
		logger.Error(ctx, err, "error encoding response")
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
