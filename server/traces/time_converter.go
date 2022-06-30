package traces

import (
	"fmt"
	"math"
)

const (
	NANOSECOND_SCALE  = 0
	MICROSECOND_SCALE = 1
	MILLISECOND_SCALE = 2
	SECOND_SCALE      = 3
	MINUTE_SCALE      = 4
	HOUR_SCALE        = 5
)

var timeUnits = []string{"ns", "μs", "ms", "s", "m", "h"}

func ConvertNanoSecondsIntoProperTimeUnit(value int) string {
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
		return convertSecondsIntoPropertTimeUnit(convertedNumber)
	}

	return fmt.Sprintf("%.0f%s", convertedNumber, unit)
}

func convertSecondsIntoPropertTimeUnit(number float64) string {
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
