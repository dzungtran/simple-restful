package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"simple-restful/cmd/grpc/seat-svc/booking"
	"simple-restful/cmd/grpc/seat-svc/pkg"
	seatUtils "simple-restful/cmd/grpc/seat-svc/pkg/utils"
	seat_models "simple-restful/exmsgs/seat/models"
	seat_svc "simple-restful/exmsgs/seat/services"
	"simple-restful/pkg/dtos"
)

// Server represents the gRPC server
type server struct {
	CinemaSettings *dtos.CinemaSettings
	RedisClient    *redis.Client
	BookingRepo    booking.BookingRepository
}

func (s *server) GetAvailableSeats(ctx context.Context, in *seat_svc.GetAvailableSeatsRequest) (*seat_svc.SeatsResult, error) {
	result := make([]*seat_svc.SetOfSeats, 0)
	total := int64(0)

	// Init available seats
	seats := s.initCinemaSeats()

	// Unavailable seats and ignore selected seats of current session
	mapUnavailableSetOfSeats, err := s.getUnavailableSetOfSeats(in.CurrentSession)
	if err != nil {
		return nil, err
	}

	// Remove unavailable seats
	seatUtils.RemoveUnavailableSetOfSeats(seats, mapUnavailableSetOfSeats)
	for _, set := range seats {
		if len(set.Cols) == 0 {
			continue
		}
		result = append(result, set)
		total += int64(len(set.Cols))
	}

	return &seat_svc.SeatsResult{
		List:           result,
		TotalAvailable: total,
	}, nil
}

// BookSeats
func (s *server) BookSeats(ctx context.Context, in *seat_svc.SeatsRequest) (*seat_svc.SeatAvailabilityResult, error) {
	err := validateSelectedSeats(in.List)
	if err != nil {
		return nil, err
	}

	sessionKey := in.Session
	if in.Session == "" {
		id, err := uuid.NewUUID()
		if err != nil {
			log.Fatalf("uuid.NewUUID() failed with %s\n", err)
		}
		sessionKey = id.String()
	}

	seatsToCheck := in.List

	// append all seat at around a bunch of seats with a max distance
	if s.CinemaSettings.MaxDistance > 0 {
		seatsInRange := seatUtils.GetSeatsInRange(
			seatsToCheck,
			s.CinemaSettings.NumRow,
			s.CinemaSettings.NumCol,
			s.CinemaSettings.MaxDistance,
		)
		if len(seatsInRange) > 0 {
			seatsToCheck = append(seatsToCheck, seatsInRange...)
		}
	}

	// Check seats are available or not: booked or pending by other session aren't allowed
	bookedSeats, err := s.BookingRepo.CheckBookedSeats(seatsToCheck)
	if err != nil {
		return nil, err
	}

	if bookedSeats != nil && len(bookedSeats) > 0 {
		return &seat_svc.SeatAvailabilityResult{
			Unavailable: bookedSeats,
			Session:     sessionKey,
			Status:      "fail",
		}, nil
	}

	pendingSeats, err := s.BookingRepo.CheckPendingSeats(seatsToCheck, sessionKey)
	if err != nil {
		return nil, err
	}

	if pendingSeats != nil && len(pendingSeats) > 0 {
		return &seat_svc.SeatAvailabilityResult{
			Unavailable: pendingSeats,
			Session:     sessionKey,
			Status:      "fail",
		}, nil
	}

	// Update pending seats by current session
	err = s.BookingRepo.UpdatePendingSession(sessionKey, in.List)
	if err != nil {
		return nil, err
	}

	if in.Confirm {
		// Move seats to booked list
		err = s.BookingRepo.BookSeats(sessionKey)
		if err != nil {
			return nil, err
		}
	}

	return &seat_svc.SeatAvailabilityResult{
		Session:   sessionKey,
		Available: in.List,
		Status:    "success",
	}, nil
}

func (s *server) Reset(ctx context.Context, in *seat_svc.EmptyRequest) (*seat_svc.ResultResponse, error) {
	err := s.BookingRepo.Reset()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *server) getUnavailableSetOfSeats(ignoreSession string) (map[int64]*seat_svc.SetOfSeats, error) {

	unavailableSetOfSeats := make([]*seat_svc.SetOfSeats, 0)

	// Check list booked seat
	bookedSeats, err := s.BookingRepo.GetBookedSeats()
	if err != nil {
		return nil, err
	}

	if bookedSeats != nil && len(bookedSeats) > 0 {
		unavailableSetOfSeats = append(unavailableSetOfSeats, bookedSeats...)
	}

	// Get list pending seat, which seats are locked by others (sessions)
	// and ignore selected seats by current session
	ignoreSessions := make([]string, 0)
	if ignoreSession != "" {
		ignoreSessions = append(ignoreSessions, ignoreSession)
	}

	pendingSeats, err := s.BookingRepo.GetPendingSeats(ignoreSessions)
	if err != nil {
		return nil, err
	}

	if pendingSeats != nil && len(pendingSeats) > 0 {
		unavailableSetOfSeats = append(unavailableSetOfSeats, pendingSeats...)
	}

	// Convert to map and remove duplicate cols
	return seatUtils.MapSetOfSeats(unavailableSetOfSeats...), nil
}

// Init available seat of cinema
func (s *server) initCinemaSeats() map[int64]*seat_svc.SetOfSeats {
	seats := make(map[int64]*seat_svc.SetOfSeats)

	for r := int64(0); r < s.CinemaSettings.NumRow; r++ {
		cols := make([]int64, 0)
		for c := int64(0); c < s.CinemaSettings.NumCol; c++ {
			cols = append(cols, c)
		}
		seats[r] = &seat_svc.SetOfSeats{
			Row:  r,
			Cols: cols,
		}
	}
	return seats
}

func validateSelectedSeats(seats []*seat_models.Seat) error {

	if seats == nil || len(seats) == 0 {
		return pkg.ErrorSelectedSeatsAreEmpty
	}

	// index: row+col, value: bool
	mapSeats := make(map[string]bool)
	for _, s := range seats {
		if _, ok := mapSeats[seatUtils.ParseSeatToKey(s)]; !ok {
			mapSeats[seatUtils.ParseSeatToKey(s)] = true
		} else {
			return pkg.ErrorDuplicateSelectedSeats
		}
	}

	return nil
}
