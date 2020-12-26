package utils

import (
	"github.com/stretchr/testify/assert"
	"reflect"
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
		min, max := GetMaxMinRange(tc.Current, tc.Length, tc.Distance)
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
		dist := CalculateSeatDistance(tc.Seat, tc.Row, tc.Col)
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
		items := GetItemsInRange(tc.Min, tc.Max, tc.Ignore)
		t.Logf("Test data: %v", tc)
		t.Logf("Result: %v", items)
		assert.True(t, reflect.DeepEqual(items, tc.Items))
	}
}


