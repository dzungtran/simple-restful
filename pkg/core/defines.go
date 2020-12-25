package core

const (
	TransactionTypeWithDraw = "withdraw"
	TransactionTypeDeposit  = "deposit"

	BankNameVCB = "VCB"
	BankNameACB = "ACB"
	BankNameVIB = "VIB"
)

var AvailableTransactionTypes = []string{
	TransactionTypeWithDraw,
	TransactionTypeDeposit,
}

var AvailableBanks = []string{
	BankNameVCB,
	BankNameACB,
	BankNameVIB,
}