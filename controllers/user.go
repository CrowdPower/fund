package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/crowdpower/fund/models"
	"github.com/crowdpower/fund/utils"
)

type UserController interface {
	PostUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	PutUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

func NewUserController() UserController {
	return &userController{}
}

type userController struct {
}

func (u *userController) PostUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("could not unmarshal PostUser request body\n%v", err)
		utils.SendError(w, "Could not parse body as JSON", http.StatusInternalServerError)
		return
	}

	log.Printf("POST user %v", user)

	utils.SendSuccess(w, nil, http.StatusNoContent)
}

func (u *userController) GetUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	if username == "" {
		utils.SendError(w, "Username required", http.StatusBadRequest)
		return
	}
	log.Printf("GET user %v", username)

	utils.SendSuccess(w, nil, http.StatusNoContent)
}

func (u *userController) PutUser(w http.ResponseWriter, r *http.Request) {
	utils.SendSuccess(w, nil, http.StatusNoContent)
}

func (u *userController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	utils.SendSuccess(w, nil, http.StatusNoContent)
}
