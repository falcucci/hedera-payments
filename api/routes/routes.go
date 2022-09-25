package routes

import (
	"github.com/falcucci/maga-coin-payments-api/api/account"
	"github.com/falcucci/maga-coin-payments-api/api/ping"
	"github.com/falcucci/maga-coin-payments-api/api/wallet"
	"github.com/gorilla/mux"
)

// Configure : Configure router API
func Configure() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", ping.Get).Methods("GET")
	r.HandleFunc("/wallet/balance/{id}", wallet.GetBalance).Methods("GET")
	r.HandleFunc("/wallet/cash-in", wallet.CashIn).Methods("POST")
	r.HandleFunc("/wallet/cash-out", wallet.CashOut).Methods("POST")
	r.HandleFunc("/account/create", account.CreateAccount).Methods("POST")
	return r
}
