package utils

import (
	"math"
	seat_models "simple-restful/exmsgs/seat/models"
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

// Return map of seats with index is row and value is col
// The seat must be in range
// @params seat seat_models.Seat Current seat position
// @params numRow int64  Number row of table
// @params numCol int64  Number col of table
func GetAllSeatInRange(seat *seat_models.Seat, numRow int64, numCol int64, maxRange int64) map[int64]int64 {
	seats := make(map[int64]int64)
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

	availableRows = GetItemsInRange(minRow, maxRow, []int64{seat.Row})
	availableCols = GetItemsInRange(minCol, maxCol, []int64{seat.Col})

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
			seats[r] = c
		}
	}

	return seats
}
