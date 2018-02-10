package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"

	"github.com/crowdpower/fund/models"
	"github.com/crowdpower/fund/storage"
	"github.com/crowdpower/fund/utils"
)

const (
	depositPageSize = 20
)

type DepositController interface {
	PostDeposit(w http.ResponseWriter, r *http.Request)
	GetDeposit(w http.ResponseWriter, r *http.Request)
	GetDeposits(w http.ResponseWriter, r *http.Request)
	GetDepositsSum(w http.ResponseWriter, r *http.Request)
}

type depositController struct {
	db storage.DB
}

func NewDepositController(db storage.DB) DepositController {
	return &depositController{db}
}

func (d *depositController) PostDeposit(w http.ResponseWriter, r *http.Request) {
	var deposit models.Deposit
	err := json.NewDecoder(r.Body).Decode(&deposit)
	if err != nil {
		log.Printf("could not unmarshal PostDeposit request body\n%v", err)
		utils.SendError(w, "Could not parse body as JSON", http.StatusBadRequest)
		return
	}

	if deposit.Amount <= 0 {
		utils.SendError(w, "Deposit amount must be greater than 0", http.StatusBadRequest)
		return
	}

	deposit.Username = mux.Vars(r)["username"]
	deposit.Id = uuid.NewV4().String()
	deposit.Time = time.Now().Format(storage.TimeFormat)

	err = d.db.CreateDeposit(&deposit)
	if err != nil {
		log.Printf("could not insert deposit %v into database\n%v", deposit, err)
		utils.SendError(w, "Error inserting deposit into database", http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, nil, http.StatusNoContent)
}

func (d *depositController) GetDeposit(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	id := r.URL.Query().Get("id")
	if id == "" {
		utils.SendError(w, "Parameter 'id' required", http.StatusBadRequest)
		return
	}

	deposit, err := d.db.GetDeposit(username, id)
	if storage.IsNotFound(err) {
		utils.SendError(w, fmt.Sprintf("Deposit %v not found for user %v", id, username), http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("could not get deposit %v for user %v from the database\n%v", id, username, err)
		utils.SendError(w, "Error getting deposit from database", http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, deposit, http.StatusOK)
}

func (d *depositController) GetDeposits(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	args := &models.DepositArgs{}
	err := utils.ParseArgs(r, args)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
	}

	if args.Count == 0 {
		args.Count = depositPageSize
	}

	deposits, err := d.db.GetDeposits(username, args)
	if err != nil {
		log.Printf("could not get deposits for user %v from the database\n%v", username, err)
		utils.SendError(w, "Error getting deposits from database", http.StatusInternalServerError)
		return
	}

	utils.SendPage(w, r, deposits, args.Offset + depositPageSize, depositPageSize, len(deposits) == args.Count)
}

func (d *depositController) GetDepositsSum(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	args := &models.DepositArgs{}
	err := utils.ParseArgs(r, args)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
	}

	sum, err := d.db.GetDepositsSum(username, args)
	if err != nil {
		log.Printf("could not get deposits sum for user %v from the database\n%v", username, err)
		utils.SendError(w, "Error getting deposits sum from database", http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, map[string]int{"sum": sum}, http.StatusOK)
}
