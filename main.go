package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/crowdpower/fund/controllers"
	"github.com/crowdpower/fund/server"
	"github.com/crowdpower/fund/storage"
)

func main() {
	getConfig()

	databaseType := viper.GetString("database.type")
	databasePath := viper.GetString("database.path")

	port := viper.GetString("server.port")
	cert := viper.GetString("server.cert")
	key := viper.GetString("server.key")

	db, err := storage.GetDB(databaseType, databasePath)
	if err != nil {
		log.Fatal("error connecting to database\n%v", err)
	}

	r := mux.NewRouter()
	uc := controllers.NewUserController(db)
	server.Route(r.PathPrefix("/v1").Subrouter(), uc)

	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServeTLS(":"+port, cert, key, r))
}

func getConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %v \n", err))
	}
}
