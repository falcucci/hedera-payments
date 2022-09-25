package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/falcucci/maga-coin-payments-api/config"
	"github.com/falcucci/maga-coin-payments-api/utils/response"
	"github.com/gorilla/mux"

	"github.com/launchbadge/hedera-sdk-go"
)

type transferAmountRequestBody struct {
	TargetAccount int64  `json:"target_account"`
	PrivateKey    string `json:"private_key"`
	Amount        int64  `json:"amount"`
}

type account struct {
	Balance uint64 `json:"balance"`
}

// GetBalance get balance
func GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fmt.Println("vars", vars["id"])
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.GenerateHTTPResponse(w, http.StatusInternalServerError, response.GenerateErrorResponse(response.InternalServerError,
			"Error at convert id", "Was encountered an error when processing your request. We apologize for the inconvenience."))
		return
	}

	// Target account to get the balance for
	accountID := hedera.AccountID{Account: int64(id)}

	client, err := hedera.Dial(config.NodeAddress)
	if err != nil {
		panic(err)
	}

	client.SetNode(hedera.AccountID{Account: config.OperatorID})
	client.SetOperator(accountID, func() hedera.SecretKey {
		operatorSecret, err := hedera.SecretKeyFromString(config.OperatorKEY)
		if err != nil {
			panic(err)
		}

		return operatorSecret
	})

	defer client.Close()

	// Get the _answer_ for the query of getting the account balance
	balance, err := client.Account(accountID).Balance().Get()
	if err != nil {
		panic(err)
	}

	account := account{Balance: balance}

	response.GenerateHTTPResponse(
		w, http.StatusOK, response.GenerateSuccessResponse(
			account, 1, 1, 1))
}

// CashIn : Create new transaction cash in at maga-coin-payment-api
func CashIn(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var transferAmountRequestBody transferAmountRequestBody
	err := decoder.Decode(&transferAmountRequestBody)
	if err != nil {
		fmt.Println(err)
		response.GenerateHTTPResponse(w, http.StatusInternalServerError,
			response.GenerateErrorResponse(response.InternalServerError,
				"Error at convert id",
				"Was encountered an error when processing your request. We apologize for the inconvenience."))
		return
	}

	account, err := transaction(transferAmountRequestBody.TargetAccount,
		config.OperatorID, transferAmountRequestBody.Amount, config.OperatorKEY)
	if err != nil {
		fmt.Println(err)
		response.GenerateHTTPResponse(w, http.StatusInternalServerError,
			response.GenerateErrorResponse(response.InternalServerError,
				"Error transaction",
				"Was encountered an error when processing your request. We apologize for the inconvenience."))
		return
	}

	response.GenerateHTTPResponse(
		w, http.StatusOK, response.GenerateSuccessResponse(
			account, 1, 1, 1))
}

// CashOut : Create new transaction cash out at maga-coin-payment-api
func CashOut(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var transferAmountRequestBody transferAmountRequestBody
	err := decoder.Decode(&transferAmountRequestBody)
	if err != nil {
		response.GenerateHTTPResponse(w, http.StatusInternalServerError,
			response.GenerateErrorResponse(response.InternalServerError,
				"Error at convert id",
				"Was encountered an error when processing your request. We apologize for the inconvenience."))
		return
	}

	account, err := transaction(config.OperatorID,
		transferAmountRequestBody.TargetAccount,
		transferAmountRequestBody.Amount, transferAmountRequestBody.PrivateKey)
	if err != nil {
		errMessage := err.Error()
		hasInsufficientPayerBalance := strings.Contains(errMessage, "InsufficientPayerBalance")
		if hasInsufficientPayerBalance {
			response.GenerateHTTPResponse(w, http.StatusPreconditionFailed,
				response.GenerateErrorResponse(response.InternalServerError,
					"Insufficient Payer Balance",
					"Was encountered an error when processing your request. We apologize for the inconvenience."))
			return
		}
		response.GenerateHTTPResponse(w, http.StatusInternalServerError,
			response.GenerateErrorResponse(response.InternalServerError,
				"Error transaction",
				"Was encountered an error when processing your request. We apologize for the inconvenience."))
		return
	}

	response.GenerateHTTPResponse(
		w, http.StatusOK, response.GenerateSuccessResponse(
			account, 1, 1, 1))
}

func transaction(in, out, amount int64, privateKey string) (account, error) {
	var account account

	// Read and decode the operator secret key
	operatorAccountID := hedera.AccountID{Account: out}
	operatorSecret, err := hedera.SecretKeyFromString(privateKey)
	if err != nil {
		return account, err
	}

	// Read and decode target account
	fmt.Println("transferAmountRequestBody.TargetAccount", in)
	targetAccountID, err :=
		hedera.AccountIDFromString(fmt.Sprintf("0.0.%d", in))
	if err != nil {
		return account, err
	}

	//
	// Connect to Hedera
	//
	client, err := hedera.Dial(config.NodeAddress)
	if err != nil {
		return account, err
	}

	client.SetNode(hedera.AccountID{Account: config.NodeID})
	client.SetOperator(operatorAccountID, func() hedera.SecretKey {
		return operatorSecret
	})

	defer client.Close()

	//
	// Get balance for target account
	//
	balance, err := client.Account(operatorAccountID).Balance().Get()
	if err != nil {
		return account, err
	}

	//
	// Transfer 100 cryptos to target
	//
	nodeAccountID := hedera.AccountID{Account: config.NodeID}
	transaction, err := client.TransferCrypto().
		// Move 100 out of operator account
		Transfer(operatorAccountID, -amount).
		// And place in our new account
		Transfer(targetAccountID, amount).
		Operator(operatorAccountID).
		Node(nodeAccountID).
		Memo("[test] hedera-sdk-go v2").
		Fee(1000000000).
		Sign(operatorSecret). // Sign it once as operator
		Sign(operatorSecret). // And again as sender
		Execute()

	if err != nil {
		return account, err
	}

	fmt.Printf("transferred; transaction = %v\n", transaction.String())

	//
	// Get receipt to prove we sent ok
	//
	fmt.Printf("wait for 2s...\n")
	time.Sleep(2 * time.Second)

	receipt, err := client.Transaction(transaction).Receipt().Get()
	if err != nil {
		return account, err
	}

	if receipt.Status != hedera.StatusSuccess {
		panic(fmt.Errorf("transaction has a non-successful status: %v", receipt.Status.String()))
	}

	fmt.Printf("wait for 2s...\n")
	time.Sleep(2 * time.Second)

	//
	// Get balance for target account (again)
	//
	balance, err = client.Account(targetAccountID).Balance().Get()
	if err != nil {
		return account, err
	}

	account.Balance = balance

	return account, nil
}
