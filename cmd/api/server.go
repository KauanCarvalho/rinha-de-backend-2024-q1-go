package main

import (
	"fmt"
	"net/http"
)

func (app *application) server() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
