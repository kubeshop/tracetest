package telemetry

import (
	"fmt"
	"os"
	"runtime"

	pyroscope "github.com/grafana/pyroscope-go"
)

const (
	MutexProfileFractionRate = 5 // it samples 1/5 of data
	BlockProfileFractionRate = 5
)

func StartProfiler(applicationName, applicationEnvironment, telemetryServerAddress string, samplingPercentage int) {
	fmt.Printf("Starting profiler for environment [%s] with sampling [%d]. Telemetry server address: %s\n",
		applicationEnvironment, samplingPercentage, telemetryServerAddress)

	samplingRate := 100 / samplingPercentage

	runtime.SetMutexProfileFraction(samplingRate)
	runtime.SetBlockProfileRate(samplingRate)

	pyroscope.Start(pyroscope.Config{
		ApplicationName: applicationName,

		// replace this with the address of pyroscope server
		ServerAddress: telemetryServerAddress,

		// you can disable logging by setting this to nil
		// Logger: pyroscope.StandardLogger,
		Logger: nil,

		// you can provide static tags via a map:
		Tags: map[string]string{
			"hostname":    os.Getenv("HOSTNAME"),
			"environment": applicationEnvironment,
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
