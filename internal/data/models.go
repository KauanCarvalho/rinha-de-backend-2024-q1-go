package data

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Customers    CustomerModel
	Transactions TransactionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Customers:    CustomerModel{DB: db},
		Transactions: TransactionModel{DB: db},
	}
}
