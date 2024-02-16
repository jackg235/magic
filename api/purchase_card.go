package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"magic-ledger/ledger"
	"magic-ledger/logger"
	"net/http"
	"strings"
)

type PurchaseCardRequest struct {
	UserName *string `json:"user_name"`

	MerchantId *string `json:"merchant_id"`

	// the amount purchased
	Amount *int64 `json:"amount,string,omitempty"`

	// amount of purchase that is revenue
	RevenueTake *int64 `json:"revenue_take,string,omitempty"`

	// amount of purchase that is expensed (ex. cc fees)
	Expenses *int64 `json:"expenses,string,omitempty"`
}

type PurchaseCardResponse struct {
	Transaction interface{} `json:"transaction"`
}

func PurchaseCard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	logger.Info(ctx, "received request to purchase card")
	decoder := json.NewDecoder(r.Body)
	var req PurchaseCardRequest
	err := decoder.Decode(&req)
	if err != nil {
		logger.Error(ctx, err, "error decoding request")
		http.Error(w, "unable to decode PurchaseCard request", http.StatusBadRequest)
		return
	}
	logger.Info(ctx, "got PurchaseCard request %v", req)
	if req.UserName == nil || req.MerchantId == nil || req.Amount == nil {
		logger.Error(ctx, nil, "none of userName, merchantId, or amount can be null")
		http.Error(w, "none of userName, merchantId, or amount can be null", http.StatusBadRequest)
		return
	}

	// check if the provided merchant id corresponds to an existing account
	merchantAccount, err := ledger.GetAccount(ctx, *req.MerchantId)
	if err != nil {
		logger.Error(ctx, err, "error getting merchant ledger account")
		http.Error(w, fmt.Sprintf("error getting merchant ledger account: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if merchantAccount == nil || merchantAccount.Metadata[balanceTypeKey] == nil {
		log.Print("merchant account nil")
		http.Error(w, fmt.Sprintf("no ledger account associated with address %s", *req.MerchantId), http.StatusBadRequest)
		return
	}

	cardId := fmt.Sprintf("cards:%s", strings.Replace(uuid.NewString(), "-", "", -1))
	metadata := map[string]interface{}{
		transactionTypeKey: purchaseCardTransaction,
		cardIdKey:          cardId,
		nameKey:            req.UserName,
		merchantIdKey:      req.MerchantId,
	}
	cardCreditAmount := *req.Amount
	if req.RevenueTake != nil && *req.RevenueTake != 0 {
		cardCreditAmount = *req.Amount - *req.RevenueTake
	}
	assetDebitAmount := *req.Amount
	if req.Expenses != nil && *req.Expenses != 0 {
		assetDebitAmount = *req.Amount - *req.Expenses
	}
	postings := []ledger.TransactionPosting{
		{
			Src:    worldAccountName,
			Dest:   cardId,
			Amount: cardCreditAmount,
		},
		{
			Src:    worldAccountName,
			Dest:   assetsAccountName,
			Amount: assetDebitAmount,
		},
	}
	if req.RevenueTake != nil && *req.RevenueTake != 0 {
		postings = append(postings, ledger.TransactionPosting{
			Src:    worldAccountName,
			Dest:   revenueAccountName,
			Amount: *req.RevenueTake,
		})
	}
	if req.Expenses != nil && *req.Expenses != 0 {
		postings = append(postings, ledger.TransactionPosting{
			Src:    worldAccountName,
			Dest:   expensesAccountName,
			Amount: *req.Expenses,
		})
	}
	txn, err := ledger.CreateTransactionWithPostings(ctx, metadata, postings)
	if err != nil {
		http.Error(w, "error creating transaction", http.StatusBadRequest)
		return
	}

	// add metadata to the account we just created
	accountMetadata := map[string]interface{}{
		nameKey:           *req.UserName,
		merchantIdKey:     *req.MerchantId,
		balanceTypeKey:    balanceTypeCredit,
		ledgerableTypeKey: ledgerableTypeExternal,
	}
	err = ledger.AddMetaDataToAccount(ctx, cardId, accountMetadata)
	if err != nil {
		http.Error(w, "error adding metadata to account", http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(
		PurchaseCardResponse{
			Transaction: txn,
		},
	)
	if err != nil {
		logger.Error(ctx, err, "error encoding response")
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
