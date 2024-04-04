package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Customer struct {
	ID                int64     `json:"-"`
	BankLimit         int64     `json:"limite"`
	BankBalance       int64     `json:"total"`
	BankStatementDate time.Time `json:"data_extrato"`
}

type CustomerModel struct {
	DB *sql.DB
}

func (m CustomerModel) GetForUpdate(ctx context.Context, tx *sql.Tx, id int64) (*Customer, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, bank_limit, bank_balance
		FROM customers
		WHERE id = $1
		FOR UPDATE`

	var customer Customer

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&customer.ID,
		&customer.BankLimit,
		&customer.BankBalance,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &customer, nil
}

func (m CustomerModel) Get(id int64) (*Customer, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, bank_limit, bank_balance
		FROM customers
		WHERE id = $1`

	var customer Customer

	err := m.DB.QueryRow(query, id).Scan(
		&customer.ID,
		&customer.BankLimit,
		&customer.BankBalance,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &customer, nil
}

func (m CustomerModel) Update(ctx context.Context, tx *sql.Tx, customer *Customer) error {
	query := `
		UPDATE customers
		SET bank_balance = $2
		WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, customer.ID, customer.BankBalance)
	if err != nil {
		return err
	}

	return nil
}

func ValidateNewBalance(customer Customer, operationType string, amount int64) (int64, error) {
	var newBalance int64

	if operationType == "c" {
		newBalance = customer.BankBalance + amount
	} else if operationType == "d" {
		newBalance = customer.BankBalance - amount
	} else {
		return 0, errors.New("invalid operation type")
	}

	if newBalance < -customer.BankLimit {
		return 0, errors.New("insufficient funds")
	}

	return newBalance, nil
}
