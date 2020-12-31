package tests

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"simple-restful/cmd/grpc/seat-svc/pkg/utils"
	seat_models "simple-restful/exmsgs/seat/models"
	"testing"
)

func TestGetMaxMinRange(t *testing.T) {
	tcs := []struct {
		Current  int64
		Distance int64
		Length   int64
		Min      int64
		Max      int64
	}{
		{
			Current:  2,
			Distance: 2,
			Length:   10,
			Min:      0,
			Max:      4,
		},
		{
			Current:  5,
			Distance: 3,
			Length:   10,
			Min:      2,
			Max:      8,
		},
		{
			Current:  7,
			Distance: 5,
			Length:   8,
			Min:      2,
			Max:      7,
		},
	}

	for _, tc := range tcs {
		min, max := utils.GetMaxMinRange(tc.Current, tc.Length, tc.Distance)
		t.Logf("Test data: %v", tc)
		assert.Equal(t, min, tc.Min)
		assert.Equal(t, max, tc.Max)
	}
}

func TestCalculateSeatDistance(t *testing.T) {
	tcs := []struct {
		Seat     *seat_models.Seat
		Row      int64
		Col      int64
		Distance int64
	}{
		{
			Seat: &seat_models.Seat{
				Row: 0,
				Col: 0,
			},
			Distance: 6,
			Row:      5,
			Col:      1,
		},
		{
			Seat: &seat_models.Seat{
				Row: 2,
				Col: 0,
			},
			Distance: 2,
			Row:      3,
			Col:      1,
		},
		{
			Seat: &seat_models.Seat{
				Row: 10,
				Col: 4,
			},
			Distance: 10,
			Row:      3,
			Col:      1,
		},
	}

	for _, tc := range tcs {
		dist := utils.CalculateSeatDistance(tc.Seat, tc.Row, tc.Col)
		t.Logf("Test data: %v", tc)
		assert.Equal(t, dist, tc.Distance)
	}
}

func TestGetItemsInRange(t *testing.T) {
	tcs := []struct {
		Items  []int64
		Min    int64
		Max    int64
		Ignore []int64
	}{
		{
			Min:    0,
			Max:    4,
			Ignore: []int64{3},
			Items:  []int64{0, 1, 2, 4},
		},
		{
			Min:    0,
			Max:    4,
			Ignore: []int64{3, 2},
			Items:  []int64{0, 1, 4},
		},
		{
			Min:    2,
			Max:    4,
			Ignore: []int64{3, 2},
			Items:  []int64{4},
		},
		{
			Min:    2,
			Max:    6,
			Ignore: []int64{3, 2, 5, 1},
			Items:  []int64{4, 6},
		},
		{
			Min:    2,
			Max:    6,
			Ignore: []int64{3, 2, 5, 1, 4, 6},
			Items:  []int64{},
		},
	}

	for _, tc := range tcs {
		items := utils.GetItemsInRange(tc.Min, tc.Max, tc.Ignore)
		t.Logf("Test data: %v", tc)
		t.Logf("Result: %v", items)
		assert.True(t, reflect.DeepEqual(items, tc.Items))
	}
}

func TestGetAllSeatInRange(t *testing.T) {
	tcs := []struct {
		Seat            *seat_models.Seat
		Row             int64
		Col             int64
		Distance        int64
		MapSeatsInRange map[int64][]int64
	}{
		{
			Seat: &seat_models.Seat{
				Col: 2,
				Row: 2,
			},
			Row:      4,
			Col:      4,
			Distance: 2,
			MapSeatsInRange: map[int64][]int64{
				0: {2},
				1: {1, 2, 3},
				2: {0, 1, 3},
				3: {1, 2, 3},
			},
		},
		{
			Seat: &seat_models.Seat{
				Col: 2,
				Row: 2,
			},
			Row:      3,
			Col:      3,
			Distance: 2,
			MapSeatsInRange: map[int64][]int64{
				0: {2},
				1: {1, 2},
				2: {0, 1},
			},
		},
		{
			Seat: &seat_models.Seat{
				Col: 5,
				Row: 5,
			},
			Row:      3,
			Col:      3,
			Distance: 2,
			MapSeatsInRange: nil,
		},
	}

	for _, tc := range tcs {
		items := utils.GetAllSeatInRange(tc.Seat, tc.Row, tc.Col, tc.Distance)
		t.Logf("Test data: %v", tc)
		t.Logf("Expected: %#v", tc.MapSeatsInRange)
		t.Logf("Result: %#v", items)
		assert.True(t, reflect.DeepEqual(items, tc.MapSeatsInRange))
	}
}
