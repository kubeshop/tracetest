package telemetry

import (
	"os"
	"runtime"

	pyroscope "github.com/grafana/pyroscope-go"
)

const (
	MutexProfileFractionRate = 5 // it samples 1/5 of data
	BlockProfileFractionRate = 5
)

func StartProfiler(applicationName, telemetryServer string, samplingPercentage int) {
	samplingRate := 100 / samplingPercentage

	runtime.SetMutexProfileFraction(samplingRate)
	runtime.SetBlockProfileRate(samplingRate)

	pyroscope.Start(pyroscope.Config{
		ApplicationName: applicationName,

		// replace this with the address of pyroscope server
		ServerAddress: telemetryServer,

		// you can disable logging by setting this to nil
		Logger: pyroscope.StandardLogger,

		// you can provide static tags via a map:
		Tags: map[string]string{
			"hostname": os.Getenv("HOSTNAME"),
		},

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
}
