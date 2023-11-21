package main

import (
	"fmt"
	"ll_test/app/config"
	"ll_test/app/connector"
	"log"
	"net/http"

	"github.com/spf13/viper"

	"ll_test/app/logger"
)

// - VIPER lib for config.json reading when init moment
func init() {
	err := config.LoadConfiguration("config")

	//- Uncomment this line for VSCode debugging
	// err := config.LoadConfigurationForDebugging()
	if err != nil {
		log.Fatal(err, "Cannot load config file with Viper")
	}
}

func main() {

	//- logger initialize
	log := logger.NewLogrusLogger()

	//- go-chi implementation
	r := connector.SetupRouter()

	// Start server
	defaultFields := logger.Fields{
		"service":  "littlelives",
		"hostname": "localhost",
	}
	log.Fields(defaultFields).Infof("Program is running. Access http://localhost:%s", viper.GetString("SERVER.PORT"))

	//- router
	fmt.Println("Server start at port " + ":" + viper.GetString("SERVER.PORT"))
	err := http.ListenAndServe(":"+viper.GetString("SERVER.PORT"), r)
	if err != nil {
		log.Fatal(err, "error on serve server")
	}
}
