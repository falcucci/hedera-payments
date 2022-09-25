package main

import (
	"log"
	"net/http"
	"os"

	"github.com/falcucci/maga-coin-payments-api/api/account"
	"github.com/falcucci/maga-coin-payments-api/api/ping"
	"github.com/falcucci/maga-coin-payments-api/api/wallet"
	"github.com/falcucci/maga-coin-payments-api/config"
	"github.com/gorilla/mux"
)

func main() {

	// carrega variaveis de ambiente
	config.LoadEnvVars()

	// Configure Routes
	r := mux.NewRouter()
	r.HandleFunc("/ping", ping.Get).Methods("GET")
	r.HandleFunc("/wallet/balance/{id}", wallet.GetBalance).Methods("GET")
	r.HandleFunc("/wallet/cash-in", wallet.CashIn).Methods("POST")
	r.HandleFunc("/wallet/cash-out", wallet.CashOut).Methods("POST")
	r.HandleFunc("/account/create", account.CreateAccount).Methods("POST")

	s := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}

	log.Printf("Iniciando API na porta %s", os.Getenv("PORT"))
	log.Fatal(s.ListenAndServe())

}
