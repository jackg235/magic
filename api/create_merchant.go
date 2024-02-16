package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"magic-ledger/ledger"
	"net/http"
	"strings"
)

type CreateMerchantRequest struct {
	MerchantName *string `json:"merchant_name"`
}

func CreateMerchant(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	decoder := json.NewDecoder(r.Body)
	var req CreateMerchantRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "unable to decode CreateMerchant request", http.StatusBadRequest)
		return
	}

	if req.MerchantName == nil {
		http.Error(w, "merchantName cannot be null", http.StatusBadRequest)
		return
	}

	merchantId := fmt.Sprintf("merchant:%s", strings.Replace(uuid.NewString(), "-", "", -1))
	metadata := map[string]interface{}{
		transactionTypeKey: createMerchantTransaction,
		merchantIdKey:      merchantId,
	}
	postings := []ledger.TransactionPosting{
		{
			Src:    worldAccountName,
			Dest:   merchantId,
			Amount: 0,
		},
	}
	_, err = ledger.CreateTransactionWithPostings(ctx, metadata, postings)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating transaction: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// add metadata to the account we just created
	accountMetadata := map[string]interface{}{
		nameKey:           *req.MerchantName,
		balanceTypeKey:    balanceTypeCredit,
		ledgerableTypeKey: ledgerableTypeExternal,
	}
	err = ledger.AddMetaDataToAccount(ctx, merchantId, accountMetadata)
	if err != nil {
		http.Error(w, fmt.Sprintf("error adding metadata to account %s", err.Error()), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
