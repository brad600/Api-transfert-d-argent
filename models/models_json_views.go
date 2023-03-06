package APIServer

import (
	"time"
)

// AccountJsonView detient le ID et le montant de l'argent
type AccountJsonView struct {
	AccountID int64 `json:"account_id"`
	Balance   int64 `json:"balance"`
}

// AccountIDJsonView ...
type AccountIDJsonView struct {
	ID int64 `json:"account_id"`
}

// TransactionJsonView detient les données nécessaires pour effectuer un transfert d'argent
type TransactionJsonView struct {
	Timestamp     time.Time `json:"timestamp"`
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Amount        int64     `json:"amount"`
}
