package transformers

import (
	"github.com/jinzhu/gorm"
	"simple-restful/pkg/dtos"
	"simple-restful/pkg/models"
)

const DateTimeFormat = "2006-01-02 15:04:05 -0700"

func TransformTransactionModelToTransactionDTO(account *models.Account, transaction *models.Transaction) *dtos.TransactionDTO {
	createdAt := transaction.CreatedAt.Format(DateTimeFormat)
	return &dtos.TransactionDTO{
		Id:              transaction.Id,
		AccountId:       transaction.AccountId,
		Amount:          transaction.Amount,
		Bank:            account.Bank,
		TransactionType: transaction.TransactionType,
		CreatedAt:       createdAt,
	}
}

func TransformTransactionModelsToTransactionDTOs(db *gorm.DB, transactions []*models.Transaction) ([]*dtos.TransactionDTO, error) {
	transactionsDTO := make([]*dtos.TransactionDTO, 0)
	listAccount := make([]*models.Account, 0)
	mapAccounts := make(map[uint]*models.Account)
	accountIds := make([]uint, 0)

	for _, t := range transactions {
		accountIds = append(accountIds, t.AccountId)
	}

	err := db.Where("id in (?)", accountIds).Find(&listAccount).Error
	if err != nil {
		return nil, err
	}

	for _, acc := range listAccount {
		mapAccounts[acc.Id] = acc
	}

	for _, t := range transactions {
		acc, ok := mapAccounts[t.AccountId]
		if !ok {
			continue
		}
		transactionsDTO = append(transactionsDTO, TransformTransactionModelToTransactionDTO(acc, t))
	}

	return transactionsDTO, nil
}
