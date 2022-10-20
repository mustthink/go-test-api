package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (app *application) showTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	sender := r.URL.Query().Get("from")
	receiver := r.URL.Query().Get("to")
	id := r.URL.Query().Get("id")
	block := r.URL.Query().Get("block")
	page := r.URL.Query().Get("page")
	validation := r.URL.Query().Get("v")
	if sender == "" && receiver == "" && id == "" && block == "" {
		app.clientError(w, http.StatusBadGateway)
	}

	transactions, err := app.service.GetTransactions(sender, receiver, id, block, page)
	if transactions == nil || len(transactions) == 0 {
		app.clientError(w, http.StatusNotFound)
	}
	if err != nil {
		app.serverError(w, err)
	}

	json_data, err := json.MarshalIndent(transactions, "", "    ")
	if err != nil {
		app.serverError(w, err)
	}

	if validation == "true" {
		for _, v := range transactions {
			if app.Validation(v) {
				log.Println("ID: ", v.ID, " is valid")
			} else {
				log.Println("ID: ", v.ID, " isn't valid")
			}
			time.Sleep(time.Millisecond * 200)
		}
	}

	fmt.Fprintf(w, "%v", string(json_data))
}
