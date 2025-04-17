package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/piotrzalecki/budget-api/internal/data"
)

func (rep *Repository) Tags(w http.ResponseWriter, r *http.Request) {

	tags, err := rep.App.Models.Tag.AllTags()
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    envelope{"tags": tags},
	}

	rep.WriteJSON(w, http.StatusOK, payload)

}

func (rep *Repository) TagsCreateUpdate(w http.ResponseWriter, r *http.Request) {
	var tag data.Tag

	err := rep.readJSON(w, r, &tag)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	if tag.Id == 0 {
		_, err := rep.App.Models.Tag.CreateTag(tag)
		if err != nil {
			rep.errorJson(w, err)
			return
		}
	} else {
		err := rep.App.Models.Tag.UpdateTag(tag)
		if err != nil {
			rep.errorJson(w, err)
			return
		}
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Tag saved!",
	}

	rep.WriteJSON(w, http.StatusAccepted, payload)
}

func (rep *Repository) TagsDelete(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		ID int `json:"id"`
	}

	err := rep.readJSON(w, r, &requestPayload)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	err = rep.App.Models.Tag.DeleteTag(requestPayload.ID)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Tag deleted!",
	}

	rep.WriteJSON(w, http.StatusOK, payload)
}

func (rep *Repository) TagById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	tag, err := rep.App.Models.Tag.GetTagById(idInt)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    tag,
	}
	rep.WriteJSON(w, http.StatusOK, payload)
}
