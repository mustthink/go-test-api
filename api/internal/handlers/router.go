package handlers

import "net/http"

func (app *application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/transactions", app.showTransactions)

	return mux
}
