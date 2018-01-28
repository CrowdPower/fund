package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/crowdpower/fund/controllers"
)

func Route(
	r *mux.Router,
	user controllers.UserController,
	auth controllers.AuthController,
	deposit controllers.DepositController,
	payment controllers.PaymentController) {

	r.HandleFunc("/health",
		GetHealth,
	).Methods(http.MethodGet)

	r.HandleFunc("/user",
		user.PostUser,
	).Methods(http.MethodPost)
	r.HandleFunc("/user/{username}",
		auth.Wrapper(controllers.AccessTokenType, user.GetUser),
	).Methods(http.MethodGet)
	r.HandleFunc("/user/{username}",
		auth.Wrapper(controllers.AccessTokenType, user.PutUser),
	).Methods(http.MethodPut)
	r.HandleFunc("/user/{username}",
		auth.Wrapper(controllers.AccessTokenType, user.DeleteUser),
	).Methods(http.MethodDelete)

	r.HandleFunc("/user/{username}/authorize",
		auth.GetRefreshToken,
	).Methods(http.MethodGet)
	r.HandleFunc("/user/{username}/token",
		auth.Wrapper(controllers.RefreshTokenType, auth.GetAuthToken),
	).Methods(http.MethodGet)

	r.HandleFunc("/user/{username}/deposit",
		auth.Wrapper(controllers.AccessTokenType, deposit.PostDeposit),
	).Methods(http.MethodPost)
	r.HandleFunc("/user/{username}/deposit",
		auth.Wrapper(controllers.AccessTokenType, deposit.GetDeposit),
	).Methods(http.MethodGet)
	r.HandleFunc("/user/{username}/deposits",
		auth.Wrapper(controllers.AccessTokenType, deposit.GetDeposits),
	).Methods(http.MethodGet)
	r.HandleFunc("/user/{username}/deposits/sum",
		auth.Wrapper(controllers.AccessTokenType, deposit.GetDepositsSum),
	).Methods(http.MethodGet)

	r.HandleFunc("/user/{username}/payment",
		auth.Wrapper(controllers.AccessTokenType, payment.PostPayment),
	).Methods(http.MethodPost)
	r.HandleFunc("/user/{username}/payment",
		auth.Wrapper(controllers.AccessTokenType, payment.GetPayment),
	).Methods(http.MethodGet)
	r.HandleFunc("/user/{username}/payments",
		auth.Wrapper(controllers.AccessTokenType, payment.GetPayments),
	).Methods(http.MethodGet)
	r.HandleFunc("/user/{username}/payments/sum",
		auth.Wrapper(controllers.AccessTokenType, payment.GetPaymentsSum),
	).Methods(http.MethodGet)
}

func GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
