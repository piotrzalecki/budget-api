package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/piotrzalecki/budget-api/internal/data"
)

func (rep *Repository) Budgets(w http.ResponseWriter, r *http.Request) {

	buds, err := rep.App.Models.Budget.AllBudgets()
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    envelope{"budgets": buds},
	}

	rep.WriteJSON(w, http.StatusOK, payload)
}

func (rep *Repository) BudgetsCreateUpdate(w http.ResponseWriter, r *http.Request) {
	var budget data.Budget

	err := rep.readJSON(w, r, &budget)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	if budget.Id == 0 {
		_, err := rep.App.Models.Budget.CreateBudget(budget)
		if err != nil {
			rep.errorJson(w, err)
			return
		}
	} else {
		err := rep.App.Models.Budget.UpdateBudget(budget)
		if err != nil {
			rep.errorJson(w, err)
			return
		}
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Budget saved!",
	}

	rep.WriteJSON(w, http.StatusAccepted, payload)
}

func (rep *Repository) BudgetsDelete(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		ID int `json:"id"`
	}

	err := rep.readJSON(w, r, &requestPayload)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	err = rep.App.Models.Budget.DeleteBudget(requestPayload.ID)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Budget deleted!",
	}

	rep.WriteJSON(w, http.StatusOK, payload)

}

func (rep *Repository) BudgetsById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	budget, err := rep.App.Models.Budget.GetBudgetById(idInt)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    budget,
	}
	rep.WriteJSON(w, http.StatusOK, payload)
}
