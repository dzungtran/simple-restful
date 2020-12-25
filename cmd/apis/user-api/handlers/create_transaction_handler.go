package handlers

import (
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"simple-restful/cmd/apis/user-api/requests"
	"simple-restful/pkg/core/servehttp"
	"simple-restful/pkg/models"
	"simple-restful/pkg/transformers"
)

type CreateTransactionHandler struct {
	DB *gorm.DB
}

func (h *CreateTransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := requests.CreateTransactionRequestParser(r)
	if err != nil {
		log.Printf("Error bad request: %s", err.Error())
		servehttp.ResponseJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Check user exists
	var user models.User
	err = h.DB.First(&user, req.UserId).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			servehttp.ResponseJSON(w, http.StatusNotFound, map[string]string{
				"error": "Not found.",
			})
			return
		}
		log.Printf("Error when get user: %v", err)
		servehttp.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Internal server error.",
		})
		return
	}

	var account models.Account
	err = h.DB.First(&account, req.AccountId).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			servehttp.ResponseJSON(w, http.StatusNotFound, map[string]string{
				"error": "Account not found.",
			})
			return
		}
		log.Printf("Error when get account: %v", err)
		servehttp.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Internal server error.",
		})
		return
	}

	if account.UserId != user.Id {
		servehttp.ResponseJSON(w, http.StatusForbidden, map[string]string{
			"error": "Account not belong with user.",
		})
		return
	}

	transaction := h.convertCreateTransactionRequestToTransactionModel(req, account)
	err = h.DB.Create(transaction).Error
	if err != nil {
		log.Printf("Error when get create transaction: %v", err)
		servehttp.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Internal server error.",
		})
		return
	}
	servehttp.ResponseJSON(w, http.StatusCreated, transformers.TransformTransactionModelToTransactionDTO(&account, transaction))
}

func (h *CreateTransactionHandler) convertCreateTransactionRequestToTransactionModel(
	req *requests.CreateTransactionRequest,
	account models.Account,
) *models.Transaction {
	return &models.Transaction{
		UserId:          account.UserId,
		AccountId:       account.Id,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
	}
}
