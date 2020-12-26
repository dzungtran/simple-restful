package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	seat_svc "simple-restful/exmsgs/seat/services"
	"simple-restful/pkg/dtos"
)

// Server represents the gRPC server
type server struct {
	CinemaSettings *dtos.CinemaSettings
	RedisClient    *redis.Client
}

func (s *server) GetAvailableSeats(ctx context.Context, in *seat_svc.EmptyRequest) (*seat_svc.SeatsResult, error) {
	return nil, nil
}

// BookSeats
func (s *server) BookSeats(ctx context.Context, in *seat_svc.SeatsRequest) (*seat_svc.SeatAvailabilityResult, error) {
	return nil, nil
}

func (s *server) Reset(ctx context.Context, in *seat_svc.EmptyRequest) (*seat_svc.ResultResponse, error) {
	return nil, nil
}
