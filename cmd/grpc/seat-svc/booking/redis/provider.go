package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"simple-restful/cmd/grpc/seat-svc/booking"
	"simple-restful/cmd/grpc/seat-svc/pkg"
	seatUtils "simple-restful/cmd/grpc/seat-svc/pkg/utils"
	seat_models "simple-restful/exmsgs/seat/models"
	seat_svc "simple-restful/exmsgs/seat/services"
	"simple-restful/pkg/core/utils"
	"strings"
)

const (
	BookedSeatsKey  = "booked"
	PendingSeatsKey = "pending"
	SessionKey      = "session"
)

func NewBookingRedis(rdb *redis.Client) booking.BookingRepository {
	return &BookingRedis{RedisClient: rdb}
}

type BookingRedis struct {
	RedisClient *redis.Client
}

func (b *BookingRedis) BookSeats(sessionKey string) error {
	ctx := context.Background()

	// Get pendingSeats from session
	pendingSeatKeys, err := b.getPendingSeatKeysBySessions([]string{sessionKey})
	if err != nil {
		return err
	}

	if pendingSeatKeys != nil && len(pendingSeatKeys) > 0 {
		for _, key := range pendingSeatKeys {
			// Create booked seat
			b.RedisClient.HSet(ctx, BookedSeatsKey, key, sessionKey)
		}

		// Remove pending seats
		b.RedisClient.HDel(ctx, PendingSeatsKey, pendingSeatKeys...)
	}

	// Remove Pending session
	b.RedisClient.HDel(ctx, SessionKey, sessionKey)
	return nil
}

func (b *BookingRedis) GetBookedSeats() (seats []*seat_svc.SetOfSeats, err error) {
	seats = make([]*seat_svc.SetOfSeats, 0)
	ctx := context.Background()
	// Get all fields of booked
	rs := b.RedisClient.HKeys(ctx, BookedSeatsKey)
	if rs.Err() != nil {
		return seats, rs.Err()
	}

	data, err := rs.Result()
	if err != nil {
		if err == redis.Nil {
			return seats, nil
		}
		return
	}

	mapSet := make(map[int64]*seat_svc.SetOfSeats, 0)
	for _, k := range data {
		pos := strings.Split(k, ":")
		if len(pos) != 2 {
			continue
		}

		// Parse row and col from field name
		row, col := cast.ToInt64(pos[0]), cast.ToInt64(pos[1])
		if _, ok := mapSet[row]; !ok {
			mapSet[row] = &seat_svc.SetOfSeats{
				Row:  row,
				Cols: []int64{col},
			}
		} else {
			mapSet[row].Cols = append(mapSet[row].Cols, col)
		}
	}

	// convert map to slice
	if len(mapSet) > 0 {
		for _, set := range mapSet {
			seats = append(seats, set)
		}
	}

	return
}

func (b *BookingRedis) GetPendingSeats(ignoreSessions []string) ([]*seat_svc.SetOfSeats, error) {
	// Get pending seat of by sessions
	seats := make([]*seat_svc.SetOfSeats, 0)
	ctx := context.Background()
	pendingSeatsKeys, err := b.getPendingSeatKeysBySessions(ignoreSessions)
	if err != nil {
		return nil, err
	}

	// Get all pending seat
	rs := b.RedisClient.HKeys(ctx, PendingSeatsKey)
	if rs.Err() != nil {
		return seats, rs.Err()
	}

	data, err := rs.Result()
	if err != nil {
		if err == redis.Nil {
			return seats, nil
		}
		return nil, err
	}

	mapSet := make(map[int64]*seat_svc.SetOfSeats, 0)
	for _, k := range data {

		// ignore current session
		if pendingSeatsKeys != nil && utils.IsStringSliceContains(pendingSeatsKeys, k) {
			continue
		}

		pos := strings.Split(k, ":")
		if len(pos) != 2 {
			continue
		}

		// Parse row and col from field name
		row, col := cast.ToInt64(pos[0]), cast.ToInt64(pos[1])
		if _, ok := mapSet[row]; !ok {
			mapSet[row] = &seat_svc.SetOfSeats{
				Row:  row,
				Cols: []int64{col},
			}
		} else {
			mapSet[row].Cols = append(mapSet[row].Cols, col)
		}
	}

	// convert map to slice
	if len(mapSet) > 0 {
		for _, set := range mapSet {
			seats = append(seats, set)
		}
	}

	return nil, nil
}

