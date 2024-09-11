package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/piotrzalecki/budget-api/internal/data"

	"github.com/go-chi/chi/v5"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type envelope map[string]interface{}

func (rep *Repository) Login(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		UserName string `json:"email"`
		Password string `json:"password"`
	}

	var creds credentials
	var payload jsonResponse

	err := rep.readJSON(w, r, &creds)
	if err != nil {
		rep.App.ErrorLogger.Println(err)
		payload.Error = true
		payload.Message = "invalid json"
		rep.WriteJSON(w, http.StatusBadRequest, payload)
	}

	user, err := rep.App.Models.User.GetByEmail(creds.UserName)
	if err != nil {
		rep.errorJson(w, errors.New("invalid username/password"))
		return
	}

	validPassword, err := user.CheckPassword(creds.Password)
	if err != nil || !validPassword {
		rep.errorJson(w, errors.New("02 invalid username/password"))
		return
	}

	if user.Active == 0 {
		rep.errorJson(w, errors.New("user is not active"))
		return
	}

	token, err := rep.App.Models.Token.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	err = rep.App.Models.Token.Insert(*token, *user)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload = jsonResponse{
		Error:   false,
		Message: "logged in",
		Data:    envelope{"token": token, "user": user},
	}

	err = rep.WriteJSON(w, http.StatusOK, payload)

	if err != nil {
		rep.App.ErrorLogger.Println(err)
	}

}

func (rep *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Token string `json:"token"`
	}

	err := rep.readJSON(w, r, &requestPayload)
	if err != nil {
		rep.errorJson(w, errors.New("invalid json"))
		return
	}

	err = rep.App.Models.User.Token.DeleteByToken(requestPayload.Token)
	if err != nil {
		rep.errorJson(w, errors.New("error deleting token"))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "logged out",
	}

	_ = rep.WriteJSON(w, http.StatusOK, payload)
}

func (rep *Repository) AllUsers(w http.ResponseWriter, r *http.Request) {
	var users data.User
	all, err := users.GetAll()
	if err != nil {
		rep.App.ErrorLogger.Println(err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    envelope{"users": all},
	}

	rep.WriteJSON(w, http.StatusOK, payload)
}

func (rep *Repository) EditUser(w http.ResponseWriter, r *http.Request) {
	var user data.User
	err := rep.readJSON(w, r, &user)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	if user.ID == 0 {
		if _, err := rep.App.Models.User.Insert(user); err != nil {
			rep.errorJson(w, err)
			return
		}
	} else {
		u, err := rep.App.Models.User.GetByID(user.ID)
		if err != nil {
			rep.errorJson(w, err)
			return
		}

		u.Email = user.Email
		u.FirstName = user.FirstName
		u.LastName = user.LastName
		u.Active = user.Active

		if err := u.Update(); err != nil {
			rep.errorJson(w, err)
			return
		}

		if user.Password != "" {
			err := u.ResetPassword(user.Password)
			if err != nil {
				rep.errorJson(w, err)
				return
			}
		}
	}

	payload := jsonResponse{
		Error:   false,
		Message: "User saved!",
	}

	_ = rep.WriteJSON(w, http.StatusAccepted, payload)
}

func (rep *Repository) GetUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	user, err := rep.App.Models.User.GetByID(userId)
	if err != nil {
		rep.errorJson(w, err)
		return
	}
	_ = rep.WriteJSON(w, http.StatusOK, user)
}

func (rep *Repository) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Id int `json:"id"`
	}

	err := rep.readJSON(w, r, &requestPayload)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	err = rep.App.Models.User.DeleteById(requestPayload.Id)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "User deleted!",
	}
	_ = rep.WriteJSON(w, http.StatusOK, payload)
}

func (rep *Repository) LogUserOutAdnSetInactive(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	user, err := rep.App.Models.User.GetByID(userId)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	user.Active = 0
	err = user.Update()
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	err = rep.App.Models.Token.DeleteTokensFroUser(userId)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "user logged out and set to inactive",
	}

	_ = rep.WriteJSON(w, http.StatusAccepted, payload)
}

func (rep *Repository) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Token string `json:"token"`
	}

	err := rep.readJSON(w, r, &requestPayload)
	if err != nil {
		rep.errorJson(w, err)
		return
	}

	valid := false
	valid, _ = rep.App.Models.Token.ValidToken(requestPayload.Token)

	payload := jsonResponse{
		Error: false,
		Data:  valid,
	}

	_ = rep.WriteJSON(w, http.StatusOK, payload)
}
