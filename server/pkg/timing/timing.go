package timing

import (
	"math"
	"time"
)

var Now = func() time.Time {
	return time.Now().UTC()
}

func TimeDiff(start, end time.Time) time.Duration {
	var endDate time.Time
	if !dateIsZero(end) {
		endDate = end
	} else {
		endDate = Now()
	}
	return endDate.Sub(start)
}

func DurationInMillieconds(d time.Duration) int {
	return int(d.Milliseconds())
}

func DurationInNanoseconds(d time.Duration) int {
	return int(d.Nanoseconds())
}

func DurationInSeconds(d time.Duration) int {
	return int(math.Ceil(d.Seconds()))
}

func dateIsZero(in time.Time) bool {
	return in.IsZero() || in.Unix() == 0
}
