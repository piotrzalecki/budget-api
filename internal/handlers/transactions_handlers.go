package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/piotrzalecki/budget-api/internal/data"
)

type Summary struct {
	Active    float32
	NonActive float32
}

func (rep *Repository) TransactionsAll(w http.ResponseWriter, r *http.Request) {

	trans, err := rep.App.Models.Transaction.AllTransactions()
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    envelope{"transactions": trans},
	}

	rep.WriteJSON(w, http.StatusOK, payload)
}

func (rep *Repository) TransactionsById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	trans, err := rep.App.Models.Transaction.GetTransactionById(idInt)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    trans,
	}
	rep.WriteJSON(w, http.StatusOK, payload)

}

func (rep *Repository) TransactionsCreateUpdate(w http.ResponseWriter, r *http.Request) {
	var transaction data.Transaction

	err := rep.readJSON(w, r, &transaction)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	if transaction.Id == 0 {
		newId, err := rep.App.Models.Transaction.CreateTransaction(transaction)
		if err != nil {
			rep.errorJson(w, err)
			return
		}
		rep.App.Models.Log.AddLog(fmt.Sprintf("transaction id %d created", newId))
	} else {
		err := rep.App.Models.Transaction.UpdateTransaction(transaction)
		if err != nil {
			rep.errorJson(w, err)
			return
		}
		rep.App.Models.Log.AddLog(fmt.Sprintf("transaction id %d updated", transaction.Id))
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Transaction saved!",
	}

	rep.WriteJSON(w, http.StatusAccepted, payload)
}

func (rep *Repository) TransactionsDelete(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		ID int `json:"id"`
	}

	err := rep.readJSON(w, r, &requestPayload)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	err = rep.App.Models.Transaction.DeleteTransaction(requestPayload.ID)
	if err != nil {
		rep.errorJson(w, err)
		return
	}
	rep.App.Models.Log.AddLog(fmt.Sprintf("transaction id %d deleted", requestPayload.ID))
	payload := jsonResponse{
		Error:   false,
		Message: "Transaction deleted!",
	}

	rep.WriteJSON(w, http.StatusOK, payload)

}

func (rep *Repository) TransactionsSetStatusAllActive(w http.ResponseWriter, r *http.Request) {
	err := rep.App.Models.Transaction.TransactionsSetAllActive()
	if err != nil {
		rep.errorJson(w, err)
		return
	}
	rep.App.Models.Log.AddLog("status of all transactions has been set to ACTIVE")
	payload := jsonResponse{
		Error:   false,
		Message: "All transactions set to active!",
	}

	rep.WriteJSON(w, http.StatusOK, payload)
}

func (rep *Repository) TransactionsSetStatus(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		ID     int `json:"id"`
		Status bool `json:"status"`
	}

	err := rep.readJSON(w, r, &requestPayload)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	err = rep.App.Models.Transaction.TransactionSetStatus(requestPayload.ID, requestPayload.Status)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Transaction status updated!",
	}
	rep.WriteJSON(w, http.StatusOK, payload)
}

// // DONE

// func transactionsSummary(trans []models.Transaction) Summary {
// 	var active float32
// 	var nonactive float32
// 	for _, tr := range trans {
// 		if tr.Active == true {
// 			active += tr.Quote
// 		} else {
// 			nonactive += tr.Quote
// 		}
// 	}

// 	return Summary{
// 		Active:    active,
// 		NonActive: nonactive,
// 	}
// }
