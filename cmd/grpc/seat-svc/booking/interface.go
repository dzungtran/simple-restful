package booking

import (
	seat_models "simple-restful/exmsgs/seat/models"
	seat_svc "simple-restful/exmsgs/seat/services"
)

type BookingRepository interface {
	BookSeats(sessionKey string) error

	GetBookedSeats() ([]*seat_svc.SetOfSeats, error)
	GetPendingSeats(ignoreSession []string) ([]*seat_svc.SetOfSeats, error)

	UpdatePendingSession(sessionKey string, seats []*seat_models.Seat) error

	CheckBookedSeats(seats []*seat_models.Seat) (unavailableSeats []*seat_models.Seat, err error)
	CheckPendingSeats(seats []*seat_models.Seat, sessionKey string) (unavailableSeats []*seat_models.Seat, err error)

	Reset() error
}
