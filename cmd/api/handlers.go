package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"backendfight.kauancarvalho/internal/data"
)

func (app *application) createTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Valor     int64  `json:"valor"`
		Tipo      string `json:"tipo"`
		Descricao string `json:"descricao"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	transaction := data.Transaction{
		Amount:      input.Valor,
		Type:        input.Tipo,
		Description: input.Descricao,
		CustomerID:  id,
	}

	err = data.ValidateTransaction(transaction)
	if err != nil {
		app.failedValidationResponse(w, r, map[string]string{"transaction": err.Error()})
		return
	}

	ctx := context.Background()
	tx, err := app.db.BeginTx(ctx, nil)
	if err != nil {
		app.serverErrorResponse(w, r)
		return
	}
	defer tx.Rollback()

	customer, err := app.models.Customers.GetForUpdate(ctx, tx, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r)
		}
		tx.Rollback()
		return
	}

	newBalance, err := data.ValidateNewBalance(*customer, input.Tipo, input.Valor)
	if err != nil {
		app.failedValidationResponse(w, r, map[string]string{"bank_balance": err.Error()})
		tx.Rollback()
		return
	}

	customer.BankBalance = newBalance

	err = app.models.Transactions.Insert(ctx, tx, &transaction)
	if err != nil {
		tx.Rollback()
		app.serverErrorResponse(w, r)
		return
	}

	err = app.models.Customers.Update(ctx, tx, customer)
	if err != nil {
		tx.Rollback()
		app.serverErrorResponse(w, r)
		return
	}

	tx.Commit()

	err = app.writeJSON(w, http.StatusOK, customer)
	if err != nil {
		app.serverErrorResponse(w, r)
	}
}

func (app *application) getBankStatement(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	customer, err := app.models.Customers.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r)
		}
		return
	}

	transactions, err := app.models.Transactions.LastNTransactionForCustomer(id, 10)
	if err != nil {
		app.serverErrorResponse(w, r)
		return
	}

	customer.BankStatementDate = time.Now()

	result := struct {
		Customer     data.Customer      `json:"saldo"`
		Transactions []data.Transaction `json:"ultimas_transacoes"`
	}{
		Customer:     *customer,
		Transactions: transactions,
	}

	err = app.writeJSON(w, http.StatusOK, result)
	if err != nil {
		app.serverErrorResponse(w, r)
	}
}
