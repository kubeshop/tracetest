package timing

import (
	"fmt"
	"time"
)

var Now = func() time.Time {
	return time.Now().UTC()
}

func TimeDiff(start, end time.Time) time.Duration {
	var endDate time.Time
	if !dateIsZero(end) {
		endDate = end.UTC()
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
	// Determine the number of digits in the timestamp
	numDigits := len(fmt.Sprintf("%d", timestamp))

	switch {
	case numDigits <= 10:
		// Timestamp is in seconds
		return time.Unix(timestamp, 0)
	case numDigits <= 13:
		// Timestamp is in milliseconds, convert to nanoseconds and then to time.Time
		return time.Unix(0, timestamp*int64(time.Millisecond))
	default:
		// Timestamp is in nanoseconds, convert directly to time.Time
		return time.Unix(0, timestamp)
	}
}

func MustParse(in string) time.Time {
	date, err := time.Parse(time.RFC3339Nano, in)
	if err != nil {
		panic(fmt.Sprintf("error parsing date: %s", err))
	}
	return date
}
