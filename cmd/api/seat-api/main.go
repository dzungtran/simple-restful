package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"simple-restful/cmd/api/seat-api/handlers"
	seat_svc "simple-restful/exmsgs/seat/services"
	"simple-restful/pkg/core/servehttp"
	"time"
)

func CreateSeatServiceClient() seat_svc.SeatServiceClient {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf("seat-svc:33033"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	return seat_svc.NewSeatServiceClient(conn)
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
		{
			Route:   "/seats/available",
			Method:  http.MethodGet,
			Handler: &handlers.GetAvailableSeatsHandler{
				SeatServiceClient: CreateSeatServiceClient(),
			},
		},
	}
}
