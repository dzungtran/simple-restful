package models

import "time"

type Transaction struct {
	Id              uint `gorm:"primary_key"`
	UserId          uint
	AccountId       uint
	Amount          float64
	TransactionType string
	CreatedAt       time.Time
}
