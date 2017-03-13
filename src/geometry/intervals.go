package geometry

import (
	"math"
)

type Interval struct {
	lower float64
	upper float64
}

func NewInterval(lower float64, upper float64) *Interval {
	i := &Interval{lower, upper}
	return i
}

func (inter1 Interval) expand(inter2 Interval) Interval {
	result := &Interval{math.Min(inter1.lower, inter2.lower), math.Max(inter1.upper, inter2.upper)}
	return *result
}