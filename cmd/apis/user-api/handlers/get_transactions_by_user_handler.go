package handlers

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"simple-restful/cmd/apis/user-api/requests"
	"simple-restful/pkg/core/servehttp"
	"simple-restful/pkg/models"
	"simple-restful/pkg/transformers"
)

type GetUserTransactionsHandler struct {
	DB *gorm.DB
}

func (h *GetUserTransactionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := requests.GetTransactionsRequestParser(r)
	if err != nil {
		log.Printf("Error bad request: %s", err.Error())
		servehttp.ResponseJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

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

	// Get transaction model list
	transactions, err := h.getListTransaction(req)
	if err != nil {
		log.Printf("Error when get transactions: %s", err.Error())
		servehttp.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Internal server error.",
		})
		return
	}

	// Transform transaction model to response
	transactionsDTO, err := transformers.TransformTransactionModelsToTransactionDTOs(h.DB, transactions)
	if err != nil {
		log.Printf("Error when transform transactions to dto: %s", err.Error())
		servehttp.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Internal server error.",
		})
		return
	}
	servehttp.ResponseSuccessJSON(w, transactionsDTO)
}

func (h *GetUserTransactionsHandler) getListTransaction(req *requests.GetTransactionsRequest) ([]*models.Transaction, error) {
	transactions := make([]*models.Transaction, 0)
	andConds := map[string]interface{}{
		"user_id": req.UserId,
	}

	if req.AccountId > 0 {
		andConds["account_id"] = req.AccountId
	}

	err := h.DB.Where(andConds).Limit(200).Find(&transactions).Error
	if err != nil {
		log.Printf("Get transactions error: %v", err)
		return nil, fmt.Errorf("some error has occurred")
	}
	return transactions, nil
}
