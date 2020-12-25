package requests

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"simple-restful/pkg/core"
	"simple-restful/pkg/core/utils"
	"strconv"
)

type GetTransactionsRequest struct {
	UserId    uint `json:"user_id"`
	AccountId uint `json:"account_id"`
}

type CreateTransactionRequest struct {
	UserId          uint    `json:"user_id"`
	AccountId       uint    `json:"account_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
}

func GetTransactionsRequestParser(r *http.Request) (*GetTransactionsRequest, error) {
	req := &GetTransactionsRequest{}
	vars := mux.Vars(r)

	userId, err := strconv.Atoi(vars["user_id"])
	if err != nil || userId <= 0 {
		return nil, fmt.Errorf("user is invalid")
	}
	req.UserId = uint(userId)

	accIdStr := r.URL.Query().Get("account_id")
	if accIdStr != "" {
		accountId, err := strconv.ParseInt(accIdStr, 10, 32)
		if err != nil || accountId < 0 {
			return nil, fmt.Errorf("account_id is invalid")
		}
		req.AccountId = uint(accountId)
	}

	return req, nil
}

func CreateTransactionRequestParser(r *http.Request) (*CreateTransactionRequest, error) {
	req := &CreateTransactionRequest{}
	vars := mux.Vars(r)

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return nil, fmt.Errorf("json payload is invalid")
	}

	if req.AccountId <= 0 {
		return nil, fmt.Errorf("account_id is invalid")
	}

	if req.TransactionType == "" {
		return nil, fmt.Errorf("transaction type cannot be empty")
	}

	if !utils.IsStringSliceContains(core.AvailableTransactionTypes, req.TransactionType) {
		return nil, fmt.Errorf("transaction type is invalid, must be is one of %v", core.AvailableTransactionTypes)
	}

	userId, err := strconv.Atoi(vars["user_id"])
	if err != nil || userId <= 0 {
		return nil, fmt.Errorf("user is invalid")
	}
	req.UserId = uint(userId)

	return req, nil
}
