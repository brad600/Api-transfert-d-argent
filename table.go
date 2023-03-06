package store

import (
	"POD_TRANSFERT_ARGENT_API/models"
)

// Store ...
type table interface {
	InsertAccount(balance int64) (models.Account, error)
	DeleteAccount(accountId int64) error
	GetAccount(accountId int64) (models.Account, error)
	TransferMoney(accountToId, accountFromId, amount int64) error
	GetTransactionsHistory(accountId, nLastDays, limit int64) ([]models.Transaction, error)
}
