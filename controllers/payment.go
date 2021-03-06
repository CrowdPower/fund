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
	paymentPageSize = 20
)

type PaymentController interface {
	PostPayment(w http.ResponseWriter, r *http.Request)
	GetPayment(w http.ResponseWriter, r *http.Request)
	GetPayments(w http.ResponseWriter, r *http.Request)
	GetPaymentsSum(w http.ResponseWriter, r *http.Request)
}

type paymentController struct {
	db storage.DB
}

func NewPaymentController(db storage.DB) PaymentController {
	return &paymentController{db}
}

func (d *paymentController) PostPayment(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment
	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		log.Printf("could not unmarshal PostPayment request body\n%v", err)
		utils.SendError(w, "Could not parse body as JSON", http.StatusBadRequest)
		return
	}

	if payment.Amount <= 0 {
		utils.SendError(w, "Payment amount must be greater than 0", http.StatusBadRequest)
		return
	}

	if payment.Url == "" {
		utils.SendError(w, "Payment url cannot be empty", http.StatusBadRequest)
		return
	}

	payment.Username = mux.Vars(r)["username"]
	payment.Id = uuid.NewV4().String()
	payment.Time = time.Now().Format(storage.TimeFormat)

	err = d.db.CreatePayment(&payment)
	if err != nil {
		if storage.IsInsufficientFunds(err) {
			utils.SendError(w, "Insufficient funds", http.StatusBadRequest)
			return
		}
		log.Printf("could not insert payment %v into database\n%v", payment, err)
		utils.SendError(w, "Error inserting payment into database", http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, nil, http.StatusNoContent)
}

func (d *paymentController) GetPayment(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	id := r.URL.Query().Get("id")
	if id == "" {
		utils.SendError(w, "Parameter 'id' required", http.StatusBadRequest)
		return
	}

	payment, err := d.db.GetPayment(username, id)
	if storage.IsNotFound(err) {
		utils.SendError(w, fmt.Sprintf("Payment %v not found for user %v", id, username), http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("could not get payment %v for user %v from the database\n%v", id, username, err)
		utils.SendError(w, "Error getting payment from database", http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, payment, http.StatusOK)
}

func (d *paymentController) GetPayments(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	args := &models.PaymentArgs{}
	err := utils.ParseArgs(r, args)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if args.Count == 0 {
		args.Count = paymentPageSize
	}

	payments, err := d.db.GetPayments(username, args)
	if err != nil {
		log.Printf("could not get payments for user %v from the database\n%v", username, err)
		utils.SendError(w, "Error getting payments from database", http.StatusInternalServerError)
		return
	}

	utils.SendPage(w, r, payments, args.Offset+paymentPageSize, paymentPageSize, len(payments) == args.Count)
}

func (d *paymentController) GetPaymentsSum(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	args := &models.PaymentArgs{}
	err := utils.ParseArgs(r, args)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sum, err := d.db.GetPaymentsSum(username, args)
	if err != nil {
		log.Printf("could not get payments sum for user %v from the database\n%v", username, err)
		utils.SendError(w, "Error getting payments sum from database", http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, map[string]int{"sum": sum}, http.StatusOK)
}
