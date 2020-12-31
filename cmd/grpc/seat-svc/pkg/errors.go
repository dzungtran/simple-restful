package pkg

import "errors"

var (
	ErrorDuplicateSelectedSeats = errors.New("duplicate selected seats")
	ErrorSessionNotFound = errors.New("session not found")
	ErrorSelectedSeatsAreEmpty = errors.New("selected seats are empty")
)