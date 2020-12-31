package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	bookingRedis "simple-restful/cmd/grpc/seat-svc/booking/redis"
	seat_models "simple-restful/exmsgs/seat/models"
	seat_svc "simple-restful/exmsgs/seat/services"
	"simple-restful/pkg/dtos"
	"testing"
)

func initTestClient() *server {
	viper.SetConfigName("app")
	viper.SetConfigType("toml")
	viper.SetConfigFile("./conf.toml")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.connection"),
		DB:   viper.GetInt("redis.db"), // use default DB
	})
	settings = &dtos.CinemaSettings{
		NumRow:      viper.GetInt64("cinema.row"),
		NumCol:      viper.GetInt64("cinema.col"),
		MaxDistance: viper.GetInt64("general.distance"),
	}

	return &server{
		CinemaSettings: settings,
		RedisClient:    rdb,
		BookingRepo:    bookingRedis.NewBookingRedis(rdb),
	}
}

func TestServer_BookingSeats(t *testing.T) {
	seatClient := initTestClient()
	ctx := context.Background()
	total := seatClient.CinemaSettings.NumCol * seatClient.CinemaSettings.NumRow

	t.Logf("Config: %#v", seatClient.CinemaSettings)

	t.Run("Get all available seat", func(t *testing.T) {
		seats, err := seatClient.GetAvailableSeats(ctx, &seat_svc.GetAvailableSeatsRequest{})
		assert.Nil(t, err)
		assert.Equal(t, seats.TotalAvailable, total)

		// logging
		data, _ := json.Marshal(seats)
		t.Logf("SEAT: %v", string(data))
	})

	t.Run("Get available seats with 1 valid booking", func(t *testing.T) {
		// book some seats
		rs, err := seatClient.BookSeats(ctx, &seat_svc.SeatsRequest{
			List: []*seat_models.Seat{
				{
					Row: 1,
					Col: 1,
				},
			},
			Confirm: true,
		})
		assert.Nil(t, err)
		assert.Equal(t, rs.Status, "success")

		seats, err := seatClient.GetAvailableSeats(ctx, &seat_svc.GetAvailableSeatsRequest{})
		assert.Nil(t, err)
		assert.Equal(t, total-1, seats.TotalAvailable)

		// logging
		data, _ := json.Marshal(seats)
		t.Logf("SEAT: %v", string(data))
	})

	t.Run("Get available seats with 2 invalid bookings", func(t *testing.T) {
		// book some seats
		rs, err := seatClient.BookSeats(ctx, &seat_svc.SeatsRequest{
			List: []*seat_models.Seat{
				{
					Row: 2,
					Col: 1,
				},
				{
					Row: 3,
					Col: 1,
				},
			},
			Confirm: true,
		})
		assert.Nil(t, err)
		assert.Equal(t, rs.Status, "fail")
		data, _ := json.Marshal(rs)
		t.Logf("BOOK SEAT: %v", string(data))

		seats, err := seatClient.GetAvailableSeats(ctx, &seat_svc.GetAvailableSeatsRequest{})
		assert.Nil(t, err)
		assert.Equal(t, total-1, seats.TotalAvailable)

		// logging
		data, _ = json.Marshal(seats)
		t.Logf("SEAT: %v", string(data))
	})

	t.Run("Get available seats with more 2 valid bookings", func(t *testing.T) {
		// book some seats
		rs, err := seatClient.BookSeats(ctx, &seat_svc.SeatsRequest{
			List: []*seat_models.Seat{
				{
					Row: 4,
					Col: 2,
				},
				{
					Row: 4,
					Col: 4,
				},
			},
			Confirm: true,
		})
		assert.Nil(t, err)
		assert.Equal(t, rs.Status, "success")
		data, _ := json.Marshal(rs)
		t.Logf("BOOK SEAT: %v", string(data))

		seats, err := seatClient.GetAvailableSeats(ctx, &seat_svc.GetAvailableSeatsRequest{})
		assert.Nil(t, err)
		assert.Equal(t, total-3, seats.TotalAvailable)

		// logging
		data, _ = json.Marshal(seats)
		t.Logf("SEAT: %v", string(data))
	})

	t.Run("Run book seats with duplicate", func(t *testing.T) {
		// book some seats
		rs, err := seatClient.BookSeats(ctx, &seat_svc.SeatsRequest{
			List: []*seat_models.Seat{
				{
					Row: 1,
					Col: 2,
				},
				{
					Row: 1,
					Col: 2,
				},
			},
			Confirm: true,
		})
		assert.NotNil(t, err)
		t.Logf("BOOK SEAT: %v", rs)
	})

	// clean test data
	seatClient.Reset(ctx, &seat_svc.EmptyRequest{})
}
