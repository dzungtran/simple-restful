package main

import (
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"simple-restful/cmd/apis/user-api/handlers"
	"simple-restful/pkg/core/servehttp"
	"time"
)

const LocalPort = 8080

var db *gorm.DB
var err error

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

	mysqlConnStr := viper.GetString("mysql_conn")
	log.Print("Mysql Connection: ", mysqlConnStr)

	db, err = gorm.Open("mysql", mysqlConnStr)
	if err != nil {
		log.Fatalf("Cannot connect to db, details: %v", err.Error())
	}
}

func main() {
	initConfig()
	if db != nil {
		defer db.Close()
	}

	// Init App server
	appServer := &servehttp.AppServer{}
	appServer.Init()
	for _, handler := range getListHandler() {
		appServer.RegisterHandler(handler.Method, handler.Route, handler.Handler)
	}

	localPort := viper.GetInt("http_port")
	srv := &http.Server{
		Handler:      appServer.GetRouter(),
		Addr:         fmt.Sprintf(":%d", localPort),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println(fmt.Sprintf("Starting server with port :%d", localPort))
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

		// Get user transactions
		{
			Route:  "/api/users/{user_id:[0-9]+}/transactions",
			Method: http.MethodGet,
			Handler: &handlers.GetUserTransactionsHandler{
				DB: db,
			},
		},

		// Create user transaction
		{
			Route:  "/api/users/{user_id:[0-9]+}/transactions",
			Method: http.MethodPost,
			Handler: &handlers.CreateTransactionHandler{
				DB: db,
			},
		},

		// SHOULD BE avoid method `PUT` for update transaction and method `DELETE` for delete transaction.
		// Because they will change user balance without a reason and you can't see them on transaction timeline.
	}
}
