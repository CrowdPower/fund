package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/crowdpower/fund/controllers"
)

func Route(r *mux.Router, user controllers.UserController) {
	r.HandleFunc("/health", GetHealth).Methods(http.MethodGet)
	r.HandleFunc("/user", user.PostUser).Methods(http.MethodPost)
	r.HandleFunc("/user/{username}", user.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/user/{username}", user.PutUser).Methods(http.MethodPut)
	r.HandleFunc("/user/{username}", user.DeleteUser).Methods(http.MethodDelete)
}

func GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
