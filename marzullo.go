// Package rmarzullo implements Marzullo's algorithm for
// finding a quorum-consistent interval across multiple sources.
//
// All time values are expressed in nanoseconds.
package rmarzullo

import (
	"fmt"
	"sort"
)

type Interval struct {
	LowNs  int64
	HighNs int64
}
type IntervalWithSource struct {
	Interval // All time values are expressed in nanoseconds.
	SourceId int
}
type Marzullo struct {
	MinOverlap  int   // the minimum number of overlapping intervals required , Quoram size
	ToleranceNs int64 // tolerance in Nanoseconds
}
type BoundaryPoint struct {
	SourceId     int
	Offset       int64
	BoundaryType int // 0 - lower boundary, 1 - upper boundary
}

const (
	LOWER_BOUNDARY int = iota
	UPPER_BOUNDARY
) // 0 - lower boundary, 1 - upper boundary
func (m *Marzullo) SmallestInterval(intervals []IntervalWithSource) (Interval, error) {
	boundaryPoints, err := boundaryPoints(intervals, m.ToleranceNs)
	if err != nil {
		// Handle the error appropriately
		return Interval{}, err
	}
	count := 0
	best := 0
	var bestInterval Interval
	for i := 0; i < len(boundaryPoints)-1; i++ {
		boundary := boundaryPoints[i]
		if boundary.BoundaryType == LOWER_BOUNDARY {
			count++
		} else {
			count--
		}
		if count < m.MinOverlap {
			continue

		}
		if count > best {
			best = count
			bestInterval = Interval{
				LowNs:  boundary.Offset,
				HighNs: boundaryPoints[i+1].Offset,
			}
		} else if count == best {
			currentInterval := Interval{
				LowNs:  boundary.Offset,
				HighNs: boundaryPoints[i+1].Offset,
			}
			if (currentInterval.HighNs - currentInterval.LowNs) < (bestInterval.HighNs -
				bestInterval.LowNs) {

				bestInterval = currentInterval
			}
		}
	}
	if best < m.MinOverlap {
		return Interval{}, fmt.Errorf("could not find an interval with minimum overlap %v",
			m.MinOverlap)
	}
	return bestInterval, nil
}
func boundaryPoints(intervals []IntervalWithSource, tolerance int64) ([]BoundaryPoint, error) {
	result := make([]BoundaryPoint, 0, len(intervals)*2)
	seenSources := make(map[int]struct{})
	for _, interval := range intervals {
		var source = interval.SourceId
		if _, ok := seenSources[interval.SourceId]; ok {
			return nil, fmt.Errorf("duplicate source %d", interval.SourceId)
		}
		seenSources[interval.SourceId] = struct{}{}
		var boundary = BoundaryPoint{
			SourceId:     source,
			Offset:       interval.LowNs - tolerance,
			BoundaryType: LOWER_BOUNDARY,
		}
		result = append(result, boundary)
		boundary = BoundaryPoint{
			SourceId:     source,
			Offset:       interval.HighNs + tolerance,
			BoundaryType: UPPER_BOUNDARY,
		}
		result = append(result, boundary)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Offset < result[j].Offset {
			return true
		}
		if result[i].Offset == result[j].Offset {
			if result[i].BoundaryType == result[j].BoundaryType {
				return result[i].SourceId < result[j].SourceId
			}
			return result[i].BoundaryType < result[j].BoundaryType

		}
		return false
	})
	return result, nil
}
