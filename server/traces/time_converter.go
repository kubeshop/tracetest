package traces

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

const (
	NANOSECOND_SCALE  = 0
	MICROSECOND_SCALE = 1
	MILLISECOND_SCALE = 2
	SECOND_SCALE      = 3
	MINUTE_SCALE      = 4
	HOUR_SCALE        = 5
)

var timeUnits = []string{"ns", "us", "ms", "s", "m", "h"}

func ConvertNanoSecondsIntoProperTimeUnit(value int) string {
	if value <= 0 {
		return "0ns"
	}

	// Scale is basically how many times we can divide value by 1000.
	// To achieve that, we use Log1000(value). However, go doesn't support that,
	// so we need to use the math transformation:
	// Log1000(value) = Log(value) / Log(1000)
	//
	// As we don't care about the floating point part of the result, we use Floor to remove it
	scale := math.Floor(math.Log(float64(value)) / math.Log(1000))
	if scale > SECOND_SCALE {
		scale = SECOND_SCALE
	}

	divisor := math.Pow(1000, scale)
	convertedNumber := float64(value) / divisor
	unit := timeUnits[int(scale)]

	if scale == SECOND_SCALE {
		return convertSecondsIntoProperTimeUnit(convertedNumber)
	}

	return fmt.Sprintf("%.0f%s", convertedNumber, unit)
}

func convertSecondsIntoProperTimeUnit(number float64) string {
	scale := SECOND_SCALE
	if number >= 60 {
		number = number / 60
		scale = MINUTE_SCALE
	}

	if number >= 60 {
		number = number / 60
		scale = HOUR_SCALE
	}

	unit := timeUnits[scale]

	if isWholeNumber(number) {
		return fmt.Sprintf("%.0f%s", number, unit)
	}

	return fmt.Sprintf("%.1f%s", number, unit)
}

func isWholeNumber(number float64) bool {
	return math.Mod(number, 1.0) == 0
}

var timeFieldRegex = regexp.MustCompile(`^([0-9]+(\.[0-9]+)?)(ns|us|ms|s|m|h)$`)
var unitScale = map[string]time.Duration{
	"ns": time.Nanosecond,
	"us": time.Microsecond,
	"ms": time.Millisecond,
	"s":  time.Second,
	"m":  time.Minute,
	"h":  time.Hour,
}

func ConvertTimeFieldIntoNanoSeconds(value string) int {
	result := timeFieldRegex.FindAllStringSubmatch(value, -1)
	if len(result) == 0 {
		// Maybe it is a number, try to convert it for backwards compatibility sake
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0
		}

		return intValue
	}
	number, err := strconv.ParseFloat(result[0][1], 64)
	if err != nil {
		return 0
	}

	unit := result[0][3]
	scale := float64(unitScale[unit])

	return int(number * scale)
}
