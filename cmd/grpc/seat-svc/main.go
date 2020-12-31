package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	bookingRedis "simple-restful/cmd/grpc/seat-svc/booking/redis"
	seat_svc "simple-restful/exmsgs/seat/services"
	"simple-restful/pkg/dtos"
)

var rdb *redis.Client
var settings *dtos.CinemaSettings

// to init db connection or api configs
func initConfig() {
	configFilePath := os.Getenv("CONFIG_PATH")
	log.Print("FILE: ", configFilePath)

	//viper.SetConfigName("app")
	viper.SetConfigType("toml")
	viper.SetConfigFile(configFilePath)
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	log.Print("Connecting redis: ", viper.GetString("redis.connection"))
	rdb = redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.connection"),
		DB:   viper.GetInt("redis.db"), // use default DB
	})

	settings = &dtos.CinemaSettings{
		NumRow:      viper.GetInt64("cinema.row"),
		NumCol:      viper.GetInt64("cinema.col"),
		MaxDistance: viper.GetInt64("general.distance"),
		SessionTTL:  viper.GetInt64("general.session_ttl"),
	}
}

func main() {
	initConfig()
	localPort := os.Getenv("APP_PORT")
	log.Println(fmt.Sprintf("Starting GPRC server with port :%v", localPort))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", localPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	s := server{
		CinemaSettings: settings,
		RedisClient:    rdb,
		BookingRepo:    bookingRedis.NewBookingRedis(rdb),
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	seat_svc.RegisterSeatServiceServer(grpcServer, &s)
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
