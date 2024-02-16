package api

type TransactionType string
type LedgerableType string
type BalanceType string

const (
	cardIdKey                                         = "card_id"
	nameKey                                           = "name"
	merchantIdKey                                     = "merchant_id"
	balanceTypeKey                                    = "balance_type"
	ledgerableTypeKey                                 = "ledgerable_type"
	purchaseIdKey                                     = "purchase_id"
	transactionTypeKey                                = "transaction_type"
	assetsAccountName                                 = "assets"
	revenueAccountName                                = "revenue"
	expensesAccountName                               = "expenses"
	worldAccountName                                  = "world"
	purchaseCardTransaction           TransactionType = "purchase_card"
	spendCardTransaction              TransactionType = "spend_card"
	payoutMerchantTransaction         TransactionType = "payout_merchant"
	createMerchantTransaction         TransactionType = "create_merchant"
	createInternalAccountsTransaction TransactionType = "create_internal_account"

	balanceTypeCredit      BalanceType    = "credit"
	balanceTypeDebit       BalanceType    = "debit"
	ledgerableTypeInternal LedgerableType = "internal"
	ledgerableTypeExternal LedgerableType = "external"
)
