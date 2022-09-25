package account

import (
	"fmt"
	"net/http"
	"time"

	"github.com/falcucci/maga-coin-payments-api/config"
	"github.com/falcucci/maga-coin-payments-api/utils/response"

	"github.com/launchbadge/hedera-sdk-go"
)

type account struct {
	Number     int64  `json:"number"`
	PrivateKey string `json:"private_key"`
}

// CreateAccount create account
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	//
	// Generate keys
	//

	// Read and decode the operator secret key
	operatorSecret, err := hedera.SecretKeyFromString(config.OperatorKEY)
	if err != nil {
		response.GenerateHTTPResponse(w, http.StatusInternalServerError,
			response.GenerateErrorResponse(response.InternalServerError,
				"Error transaction",
				"Was encountered an error when processing your request. We apologize for the inconvenience."))
		return
	}

	// Generate a new keypair for the new account
	secret, _ := hedera.GenerateSecretKey()
	public := secret.Public()

	fmt.Printf("secret = %v\n", secret)
	fmt.Printf("public = %v\n", public)

	//
	// Connect to Hedera
	//
	client, err := hedera.Dial(config.NodeAddress)
	if err != nil {
		response.GenerateHTTPResponse(w, http.StatusInternalServerError,
			response.GenerateErrorResponse(response.InternalServerError,
				"Error transaction",
				"Was encountered an error when processing your request. We apologize for the inconvenience."))
		return
	}

	defer client.Close()

	//
	// Send transaction to create account
	//
	nodeAccountID := hedera.AccountID{Account: config.NodeID}
	operatorAccountID := hedera.AccountID{Account: config.OperatorID}
	transaction, err := client.CreateAccount().
		Key(public).
		InitialBalance(0).
		Operator(operatorAccountID).
		Node(nodeAccountID).
		Memo("[test] hedera-sdk-go v2").
		Fee(1000000000).
		Sign(operatorSecret).
		Execute()

	if err != nil {
		response.GenerateHTTPResponse(w, http.StatusInternalServerError,
			response.GenerateErrorResponse(response.InternalServerError,
				"Error transaction",
				"Was encountered an error when processing your request. We apologize for the inconvenience."))
		return
	}

	fmt.Printf("created account; transaction = %v\n", transaction.String())

	//
	// Get receipt to prove we created it ok
	//
	fmt.Printf("wait for 2s...\n")
	time.Sleep(2 * time.Second)

	receipt, err := client.Transaction(transaction).Receipt().Get()
	if err != nil {
		response.GenerateHTTPResponse(w, http.StatusInternalServerError,
			response.GenerateErrorResponse(response.InternalServerError,
				"Error transaction",
				"Was encountered an error when processing your request. We apologize for the inconvenience."))
		return
	}

	if receipt.Status != hedera.StatusSuccess {
		response.GenerateHTTPResponse(w, http.StatusInternalServerError,
			response.GenerateErrorResponse(response.InternalServerError,
				fmt.Sprintf("transaction has a non-successful status: %s", receipt.Status.String()),
				"Was encountered an error when processing your request. We apologize for the inconvenience."))
		return
	}

	account := account{
		Number:     receipt.AccountID.Account,
		PrivateKey: secret.String(),
	}

	response.GenerateHTTPResponse(
		w, http.StatusOK, response.GenerateSuccessResponse(
			account, 1, 1, 1))

	fmt.Printf("account = %v\n", *receipt.AccountID)
}
