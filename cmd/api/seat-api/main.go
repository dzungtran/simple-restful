package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"simple-restful/cmd/api/seat-api/handlers"
	"simple-restful/pkg/core/servehttp"
	"time"
)

// to init db connection or api configs
func initConfig() {
	configFilePath := flag.String("cf", "./conf/app.yaml", "Path to config file")
	log.Print("FILE: ", *configFilePath)

	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*configFilePath)
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func main() {
	//initConfig()

	// Init App server
	appServer := &servehttp.AppServer{}
	appServer.Init()
	for _, handler := range getListHandler() {
		appServer.RegisterHandler(handler.Method, handler.Route, handler.Handler)
	}

	localPort := os.Getenv("APP_PORT")
	srv := &http.Server{
		Handler:      appServer.GetRouter(),
		Addr:         fmt.Sprintf(":%v", localPort),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println(fmt.Sprintf("Starting API server with port :%v", localPort))
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	servehttp.WaitForShutdown(srv)
}

// Defines list handler to serve requests
func getListHandler() []servehttp.AppHandler {
	return []servehttp.AppHandler{
		{
			Route:   "/",
			Method:  http.MethodGet,
			Handler: &handlers.GetHelloHandler{},
		},
	}
}