func (b *BookingRedis) RemovePendingSession(sessionKey string) error {
	ctx := context.Background()
	pendingSeatKeys, err := b.getPendingSeatKeysBySessions([]string{sessionKey})
	if err != nil {
		return err
	}

	if pendingSeatKeys != nil && len(pendingSeatKeys) > 0 {
		b.RedisClient.HDel(ctx, PendingSeatsKey, pendingSeatKeys...)
	}

	// Remove session
	b.RedisClient.HDel(ctx, SessionKey, sessionKey)
	return nil
}

func (b *BookingRedis) UpdatePendingSession(sessionKey string, seats []*seat_models.Seat) error {
	ctx := context.Background()
	if seats == nil || len(seats) == 0 {
		return nil
	}

	// Clear pending session
	err := b.RemovePendingSession(sessionKey)
	if err != nil {
		return err
	}

	for _, s := range seats {
		b.RedisClient.HSet(ctx, PendingSeatsKey, seatUtils.ParseSeatToKey(s), sessionKey)
	}

	// re-create session
	jsonData, err := json.Marshal(seats)
	if err != nil {
		return err
	}

	b.RedisClient.HSet(ctx, SessionKey, sessionKey, string(jsonData))

	return nil
}

func (b *BookingRedis) CheckBookedSeats(seats []*seat_models.Seat) (bookedSeats []*seat_models.Seat, err error) {
	bookedSeats = make([]*seat_models.Seat, 0)
	ctx := context.Background()

	for _, s := range seats {
		rs := b.RedisClient.HGet(ctx, BookedSeatsKey, seatUtils.ParseSeatToKey(s))
		st, _ := rs.Result()
		if st != "" {
			bookedSeats = append(bookedSeats, s)
		}
	}

	return bookedSeats, nil
}

func (b *BookingRedis) CheckPendingSeats(seats []*seat_models.Seat, ignoredSessionKey string) (pendingSeats []*seat_models.Seat, err error) {
	pendingSeats = make([]*seat_models.Seat, 0)
	ctx := context.Background()

	for _, s := range seats {
		rs := b.RedisClient.HGet(ctx, PendingSeatsKey, seatUtils.ParseSeatToKey(s))
		st, _ := rs.Result()
		if st != "" {
			if ignoredSessionKey != "" && st == ignoredSessionKey {
				continue
			}
			pendingSeats = append(pendingSeats, s)
		}
	}

	return pendingSeats, nil
}

func (b *BookingRedis) Reset() error {
	rs := b.RedisClient.FlushDB(context.Background())
	return rs.Err()
}

func (b *BookingRedis) getPendingSeatKeysBySessions(sessions []string) ([]string, error) {
	seatKeys := make([]string, 0)
	ctx := context.Background()

	if len(sessions) == 0 {
		return seatKeys, nil
	}

	rs := b.RedisClient.HMGet(ctx, SessionKey, sessions...)
	if rs.Err() != nil {
		return nil, rs.Err()
	}

	data, err := rs.Result()
	if err != nil {
		if err == redis.Nil {
			return nil, pkg.ErrorSessionNotFound
		}
		return nil, err
	}

	for _, s := range data {
		var pendingSeats []*seat_models.Seat
		if err = json.Unmarshal([]byte(cast.ToString(s)), &pendingSeats); err != nil {
			continue
		}

		if pendingSeats == nil || len(pendingSeats) == 0 {
			continue
		}

		for _, s := range pendingSeats {
			seatKeys = append(seatKeys, seatUtils.ParseSeatToKey(s))
		}
	}

	return seatKeys, nil
}
