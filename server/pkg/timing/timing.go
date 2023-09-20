package timing

import (
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

func dateIsZero(in time.Time) bool {
	return in.IsZero() || in.Unix() == 0
}

// ParseUnix parses a unix timestamp into a time.Time
// it accepts an integer which can be either milli or nano
func ParseUnix(timestamp int64) time.Time {
	// Determine the range of the timestamp to know if it is nano or milli
	// we can assume that any timestamp less than this is milliq
	const threshold = int64(1e12)

	if timestamp < threshold {
		// is milli
		return time.Unix(0, timestamp*int64(time.Millisecond))
	}
	// is nano
	return time.Unix(0, timestamp)
}
