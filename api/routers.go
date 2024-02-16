package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"PurchaseCard",
		http.MethodPost,
		"/card/purchase",
		PurchaseCard,
	},
	Route{
		"SpendCard",
		http.MethodPost,
		"/card/spend",
		SpendCard,
	},
	Route{
		"CreateMerchant",
		http.MethodPost,
		"/merchant/create",
		CreateMerchant,
	},
	Route{
		"PayoutMerchant",
		http.MethodPost,
		"/merchant/payout",
		PayoutMerchant,
	},
	Route{
		"ListAccounts",
		http.MethodGet,
		"/accounts",
		ListAccounts,
	},
	Route{
		"ListTransactions",
		http.MethodGet,
		"/transactions",
		ListTransactions,
	},
	Route{
		"LedgerMetadata",
		http.MethodGet,
		"/ledger",
		LedgerMetadata,
	},
}
