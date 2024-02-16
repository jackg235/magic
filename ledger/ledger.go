package ledger

import (
	"context"
	"errors"
	"fmt"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"math/big"
	"net/http"
	time2 "time"
)

const (
	ledgerName     = "gift-card-ledger"
	usdAsset       = "USD"
	formanceUrl    = "REDACTED"
	formanceSecret = "REDACTED"
)

var formanceClient = formance.New(
	formance.WithServerURL(formanceUrl),
	formance.WithSecurity(shared.Security{
		Authorization: fmt.Sprintf("Bearer %s", formanceSecret),
	}),
)

func AddMetaDataToAccount(ctx context.Context, address string, metadata map[string]interface{}) error {
	res, err := formanceClient.Ledger.AddMetadataToAccount(ctx, operations.AddMetadataToAccountRequest{
		RequestBody: metadata,
		Address:     address,
		Ledger:      ledgerName,
	})
	if err != nil {
		return err
	}
	if res.StatusCode >= http.StatusBadRequest {
		return errors.New(fmt.Sprintf("failed to add metadata to account with error code %d", res.StatusCode))
	}
	return nil
}

func GetAccount(ctx context.Context, address string) (*shared.AccountWithVolumesAndBalances, error) {
	res, err := formanceClient.Ledger.GetAccount(ctx, operations.GetAccountRequest{
		Address: address,
		Ledger:  ledgerName,
	})
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(fmt.Sprintf("failed to get ledger account with error code %d", res.StatusCode))
	}
	if res.AccountResponse == nil || len(res.AccountResponse.Data.Address) == 0 {
		return nil, nil
	}
	return &res.AccountResponse.Data, nil
}

func ListAccounts(ctx context.Context) ([]shared.Account, error) {
	res, err := formanceClient.Ledger.ListAccounts(ctx, operations.ListAccountsRequest{
		Ledger:   ledgerName,
		PageSize: formance.Int64(500), // will need to paginate but formance pagination is broken
	})
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(fmt.Sprintf("failed to get ledger account with error code %d", res.StatusCode))
	}

	return res.AccountsCursorResponse.Cursor.Data, nil
}

func ListTransactions(ctx context.Context) ([]shared.Transaction, error) {
	res, err := formanceClient.Ledger.ListTransactions(ctx, operations.ListTransactionsRequest{
		Ledger:   ledgerName,
		PageSize: formance.Int64(500), // will need to paginate but formance pagination is broken

	})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(fmt.Sprintf("failed to get ledger account with error code %d", res.StatusCode))
	}

	return res.TransactionsCursorResponse.Cursor.Data, nil
}

func ListBalances(ctx context.Context) (map[string]int64, error) {
	cursor := ""
	accountToBalance := make(map[string]int64)
	for {
		res, err := formanceClient.Ledger.GetBalances(ctx, operations.GetBalancesRequest{
			Ledger: ledgerName,
			Cursor: &cursor,
		})
		if err != nil {
			return nil, err
		}
		for _, balances := range res.BalancesCursorResponse.Cursor.Data {
			for acct, balance := range balances {
				accountToBalance[acct] = balance[usdAsset].Int64()
			}
		}
		if !res.BalancesCursorResponse.Cursor.HasMore {
			break
		}
		cursor = *res.BalancesCursorResponse.Cursor.Next
	}
	return accountToBalance, nil
}

func CreateTransactionWithPostings(ctx context.Context, metadata map[string]interface{}, postings []TransactionPosting) (*shared.Transaction, error) {
	time := time2.Now()
	formancePostings := make([]shared.Posting, len(postings))
	for i, p := range postings {
		formancePostings[i] = shared.Posting{
			Amount:      big.NewInt(p.Amount),
			Asset:       usdAsset,
			Destination: p.Dest,
			Source:      p.Src,
		}
	}
	res, err := formanceClient.Ledger.CreateTransaction(ctx, operations.CreateTransactionRequest{
		PostTransaction: shared.PostTransaction{
			Metadata:  metadata,
			Postings:  formancePostings,
			Timestamp: &time,
		},
		Ledger: ledgerName,
	})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(*res.ErrorResponse.ErrorMessage)
	}
	if res.TransactionsResponse == nil || len(res.TransactionsResponse.Data) == 0 {
		return nil, errors.New("expected to create a transaction but none were created")
	}
	return &res.TransactionsResponse.Data[0], nil
}
