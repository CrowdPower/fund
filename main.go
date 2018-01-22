package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/crowdpower/fund/server"

	"github.com/spf13/viper"
)

func main() {
	getConfig()

	port := viper.GetString("server.port")
	cert := viper.GetString("server.cert")
	key := viper.GetString("server.key")
	r := server.Router()

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
