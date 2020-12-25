package dtos

// Data transformation object
type TransactionDTO struct {
	Id              uint    `json:"id"`
	AccountId       uint    `json:"account_id"`
	Amount          float64 `json:"amount"`
	Bank            string  `json:"bank"`
	TransactionType string  `json:"transaction_type"`
	CreatedAt       string  `json:"created_at"`
}
