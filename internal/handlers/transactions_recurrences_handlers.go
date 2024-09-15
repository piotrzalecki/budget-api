package handlers

import (
	"net/http"

	"github.com/piotrzalecki/budget-api/internal/data"
)

func (rep *Repository) TransactionsRecurrences(w http.ResponseWriter, r *http.Request) {

	trs, err := rep.App.Models.TransactionRecurrence.AllTransactionsRecurrences()
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    envelope{"transaction_recurrences": trs},
	}

	rep.WriteJSON(w, http.StatusOK, payload)

}

func (rep *Repository) TransactionsRecurrencesCreate(w http.ResponseWriter, r *http.Request) {
	var tr data.TransactionRecurrence

	err := rep.readJSON(w, r, &tr)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	_, err = rep.App.Models.TransactionRecurrence.CreateTransactionsRecurrences(tr)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Transaction recurrence created!",
	}

	rep.WriteJSON(w, http.StatusAccepted, payload)
}
