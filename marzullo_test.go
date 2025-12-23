package rmarzullo

import (
	"testing"
)

var marzullo = &Marzullo{MinOverlap: 0, ToleranceNs: 0}

func TestClosedIntervalSemantics(t *testing.T) {
	intervals := []IntervalWithSource{
		{Interval: Interval{LowNs: 10, HighNs: 20}, SourceId: 1},
		{Interval: Interval{LowNs: 20, HighNs: 30}, SourceId: 2},
	}
	result, err := marzullo.SmallestInterval(intervals)
	if err != nil {
		t.Fatalf("error getting smallest interval: %v", err)
	}
	if result.LowNs != 20 || result.HighNs != 20 {
		t.Fatalf("expected overlap at boundary 20, got %+v", result)
	}
}
func testSmallestInterval(
	t *testing.T,
	bounds []int64,
	expected Interval,
) {
	t.Helper()
	tuples := make([]IntervalWithSource, len(bounds)/2)
	for i := 0; i < len(bounds)/2; i++ {
		tuples[i] = IntervalWithSource{
			SourceId: i,
			Interval: Interval{
				LowNs:  bounds[i*2],
				HighNs: bounds[i*2+1],
			},
		}
	}
	interval, err := marzullo.SmallestInterval(tuples)
	if err != nil {
		t.Fatalf("error getting smallest interval: %v", err)
	}
	if interval != expected {

		t.Fatalf(
			"unexpected interval\nexpected: %+v\ngot: %+v",
			expected,
			interval,
		)
	}
}
func TestMarzullo(t *testing.T) {
	testSmallestInterval(t,
		[]int64{
			11, 13,
			10, 12,
			8, 12,
		},
		Interval{
			LowNs:  11,
			HighNs: 12,
		},
	)
	testSmallestInterval(t,
		[]int64{
			8, 12,
			11, 13,
			14, 15,
		},
		Interval{
			LowNs:  11,
			HighNs: 12,
		},
	)
	testSmallestInterval(t,
		[]int64{
			-10, 10,
			-1, 1,
			0, 0,
		},
		Interval{
			LowNs:  0,
			HighNs: 0,
		},
	)
	// The upper bound of the first interval overlaps inclusively with the lower of the last.
	testSmallestInterval(t,
		[]int64{
			8, 12,
			10, 11,
			8, 10,
		},
		Interval{
			LowNs:  10,
			HighNs: 10,
		},
	)

	// The first smallest interval is selected. The alternative with equal overlap is 10..12.
	// However, while this shares the same number of sources, it is not the smallest interval.
	testSmallestInterval(t,
		[]int64{
			8, 12,
			10, 12,
			8, 9,
		},
		Interval{
			LowNs:  8,
			HighNs: 9,
		},
	)
	// The last smallest interval is selected.
	testSmallestInterval(t,
		[]int64{
			7, 9,
			7, 12,
			10, 11,
		},
		Interval{
			LowNs:  10,
			HighNs: 11,
		},
	)
	// Negative offsets.
	testSmallestInterval(t,
		[]int64{
			-9, -7,
			-12, -7,
			-11, -10,
		},
		Interval{
			LowNs:  -11,
			HighNs: -10,
		},
	)
	// A cluster of one with no remote sources.
	testSmallestInterval(t,
		[]int64{},
		Interval{
			LowNs:  0,
			HighNs: 0,
		},
	)
	// A cluster of two with one remote source.
	testSmallestInterval(t,
		[]int64{
			1, 3,
		},
		Interval{
			LowNs: 1,

			HighNs: 3,
		},
	)
	// A cluster of three with agreement.
	testSmallestInterval(t,
		[]int64{
			1, 3,
			2, 2,
		},
		Interval{
			LowNs:  2,
			HighNs: 2,
		},
	)
	// A cluster of three with no agreement.
	testSmallestInterval(t,
		[]int64{
			1, 3,
			4, 5,
		},
		Interval{
			LowNs:  4,
			HighNs: 5,
		},
	)
}
