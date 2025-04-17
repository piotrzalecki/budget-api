package handlers

import (
	"net/http"
)

func (rep *Repository) Logs(w http.ResponseWriter, r *http.Request) {
	logs, err := rep.App.Models.Log.AllLogs()
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    envelope{"logs": logs},
	}

	rep.WriteJSON(w, http.StatusOK, payload)
}
