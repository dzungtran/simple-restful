package models

type Account struct {
	Id     uint `gorm:"primary_key"`
	Name   string
	Bank   string
	UserId uint
}
