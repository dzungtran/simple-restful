package utils

import (
	"fmt"
	"math"
	seat_models "simple-restful/exmsgs/seat/models"
	seat_svc "simple-restful/exmsgs/seat/services"
	"simple-restful/pkg/core/utils"
)

func CalculateSeatDistance(seat *seat_models.Seat, row, col int64) int64 {
	distance := int64(math.Abs(float64(seat.Col-col)) + math.Abs(float64(seat.Row-row)))
	return distance
}

func GetMaxMinRange(pos int64, len int64, distance int64) (minRange int64, maxRange int64) {
	minRange = int64(math.Max(0, float64(pos-distance)))
	maxRange = int64(math.Min(float64(len-1), float64(pos+distance)))
	return
}

// Just get items in range min - max
func GetItemsInRange(min int64, max int64, ignores []int64) []int64 {
	items := make([]int64, 0)
	for i := min; i <= max; i++ {
		if utils.IsInt64SliceContains(ignores, i) {
			continue
		}
		items = append(items, i)
	}
	return items
}

// Return map of seats with index is row and value is slice of cols
// The seat must be in range
// @params seat seat_models.Seat Current seat position
// @params numRow int64  Number row of table
// @params numCol int64  Number col of table
func GetAllSeatInRange(seat *seat_models.Seat, numRow int64, numCol int64, maxRange int64) map[int64][]int64 {
	seats := make(map[int64][]int64)
	if maxRange == 0 || seat == nil || numRow < 0 || numCol < 0 {
		return nil
	}

	// Out of range
	if seat.Col >= numCol || seat.Row >= numRow {
		return nil
	}

	availableRows := make([]int64, 0)
	availableCols := make([]int64, 0)

	minCol, maxCol := GetMaxMinRange(seat.Col, numCol, maxRange)
	minRow, maxRow := GetMaxMinRange(seat.Row, numRow, maxRange)

	availableRows = GetItemsInRange(minRow, maxRow, []int64{})
	availableCols = GetItemsInRange(minCol, maxCol, []int64{})

	for _, r := range availableRows {
		for _, c := range availableCols {
			// calculate distance, if not in range, continue
			if CalculateSeatDistance(seat, r, c) > maxRange {
				continue
			}
			// ignore current seat
			if c == seat.Col && r == seat.Row {
				continue
			}

			if _, ok := seats[r]; !ok {
				seats[r] = make([]int64, 0)
			}
			seats[r] = append(seats[r], c)
		}
	}

	if len(seats) == 0 {
		return nil
	}

	return seats
}

// Get all seat around a bunch of seats with max distance
func GetSeatsInRange(seats []*seat_models.Seat, numRow, numCol, maxDistance int64) []*seat_models.Seat {
	mapRowSeats := make(map[int64][]int64)
	inRangeSeats := make([]*seat_models.Seat, 0)

	for _, st := range seats {
		tmpMap := GetAllSeatInRange(st, numRow, numCol, maxDistance)
		if tmpMap == nil || len(tmpMap) == 0 {
			continue
		}

		for r, sliceCol := range tmpMap {
			if _, ok := mapRowSeats[r]; ok {
				mapRowSeats[r] = sliceCol
			} else {
				// merge slice of col
				mapRowSeats[r] = utils.Int64SliceUnique(append(mapRowSeats[r], sliceCol...))
			}
		}
	}

	// convert map seats to slice seat
	if len(mapRowSeats) > 0 {
		for r, cols := range mapRowSeats {
			if len(cols) == 0 {
				continue
			}
			for _, c := range cols {
				inRangeSeats = append(inRangeSeats, &seat_models.Seat{
					Row: r,
					Col: c,
				})
			}
		}
	}

	return inRangeSeats
}

// convert slice set to map set, and merge them into one
func MapSetOfSeats(sets ...*seat_svc.SetOfSeats) map[int64]*seat_svc.SetOfSeats {
	mapSets := make(map[int64]*seat_svc.SetOfSeats)

	for _, s := range sets {
		if s == nil || s.Cols == nil || len(s.Cols) == 0 {
			continue
		}
		if _, ok := mapSets[s.Row]; !ok {
			mapSets[s.Row] = &seat_svc.SetOfSeats{
				Row:  s.Row,
				Cols: s.Cols,
			}
			continue
		}
		mapSets[s.Row].Cols = utils.Int64SliceUnique(append(mapSets[s.Row].Cols, s.Cols...))
	}
	return mapSets
}

func RemoveUnavailableSetOfSeats(
	list map[int64]*seat_svc.SetOfSeats,
	unavailable map[int64]*seat_svc.SetOfSeats,
) map[int64]*seat_svc.SetOfSeats {

	for _, s := range unavailable {
		if _, ok := list[s.Row]; !ok {
			continue
		}

		if len(s.Cols) == 0 {
			continue
		}

		// remove unavailable seat by cols
		newCols := make([]int64, 0)
		for _, c := range list[s.Row].Cols {
			if !utils.IsInt64SliceContains(s.Cols, c) {
				newCols = append(newCols, c)
			}
		}
		list[s.Row].Cols = newCols
	}

	return list
}

func ParseSeatToKey(seat *seat_models.Seat) string {
	return fmt.Sprintf("%v:%v", seat.Row, seat.Col)
}