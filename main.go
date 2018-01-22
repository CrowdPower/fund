package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/crowdpower/fund/controllers"
	"github.com/crowdpower/fund/server"
)

func main() {
	getConfig()

	port := viper.GetString("server.port")
	cert := viper.GetString("server.cert")
	key := viper.GetString("server.key")

	r := mux.NewRouter()
	uc := controllers.NewUserController()
	server.Route(r.PathPrefix("/v1").Subrouter(), uc)

	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServeTLS(":"+port, cert, key, r))
}

func getConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetDefault("server", map[string]string{
		"port": "8080",
		"cert": "server.crt",
		"key":  "server.key",
	})

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %v \n", err))
	}
}
