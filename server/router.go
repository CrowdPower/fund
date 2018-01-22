package server

import (
    "net/http"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
    r := mux.NewRouter()

    r.HandleFunc("/health", GetHealth).Methods(http.MethodGet)

    return r
}

func GetHealth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}
