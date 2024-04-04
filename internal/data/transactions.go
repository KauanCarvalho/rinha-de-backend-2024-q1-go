package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Transaction struct {
	ID          int64     `json:"-"`
	Amount      int64     `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
	CustomerID  int64     `json:"-"`
}

type TransactionModel struct {
	DB *sql.DB
}

func (m TransactionModel) Insert(ctx context.Context, tx *sql.Tx, transaction *Transaction) error {
	query := `
		INSERT INTO transactions (amount, type, description, customer_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	args := []interface{}{transaction.Amount, transaction.Type, transaction.Description, transaction.CustomerID}

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m TransactionModel) LastNTransactionForCustomer(cutomerId int64, n int) ([]Transaction, error) {
	query := `
		SELECT amount, type, description, created_at FROM transactions
		WHERE customer_id = $1
		ORDER BY id DESC LIMIT $2`

	rows, err := m.DB.Query(query, cutomerId, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []Transaction{}

	for rows.Next() {
		var transaction Transaction

		err := rows.Scan(
			&transaction.Amount,
			&transaction.Type,
			&transaction.Description,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func ValidateTransaction(t Transaction) error {
	if t.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if t.Type != "c" && t.Type != "d" {
		return errors.New("invalid operation type")
	}

	if descriptionLength := len(t.Description); descriptionLength > 10 || descriptionLength < 1 {
		return errors.New("description must be between 1 and 10 characters")
	}

	return nil
}
