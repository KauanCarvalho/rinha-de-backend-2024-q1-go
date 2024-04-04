package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/clientes/:id/transacoes", app.createTransaction)
	router.HandlerFunc(http.MethodGet, "/clientes/:id/extrato", app.getBankStatement)

	return router
}
